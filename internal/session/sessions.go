package session

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
)

func NewManager(db *sql.DB) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Lifetime = 3 * time.Hour
	sessionManager.IdleTimeout = 20 * time.Minute
	sessionManager.Cookie.Name = "GO_SESSION_ID"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = true
	sessionManager.Store = sqlite3store.New(db)
	return sessionManager
}
