package api

import (
	"context"
	"net/http"
	"slices"

	"adelhub.com/voiceline/internal/db"
)

type ContextKey string

const ContextUserKey ContextKey = "user"

type Middleware func(http.Handler) http.Handler
type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {

	for _, middleware := range slices.Backward(middlewares) {
		handler = middleware(handler)
	}

	return handler
}

func ChainHandlerFunc(handler http.HandlerFunc, middlewares ...MiddlewareFunc) http.HandlerFunc {
	for _, middleware := range slices.Backward(middlewares) {
		handler = middleware(handler)
	}

	return handler
}

func (api *Api) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := api.Sm.Get(r.Context(), userSessionKey).(db.User)
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		authCtx := context.WithValue(r.Context(), ContextUserKey, user)
		next.ServeHTTP(w, r.WithContext(authCtx))
	}
}
