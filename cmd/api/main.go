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

	"adelhub.com/voiceline/internal/config"
	"github.com/golang-migrate/migrate/v4"

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

	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:           http.NewServeMux(),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	startServer(server, logger, cfg.Port)

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

func startServer(server *http.Server, logger *slog.Logger, port int) {
	go func() {
		logger.Info(fmt.Sprintf("server is running on port %d", port))
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
