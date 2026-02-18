package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"trello-clone/internal/testutil"
)

func TestRequireAuthNoCookie(t *testing.T) {
	db := testutil.SetupDB(t)
	store := &Store{DB: db}
	middleware := RequireAuth(store)

	called := false
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
	if called {
		t.Fatal("next handler should not have been called")
	}
}

func TestRequireAuthInvalidToken(t *testing.T) {
	db := testutil.SetupDB(t)
	store := &Store{DB: db}
	middleware := RequireAuth(store)

	called := false
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{Name: "session", Value: "invalidtoken"})
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
	if called {
		t.Fatal("next handler should not have been called")
	}
}

func TestRequireAuthValid(t *testing.T) {
	db := testutil.SetupDB(t)
	store := &Store{DB: db}
	ctx := context.Background()

	hash, _ := HashPassword("password123")
	u, _ := store.CreateUser(ctx, "valid@example.com", hash, "Valid")
	sess, _ := store.CreateSession(ctx, u.ID)

	middleware := RequireAuth(store)

	var gotUser *User
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUser = UserFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{Name: "session", Value: sess.Token})
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if gotUser == nil {
		t.Fatal("expected user in context")
	}
	if gotUser.ID != u.ID {
		t.Fatalf("user ID = %s, want %s", gotUser.ID, u.ID)
	}
}

func TestUserFromContextNil(t *testing.T) {
	ctx := context.Background()
	u := UserFromContext(ctx)
	if u != nil {
		t.Fatal("expected nil user from empty context")
	}
}
