package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"adelhub.com/voiceline/internal/db"
	"github.com/mattn/go-sqlite3"
)

func (api *Api) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	err := ReadJSON(w, r, &req)
	if err != nil {
		api.badRequestResponse(w, err.Error())
		return
	}

	err = api.Validate.Struct(req)
	if err != nil {
		api.validationErrorResponse(w, err)
		return
	}
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		api.serverError(w, err)
		return
	}
	_, err = api.Store.CreateUser(context.TODO(), db.CreateUserParams{Email: req.Email, PasswordHash: sql.NullString{Valid: true, String: hashedPassword}})
	if err != nil {
		var sqlErr sqlite3.Error
		if errors.As(err, &sqlErr) {
			if sqlErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				api.badRequestResponse(w, "A user with that email already exists")
			}
		} else {
			api.serverError(w, err)
		}
	}
	w.WriteHeader(http.StatusCreated)
}

func (api *Api) passwordLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req PasswordLoginRequest

	err := ReadJSON(w, r, &req)
	if err != nil {
		api.badRequestResponse(w, err.Error())
		return
	}

	loginErr := "incorrect username or password"
	user, err := api.Store.GetUserByEmail(context.TODO(), req.Email)
	if err != nil {
		api.badRequestResponse(w, loginErr)
		return
	}
	if !user.PasswordHash.Valid {
		api.badRequestResponse(w, loginErr)
		return
	}
	if err := CheckPassword(user.PasswordHash.String, req.Password); err != nil {
		api.badRequestResponse(w, loginErr)
		return
	}

	api.setUserSession(w, r, user)
}
