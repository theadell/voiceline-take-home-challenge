package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/alexedwards/scs/v2"
	"golang.org/x/oauth2"
)

func TestOauth2AuthRequestHandler(t *testing.T) {

	sessionManager := scs.New()
	sessionManager.Cookie.Name = "GO_SESSION_ID"

	mockProvider := &oauth2.Config{
		ClientID:     "mock-client-id",
		ClientSecret: "mock-client-secret",
		Endpoint: oauth2.Endpoint{
			AuthURL: "https://example.com/auth",
		},
		RedirectURL: "https://example.com/callback",
	}

	tests := []struct {
		name           string
		provider       string
		queryNext      string
		mockProviders  map[string]*oauth2.Config
		expectedStatus int
	}{
		{
			name:           "unsupported provider",
			provider:       "unsupported_provider",
			mockProviders:  map[string]*oauth2.Config{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "successful redirection",
			provider:       "google",
			mockProviders:  map[string]*oauth2.Config{"google": mockProvider},
			queryNext:      "/profile",
			expectedStatus: http.StatusSeeOther,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the API instance
			apiInstance := New(Dependencies{
				Providers: tt.mockProviders,
				Sm:        sessionManager,
			})

			req := httptest.NewRequest(http.MethodGet, "/auth/"+tt.provider+"?next="+tt.queryNext, nil)
			req.SetPathValue("provider", tt.provider)

			recorder := httptest.NewRecorder()

			handler := apiInstance.Sm.LoadAndSave(http.HandlerFunc(apiInstance.oauth2AuthRequestHandler))

			handler.ServeHTTP(recorder, req)

			mustEqual(t, recorder.Code, tt.expectedStatus)

			if tt.expectedStatus == http.StatusSeeOther {
				redirectURL := recorder.Header().Get("Location")

				parsedURL, err := url.Parse(redirectURL)
				mustOk(t, err)

				queryParams := parsedURL.Query()

				mustTrue(t, queryParams.Get("state") != "", "missing 'state' parameter in redirect URL")
				mustTrue(t, queryParams.Get("nonce") != "", "missing 'nonce' parameter in redirect URL")
				mustTrue(t, queryParams.Get("code_challenge") != "", "missing 'code_challenge' parameter in redirect URL")
				mustTrue(t, queryParams.Get("code_challenge_method") == "S256", "code_challenge_method should be 'S256'")
			}
		})
	}
}
