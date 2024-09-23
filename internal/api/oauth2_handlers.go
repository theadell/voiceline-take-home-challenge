package api

import (
	"context"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
)

func (api *Api) oauth2AuthRequestHandler(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	if !api.isProviderSupported(provider) {
		api.badRequestResponse(w, "unsupported oauth2 provider")
		return
	}
	next := r.URL.Query().Get("next")
	if next == "" {
		next = "/"
	}

	state := OAuth2State{
		Provider: provider,
		State:    oauth2.GenerateVerifier(),
		Verifier: oauth2.GenerateVerifier(),
		Nonce:    oauth2.GenerateVerifier(),
		Next:     next,
	}

	provder, ok := api.Providers[provider]
	if !ok {
		api.serverError(w, errors.New("failed to obtain oauth2 provider"), "proider", provder, "providers map", api.Providers)
		return
	}

	authURL := provder.AuthCodeURL(state.State, oauth2.S256ChallengeOption(state.Verifier), oauth2.SetAuthURLParam("nonce", state.Nonce))
	api.Sm.Put(r.Context(), oauth2StateKey, state)
	http.Redirect(w, r, authURL, http.StatusSeeOther)
}

func (api *Api) oauth2CallbackHandler(w http.ResponseWriter, r *http.Request) {

	oauth2State, ok := api.Sm.Pop(r.Context(), oauth2StateKey).(OAuth2State)
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	provider, providerExists := api.Providers[oauth2State.Provider]
	if !ok || oauth2State == (OAuth2State{}) || !providerExists || state != oauth2State.State || code == "" {
		api.badRequestResponse(w, "invalid oauth2 callback request")
		return
	}

	token, err := provider.Exchange(context.TODO(), code, oauth2.VerifierOption(oauth2State.Verifier))
	if err != nil {
		api.serverError(w, err)
		return
	}

	rawIdToken, ok := token.Extra("id_token").(string)
	if !ok {
		api.serverError(w, errors.New("token response doesn't contain id token"))
		return
	}
	idToken, err := api.ProvidersValidator.ValidateJWT(rawIdToken)
	if err != nil {
		api.serverError(w, errors.New("invalid id token"), "error", err, "token", rawIdToken)
		return
	}

	if nonce := idToken.GetStringClaim("nonce"); nonce != oauth2State.Nonce || idToken.Email == "" || idToken.Subject == "" {
		api.serverError(w, errors.New("invalid id token claims"), "state", oauth2State, "token", idToken)
		return
	}

	user, err := api.Store.LoginWithProvider(context.TODO(), oauth2State.Provider, idToken.Subject, idToken.Email)
	if err != nil {
		api.serverError(w, errors.New("failed to create or link user"), "error", err)
		return
	}

	api.Sm.Put(r.Context(), userSessionKey, user)

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	http.Redirect(w, r, oauth2State.Next, http.StatusSeeOther)
}

func (api *Api) oauth2SessionHandler(w http.ResponseWriter, r *http.Request) {
	var req Oauth2SessionRequest

	err := ReadJSON(w, r, &req)
	if err != nil {
		api.badRequestResponse(w, err.Error())
		return
	}

	idToken, err := api.ProvidersValidator.ValidateJWT(req.IdToken)
	if err != nil {
		api.badRequestResponse(w, err.Error())
		return
	}

	if !api.isProviderSupported(req.Provider) || idToken.Email == "" || idToken.Subject == "" {
		api.badRequestResponse(w, "invalid id tokens calims")
		return
	}

	user, err := api.Store.LoginWithProvider(context.TODO(), req.Provider, idToken.Subject, idToken.Email)
	if err != nil {
		api.serverError(w, errors.New("failed to create or link user"), "error", err)
		return
	}

	api.setUserSession(w, r, user)
}
