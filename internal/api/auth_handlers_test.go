package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"adelhub.com/voiceline/internal/db"
	"adelhub.com/voiceline/internal/db/mock"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/validator/v10"
	"github.com/mattn/go-sqlite3"
	"go.uber.org/mock/gomock"
)

func TestRegisterUserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock.NewMockQuerier(ctrl)
	validate := validator.New()
	store := db.NewSqlStore(nil, mockQuerier)
	deps := Dependencies{
		Store:    store,
		Validate: validate,
	}

	apiInstance := New(deps)

	tests := []struct {
		name           string
		requestBody    CreateUserRequest
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "successful registration",
			requestBody: CreateUserRequest{
				Email:    "user@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockQuerier.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(db.User{}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "validation error - missing email",
			requestBody: CreateUserRequest{
				Password: "password123",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "duplicate email error",
			requestBody: CreateUserRequest{
				Email:    "existing@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockQuerier.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(db.User{}, sqlite3.Error{
						ExtendedCode: sqlite3.ErrConstraintPrimaryKey,
					})
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockSetup()

			body, err := json.Marshal(tt.requestBody)
			mustOk(t, err)

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			apiInstance.registerUserHandler(recorder, req)

			mustEqual(t, recorder.Code, tt.expectedStatus)
		})
	}
}

func TestPasswordLoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := mock.NewMockQuerier(ctrl)
	sessionManager := scs.New()
	sessionManager.Lifetime = 3 * time.Hour
	sessionManager.IdleTimeout = 20 * time.Minute
	sessionManager.Cookie.Name = "GO_SESSION_ID"

	password := "password123"
	hashedPassword, err := HashPassword(password)
	mustOk(t, err)

	store := db.NewSqlStore(nil, mockQuerier)
	deps := Dependencies{
		Store:    store,
		Sm:       sessionManager,
		Validate: validator.New(),
	}

	apiInstance := New(deps)

	tests := []struct {
		name           string
		requestBody    PasswordLoginRequest
		mockSetup      func()
		expectedStatus int
		expectSession  bool
	}{
		{
			name: "successful login",
			requestBody: PasswordLoginRequest{
				Email:    "user@example.com",
				Password: password,
			},
			mockSetup: func() {
				mockQuerier.EXPECT().
					GetUserByEmail(gomock.Any(), "user@example.com").
					Return(db.User{
						Email: "user@example.com",
						PasswordHash: sql.NullString{
							Valid:  true,
							String: hashedPassword,
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectSession:  true,
		},
		{
			name: "incorrect password",
			requestBody: PasswordLoginRequest{
				Email:    "user@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func() {
				mockQuerier.EXPECT().
					GetUserByEmail(gomock.Any(), "user@example.com").
					Return(db.User{
						Email: "user@example.com",
						PasswordHash: sql.NullString{
							Valid:  true,
							String: "password123",
						},
					}, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectSession:  false,
		},
		{
			name: "user not found",
			requestBody: PasswordLoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockQuerier.EXPECT().
					GetUserByEmail(gomock.Any(), "nonexistent@example.com").
					Return(db.User{}, errors.New("user not found"))
			},
			expectedStatus: http.StatusBadRequest,
			expectSession:  false,
		},
		{
			name: "provider user with no password",
			requestBody: PasswordLoginRequest{
				Email:    "user@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockQuerier.EXPECT().
					GetUserByEmail(gomock.Any(), "user@example.com").
					Return(db.User{
						Email:        "user@example.com",
						PasswordHash: sql.NullString{Valid: false},
					}, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectSession:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockSetup()

			body, err := json.Marshal(tt.requestBody)
			mustOk(t, err)

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			handler := sessionManager.LoadAndSave(http.HandlerFunc(apiInstance.passwordLoginHandler))

			handler.ServeHTTP(recorder, req)

			mustEqual(t, recorder.Code, tt.expectedStatus)

			if tt.expectSession {
				cookies := recorder.Result().Cookies()

				foundSession := false
				for _, cookie := range cookies {
					if cookie.Name == sessionManager.Cookie.Name {
						foundSession = true
					}
				}
				mustTrue(t, foundSession, "Session cookie should be set")
			} else {
				mustEqual(t, len(recorder.Result().Cookies()), 0)
			}
		})
	}
}

func mustOk(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatal(err)
	}
}
func mustTrue(tb testing.TB, f bool, msg string) {
	tb.Helper()
	if !f {
		tb.Fatal(msg)
	}
}

func mustEqual[T any](tb testing.TB, have, want T) {
	tb.Helper()
	if !reflect.DeepEqual(have, want) {
		tb.Fatalf("\nhave: %+v\nwant: %+v\n", have, want)
	}
}
