package api

const userSessionKey = "user-session"

type envelope map[string]any

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type PasswordLoginRequest = CreateUserRequest

type UserResponse struct {
	Email string `json:"email" `
}

type SessionResponse struct {
	User          UserResponse `json:"user"`
	SessionActive bool         `json:"session_active"`
}
