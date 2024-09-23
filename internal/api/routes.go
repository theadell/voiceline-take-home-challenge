package api

import "net/http"

func (api *Api) Router() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/signup", api.registerUserHandler)
	mux.HandleFunc("POST /auth/userinfo", ChainHandlerFunc(api.userinfoHandler, api.Authenticate))

	mux.HandleFunc("POST /auth/session", api.passwordLoginHandler)
	mux.HandleFunc("DELETE /auth/session", ChainHandlerFunc(api.logoutHandler, api.Authenticate))

	mux.HandleFunc("GET /oauth2/{provider}/auth", api.oauth2AuthRequestHandler)
	mux.HandleFunc("GET /oauth2/{provider}/callback", api.oauth2CallbackHandler)

	mux.HandleFunc("POST /oauth2/{provider}/session", api.oauth2SessionHandler)

	return Chain(mux, api.Sm.LoadAndSave)
}
