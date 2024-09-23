package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

var (
	sqlStore *SqlStore
	rnd      *rand.Rand
	ctx      = context.TODO()
)

func TestMain(m *testing.M) {

	tempFile, err := os.CreateTemp("", "testdb-*.sqlite")
	if err != nil {
		log.Fatalf("failed to create temp SQLite file: %v", err)
	}

	defer os.Remove(tempFile.Name())

	db, err := sql.Open("sqlite3", tempFile.Name())
	if err != nil {
		log.Fatalf("failed to create SQLite database: %v", err)
	}
	defer db.Close()

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("failed to create migration driver: %v", err)
	}

	mi, err := migrate.NewWithDatabaseInstance("file://../../migrations", "sqlite3", driver)
	if err != nil {
		log.Fatalf("failed to create migration object: %v", err)
	}

	err = mi.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	rnd = rand.New(rand.NewSource(time.Now().Unix()))
	q := New(db)
	sqlStore = NewSqlStore(db, q)

	os.Exit(m.Run())
}

func generateUniqueEmail() string {
	return fmt.Sprintf("user-%d@example.com", rnd.Int63())
}

func generateUniqueProviderUserID(provider string) string {
	return fmt.Sprintf("%s-user-%d", provider, rnd.Int63())
}

func mustOk(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatal(err)
	}
}

func mustEqual[T any](tb testing.TB, have, want T) {
	tb.Helper()
	if !reflect.DeepEqual(have, want) {
		tb.Fatalf("\nhave: %+v\nwant: %+v\n", have, want)
	}
}
