package api

import (
	"encoding/gob"
	"log/slog"

	"adelhub.com/voiceline/internal/db"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/validator/v10"
)

type Dependencies struct {
	Logger   slog.Logger
	Store    *db.SqlStore
	Validate *validator.Validate
	Sm       *scs.SessionManager
}

type Api struct {
	Dependencies
}

func New(deps Dependencies) *Api {
	gob.Register(db.User{})
	return &Api{Dependencies: deps}
}
