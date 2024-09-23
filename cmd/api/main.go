package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"adelhub.com/voiceline/internal/api"
	"adelhub.com/voiceline/internal/config"
	"adelhub.com/voiceline/internal/db"
	"adelhub.com/voiceline/internal/session"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	"github.com/theadell/authress"
	"golang.org/x/oauth2"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := Run(); err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
}

func Run() error {

	cfg := config.LoadConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dbConnPool, err := createDatabaseConnectionPool(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer dbConnPool.Close()

	if cfg.MigrateDB {
		if err := migrateDatabase(cfg.DatabaseURL, logger); err != nil {
			return fmt.Errorf("failed to migrate database: %w", err)
		}
	}
	providersValidator, providers, err := createOAuth2Config(*cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize OAuth2 config: %w", err)
	}

	validate := validator.New()

	store, sessionManager := createDatabaseSqlStore(dbConnPool)

	apiInstance := api.New(api.Dependencies{
		Logger:             *logger,
		Store:              store,
		Validate:           validate,
		Sm:                 sessionManager,
		Providers:          providers,
		ProvidersValidator: providersValidator,
	})

	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:           apiInstance.Router(),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	startServer(server, logger, cfg.Host, cfg.Port)

	waitForShutdown(server, logger)

	return nil
}

func createDatabaseConnectionPool(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func migrateDatabase(databaseURL string, logger *slog.Logger) error {
	logger.Info("applying database migrations")
	m, err := migrate.New("file://migrations", fmt.Sprintf("sqlite3://%s", databaseURL))
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	if err == migrate.ErrNoChange {
		logger.Info("No new migrations to apply")
	} else {
		logger.Info("Database Migrations applied successfully")
	}

	return nil
}

func createOAuth2Config(cfg config.Config) (*authress.Validator, map[string]*oauth2.Config, error) {
	providersValidator, err := authress.New(authress.WithDiscovery(cfg.OAuth2DiscoveryUrl))
	if err != nil {
		return nil, nil, err
	}

	oauth2ClientConfig := &oauth2.Config{
		ClientID:     cfg.OAuth2ClientId,
		ClientSecret: cfg.OAuth2ClientSecret,
		RedirectURL:  cfg.Oauth2CallbackUri,
		Endpoint:     providersValidator.ClientEndpoint(),
		Scopes:       []string{"email", "profile", "openid"},
	}

	providers := map[string]*oauth2.Config{
		cfg.Oauth2ProviderName: oauth2ClientConfig,
	}

	return providersValidator, providers, nil
}

func createDatabaseSqlStore(dbConnPool *sql.DB) (*db.SqlStore, *scs.SessionManager) {
	querier := db.New(dbConnPool)
	store := db.NewSqlStore(dbConnPool, querier)
	sessionManager := session.NewManager(dbConnPool)
	return store, sessionManager
}

func startServer(server *http.Server, logger *slog.Logger, network string, port int) {
	go func() {
		logger.Info("server is running", "network", network, "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed to start", "err", err)
			os.Exit(1)
		}
	}()
}

func waitForShutdown(server *http.Server, logger *slog.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "err", err)
	} else {
		logger.Info("server has shut down gracefully")
	}
}
