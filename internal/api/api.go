package api

import (
	"encoding/gob"
	"log/slog"

	"adelhub.com/voiceline/internal/db"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/validator/v10"
	"github.com/theadell/authress"
	"golang.org/x/oauth2"
)

type Dependencies struct {
	Logger             slog.Logger
	Store              *db.SqlStore
	Validate           *validator.Validate
	Sm                 *scs.SessionManager
	Providers          map[string]*oauth2.Config
	ProvidersValidator *authress.Validator
}

type Api struct {
	Dependencies
}

func New(deps Dependencies) *Api {
	gob.Register(db.User{})
	gob.Register(OAuth2State{})
	return &Api{Dependencies: deps}
}
