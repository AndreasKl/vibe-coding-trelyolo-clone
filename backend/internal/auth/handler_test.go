package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"trello-clone/internal/testutil"
)

func signupRequest(t *testing.T, handler http.HandlerFunc, body string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(http.MethodPost, "/api/auth/signup", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler(w, r)
	return w
}

func TestSignup(t *testing.T) {
	db := testutil.SetupDB(t)
	h := &Handler{Store: &Store{DB: db}}

	w := signupRequest(t, h.Signup, `{"email":"new@example.com","password":"password123","name":"New"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body = %s", w.Code, http.StatusCreated, w.Body.String())
	}

	// Check cookie
	cookies := w.Result().Cookies()
	var found bool
	for _, c := range cookies {
		if c.Name == "session" && c.Value != "" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected session cookie to be set")
	}

	// Check user JSON
	var user map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if user["email"] != "new@example.com" {
		t.Fatalf("email = %v, want new@example.com", user["email"])
	}
	// password_hash should not be in JSON (json:"-")
	if _, ok := user["password_hash"]; ok {
		t.Fatal("password_hash should not be in response")
	}
}

func TestSignupMissingFields(t *testing.T) {
	db := testutil.SetupDB(t)
	h := &Handler{Store: &Store{DB: db}}

	w := signupRequest(t, h.Signup, `{}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestSignupShortPassword(t *testing.T) {
	db := testutil.SetupDB(t)
	h := &Handler{Store: &Store{DB: db}}

	w := signupRequest(t, h.Signup, `{"email":"short@example.com","password":"1234567"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestSignupDuplicateEmail(t *testing.T) {
	db := testutil.SetupDB(t)
	h := &Handler{Store: &Store{DB: db}}

	signupRequest(t, h.Signup, `{"email":"dup@example.com","password":"password123"}`)
	w := signupRequest(t, h.Signup, `{"email":"dup@example.com","password":"password456"}`)
	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusConflict)
	}
}

func TestLogin(t *testing.T) {
	db := testutil.SetupDB(t)
	h := &Handler{Store: &Store{DB: db}}

	// Signup first
	signupRequest(t, h.Signup, `{"email":"login@example.com","password":"password123"}`)

	// Login
	r := httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"email":"login@example.com","password":"password123"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Login(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", w.Code, http.StatusOK, w.Body.String())
	}

	var found bool
	for _, c := range w.Result().Cookies() {
		if c.Name == "session" && c.Value != "" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected session cookie")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	db := testutil.SetupDB(t)
	h := &Handler{Store: &Store{DB: db}}

	signupRequest(t, h.Signup, `{"email":"wrong@example.com","password":"password123"}`)

	r := httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"email":"wrong@example.com","password":"wrongpassword"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Login(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestLoginUnknownEmail(t *testing.T) {
	db := testutil.SetupDB(t)
	h := &Handler{Store: &Store{DB: db}}

	r := httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"email":"unknown@example.com","password":"password123"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Login(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestLogout(t *testing.T) {
	db := testutil.SetupDB(t)
	h := &Handler{Store: &Store{DB: db}}

	// Signup to get a session cookie
	sw := signupRequest(t, h.Signup, `{"email":"logout@example.com","password":"password123"}`)
	var sessionCookie *http.Cookie
	for _, c := range sw.Result().Cookies() {
		if c.Name == "session" {
			sessionCookie = c
		}
	}

	// Logout
	r := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	r.AddCookie(sessionCookie)
	w := httptest.NewRecorder()
	h.Logout(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}

	// Check cookie is cleared
	for _, c := range w.Result().Cookies() {
		if c.Name == "session" && c.MaxAge != -1 {
			t.Fatalf("expected MaxAge=-1, got %d", c.MaxAge)
		}
	}
}

func TestMe(t *testing.T) {
	db := testutil.SetupDB(t)
	store := &Store{DB: db}
	h := &Handler{Store: store}

	// Signup
	sw := signupRequest(t, h.Signup, `{"email":"me@example.com","password":"password123","name":"Me User"}`)
	var sessionCookie *http.Cookie
	for _, c := range sw.Result().Cookies() {
		if c.Name == "session" {
			sessionCookie = c
		}
	}

	// Wrap through RequireAuth middleware
	middleware := RequireAuth(store)
	handler := middleware(http.HandlerFunc(h.Me))

	r := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	r.AddCookie(sessionCookie)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", w.Code, http.StatusOK, w.Body.String())
	}

	var user map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if user["email"] != "me@example.com" {
		t.Fatalf("email = %v, want me@example.com", user["email"])
	}
}
