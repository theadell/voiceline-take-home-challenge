package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"adelhub.com/voiceline/internal/db"
	"github.com/go-playground/validator/v10"
)

const maxBodySize = 1_024 * 1_024

func ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(dst)
	if err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError
		switch {

		case errors.As(err, &unmarshalTypeError):
			return errors.New("incorrect JSON type for field " + unmarshalTypeError.Field)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("unexpected end of JSON input")

		case errors.Is(err, io.EOF):
			return errors.New("request body cannot be empty")

		default:
			return errors.New("invalid JSON provided")
		}
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func (api *Api) errorResponse(w http.ResponseWriter, status int, message any) {
	env := envelope{"error": message}
	err := WriteJSON(w, status, env)
	if err != nil {
		api.Logger.Error("failed to write JSON error response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (api *Api) badRequestResponse(w http.ResponseWriter, message any) {
	api.errorResponse(w, http.StatusBadRequest, message)
}

func (api *Api) serverError(w http.ResponseWriter, err error, args ...any) {
	api.Logger.Error(err.Error(), args...)
	message := "the server could not process your request"
	api.errorResponse(w, http.StatusBadRequest, message)
}

func (api *Api) validationErrorResponse(w http.ResponseWriter, err error) {

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)
		for _, e := range validationErrors {
			errors[e.Field()] = e.Tag()
		}
		WriteJSON(w, http.StatusBadRequest, errors)
		return
	}

	http.Error(w, "Invalid input", http.StatusBadRequest)
}

func (api *Api) setUserSession(w http.ResponseWriter, r *http.Request, user db.User) {
	api.Sm.Put(r.Context(), userSessionKey, user)

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	WriteJSON(w, http.StatusOK, SessionResponse{
		User:          UserResponse{Email: user.Email},
		SessionActive: true,
	})
}
