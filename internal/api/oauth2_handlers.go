package api

import (
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
