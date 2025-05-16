package erd_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Package-level DB connection for tests
var testDB *sqlx.DB

func DBSetup(t *testing.T) (*sqlx.DB, func(), error) {
	t.Helper()

	ctx := context.Background()

	pgc, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithInitScripts(filepath.Join("..", "..", "testdata", "schema.sql")),
		postgres.WithDatabase("test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithSQLDriver("sqlx"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, nil, err
	}

	dsn, err := pgc.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = pgc.Terminate(ctx)
		return nil, nil, err
	}

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		_ = pgc.Terminate(ctx)
		return nil, nil, err
	}

	teardown := func() {
		_ = db.Close()
		_ = pgc.Terminate(ctx)
	}

	return db, teardown, nil
}

func TestMain(m *testing.M) {
	// Setup the database
	db, teardown, err := DBSetup(&testing.T{})
	if err != nil {
		panic("Failed to setup test database: " + err.Error())
	}
	defer teardown()

	// Store DB in package-level variable
	testDB = db

	// Run all tests
	code := m.Run()
	os.Exit(code)
}

// GetTestDB returns the shared database connection for tests
func GetTestDB() *sqlx.DB {
	return testDB
}
