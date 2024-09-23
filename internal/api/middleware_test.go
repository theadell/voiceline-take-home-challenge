package api

import (
	"encoding/gob"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"

	"adelhub.com/voiceline/internal/db"
	"github.com/alexedwards/scs/v2"
)

func TestChain(t *testing.T) {

	endpointHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Final"))
	})

	middlewareA := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("A-"))
			next.ServeHTTP(w, r)
		})
	}

	middlewareB := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("B-"))
			next.ServeHTTP(w, r)
		})
	}

	middlewareC := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("C-"))
			next.ServeHTTP(w, r)
		})
	}

	adapt := func(h http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) { h.ServeHTTP(w, r) }
	}

	adaptMW := func(h Middleware) MiddlewareFunc {
		return func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				h(http.HandlerFunc(next)).ServeHTTP(w, r)
			}
		}
	}

	adaptAll := func(h ...Middleware) []MiddlewareFunc {
		all := make([]MiddlewareFunc, len(h))
		for i, handler := range slices.All(h) {
			all[i] = adaptMW(handler)
		}
		return all
	}

	tests := []struct {
		name           string
		handlerCreator func(http.Handler) http.Handler
	}{

		{"Test with ChainHandlerFunc",
			func(h http.Handler) http.Handler {
				return ChainHandlerFunc(adapt(h), adaptAll(middlewareA, middlewareB, middlewareC)...)
			}},

		{"Test with Chain (http.Handler)", func(h http.Handler) http.Handler {
			return Chain(h, middlewareA, middlewareB, middlewareC)
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := test.handlerCreator(endpointHandler)

			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			response := w.Body.String()

			expected := "A-B-C-Final"

			if strings.TrimSpace(response) != expected {
				t.Errorf("Expected response %q, but got %q", expected, response)
			}
		})
	}
}

func TestAuthenticateMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		sessionData    any
		expectedStatus int
	}{
		{
			name:           "Authenticated",
			sessionData:    db.User{ID: 1, Email: "test"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unauthenticated",
			sessionData:    nil,
			expectedStatus: http.StatusUnauthorized,
		},
	}
	gob.Register(db.User{})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			sm := scs.New()
			sm.Lifetime = 1 * time.Hour

			api := &Api{Dependencies: Dependencies{Sm: sm}}

			endpointHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Success"))
			})

			handler := ChainHandlerFunc(endpointHandler, loadAndSaveMockHandlerfunc(sm, userSessionKey, tt.sessionData), api.Authenticate)

			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, but got %d", tt.expectedStatus, w.Code)
			}

		})
	}

}

func loadAndSaveMockHandlerfunc(session *scs.SessionManager, key string, value any) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			session.LoadAndSave(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				session.Put(r.Context(), key, value)
				next.ServeHTTP(w, r)
			})).ServeHTTP(w, r)
		}
	}
}
