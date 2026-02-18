package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"trello-clone/internal/database"
)

// SetupDB creates (if needed) a per-package test database, runs migrations,
// truncates all tables, and returns a *sql.DB connected to it.
// Each calling package gets its own database, so packages can run in parallel.
func SetupDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://trello:trello@127.0.0.1:5432/trello?sslmode=disable"
	}

	// Derive a unique DB name from the calling package's directory.
	_, file, _, _ := runtime.Caller(1)
	pkg := filepath.Base(filepath.Dir(file))
	dbName := "trello_test_" + sanitize(pkg)

	ctx := context.Background()

	// Connect to the base DB to create the test database.
	baseDB, err := database.Open(ctx, dsn)
	if err != nil {
		t.Fatalf("open base db: %v", err)
	}

	var exists bool
	if err := baseDB.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1)", dbName,
	).Scan(&exists); err != nil {
		baseDB.Close()
		t.Fatalf("check db: %v", err)
	}
	if !exists {
		// CREATE DATABASE cannot run inside a transaction.
		if _, err := baseDB.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName)); err != nil {
			baseDB.Close()
			t.Fatalf("create db %s: %v", dbName, err)
		}
	}
	baseDB.Close()

	// Connect to the per-package test database.
	testDSN := replaceDSNDatabase(dsn, dbName)
	db, err := database.Open(ctx, testDSN)
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}

	// Run migrations first (tables may not exist yet on first run).
	if err := database.Migrate(ctx, db); err != nil {
		db.Close()
		t.Fatalf("migrate: %v", err)
	}

	// Truncate all tables in FK-safe order, then re-populate schema_migrations.
	_, _ = db.ExecContext(ctx,
		`TRUNCATE cards, board_columns, boards, oauth_accounts, sessions, users, schema_migrations CASCADE`)
	if err := database.Migrate(ctx, db); err != nil {
		db.Close()
		t.Fatalf("re-migrate: %v", err)
	}

	t.Cleanup(func() { db.Close() })
	return db
}

func replaceDSNDatabase(dsn, dbName string) string {
	u, err := url.Parse(dsn)
	if err != nil {
		return dsn
	}
	u.Path = "/" + dbName
	return u.String()
}

func sanitize(s string) string {
	var b strings.Builder
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_' {
			b.WriteRune(c)
		}
	}
	return b.String()
}
