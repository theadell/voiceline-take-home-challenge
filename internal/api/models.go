package api

const userSessionKey = "user-session"
const oauth2StateKey = "oauth2-req-state"

type envelope map[string]any

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type PasswordLoginRequest = CreateUserRequest

type Oauth2SessionRequest struct {
	Provider string `json:"provider" validate:"required"`
	IdToken  string `json:"id_token" validate:"required"`
}

type UserResponse struct {
	Email string `json:"email" `
}

type SessionResponse struct {
	User          UserResponse `json:"user"`
	SessionActive bool         `json:"session_active"`
}

type OAuth2State struct {
	Provider string
	Next     string
	State    string
	Nonce    string
	Verifier string
}
