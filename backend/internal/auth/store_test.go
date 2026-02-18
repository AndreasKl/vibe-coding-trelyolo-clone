package auth

import (
	"context"
	"database/sql"
	"testing"

	"trello-clone/internal/testutil"
)

func TestCreateUser(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	ctx := context.Background()

	u, err := s.CreateUser(ctx, "alice@example.com", "hash123", "Alice")
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if u.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if u.Email != "alice@example.com" {
		t.Fatalf("email = %q, want alice@example.com", u.Email)
	}
	if u.Name != "Alice" {
		t.Fatalf("name = %q, want Alice", u.Name)
	}
	if u.CreatedAt.IsZero() {
		t.Fatal("expected non-zero created_at")
	}
}

func TestCreateUserDuplicateEmail(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	ctx := context.Background()

	_, err := s.CreateUser(ctx, "dup@example.com", "hash", "A")
	if err != nil {
		t.Fatalf("first create: %v", err)
	}
	_, err = s.CreateUser(ctx, "dup@example.com", "hash2", "B")
	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
}

func TestUserByEmail(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	ctx := context.Background()

	created, _ := s.CreateUser(ctx, "lookup@example.com", "hash", "Lookup")
	found, err := s.UserByEmail(ctx, "lookup@example.com")
	if err != nil {
		t.Fatalf("user by email: %v", err)
	}
	if found.ID != created.ID {
		t.Fatalf("ID mismatch: got %s, want %s", found.ID, created.ID)
	}
}

func TestUserByID(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	ctx := context.Background()

	created, _ := s.CreateUser(ctx, "byid@example.com", "hash", "ByID")
	found, err := s.UserByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("user by id: %v", err)
	}
	if found.Email != "byid@example.com" {
		t.Fatalf("email = %q, want byid@example.com", found.Email)
	}
}

func TestUserByEmailNotFound(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	ctx := context.Background()

	_, err := s.UserByEmail(ctx, "nonexistent@example.com")
	if err == nil {
		t.Fatal("expected error for missing email")
	}
}

func TestCreateSession(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	ctx := context.Background()

	u, _ := s.CreateUser(ctx, "sess@example.com", "hash", "Sess")
	sess, err := s.CreateSession(ctx, u.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	if len(sess.Token) != 64 {
		t.Fatalf("token length = %d, want 64", len(sess.Token))
	}
	if sess.UserID != u.ID {
		t.Fatalf("user_id = %s, want %s", sess.UserID, u.ID)
	}
}

func TestSessionByToken(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	ctx := context.Background()

	u, _ := s.CreateUser(ctx, "tok@example.com", "hash", "Tok")
	created, _ := s.CreateSession(ctx, u.ID)

	found, err := s.SessionByToken(ctx, created.Token)
	if err != nil {
		t.Fatalf("session by token: %v", err)
	}
	if found.UserID != u.ID {
		t.Fatalf("user_id = %s, want %s", found.UserID, u.ID)
	}
}

func TestSessionByTokenExpired(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	ctx := context.Background()

	u, _ := s.CreateUser(ctx, "exp@example.com", "hash", "Exp")
	sess, _ := s.CreateSession(ctx, u.ID)

	// Manually expire the session
	_, err := db.ExecContext(ctx,
		"UPDATE sessions SET expires_at = now() - interval '1 hour' WHERE token=$1",
		sess.Token,
	)
	if err != nil {
		t.Fatalf("expire session: %v", err)
	}

	_, err = s.SessionByToken(ctx, sess.Token)
	if err == nil {
		t.Fatal("expected error for expired session")
	}
}

func TestDeleteSession(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	ctx := context.Background()

	u, _ := s.CreateUser(ctx, "del@example.com", "hash", "Del")
	sess, _ := s.CreateSession(ctx, u.ID)

	if err := s.DeleteSession(ctx, sess.Token); err != nil {
		t.Fatalf("delete session: %v", err)
	}

	_, err := s.SessionByToken(ctx, sess.Token)
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestFindOrCreateOAuthUser_New(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	ctx := context.Background()

	u, err := s.FindOrCreateOAuthUser(ctx, "google", "gid123", "oauth@example.com", "OAuth User")
	if err != nil {
		t.Fatalf("find or create: %v", err)
	}
	if u.Email != "oauth@example.com" {
		t.Fatalf("email = %q, want oauth@example.com", u.Email)
	}
	if u.Name != "OAuth User" {
		t.Fatalf("name = %q, want OAuth User", u.Name)
	}

	// Verify oauth_account was created
	var count int
	err = db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM oauth_accounts WHERE provider='google' AND provider_id='gid123'",
	).Scan(&count)
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 1 {
		t.Fatalf("oauth_accounts count = %d, want 1", count)
	}
}

func TestFindOrCreateOAuthUser_ExistingProvider(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	ctx := context.Background()

	u1, _ := s.FindOrCreateOAuthUser(ctx, "google", "gid456", "existing@example.com", "First")
	u2, err := s.FindOrCreateOAuthUser(ctx, "google", "gid456", "existing@example.com", "Second")
	if err != nil {
		t.Fatalf("second call: %v", err)
	}
	if u2.ID != u1.ID {
		t.Fatalf("expected same user, got %s vs %s", u2.ID, u1.ID)
	}
}

func TestFindOrCreateOAuthUser_ExistingEmail(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	ctx := context.Background()

	// Create a regular user first
	existing, _ := s.CreateUser(ctx, "link@example.com", "hash", "Link")

	// OAuth with same email should link, not create new user
	u, err := s.FindOrCreateOAuthUser(ctx, "microsoft", "msid789", "link@example.com", "Link MS")
	if err != nil {
		t.Fatalf("find or create: %v", err)
	}
	if u.ID != existing.ID {
		t.Fatalf("expected linked user ID %s, got %s", existing.ID, u.ID)
	}

	// Verify oauth_account was linked
	var linkedUserID string
	err = db.QueryRowContext(ctx,
		"SELECT user_id FROM oauth_accounts WHERE provider='microsoft' AND provider_id='msid789'",
	).Scan(&linkedUserID)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if linkedUserID != existing.ID {
		t.Fatalf("linked user_id = %s, want %s", linkedUserID, existing.ID)
	}
}

// helper to create a user for tests that need one
func createTestUser(t *testing.T, db *sql.DB, email string) *User {
	t.Helper()
	s := &Store{DB: db}
	hash, _ := HashPassword("password123")
	u, err := s.CreateUser(context.Background(), email, hash, "Test User")
	if err != nil {
		t.Fatalf("create test user: %v", err)
	}
	return u
}
