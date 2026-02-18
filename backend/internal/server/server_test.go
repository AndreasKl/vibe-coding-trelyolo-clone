package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"trello-clone/internal/testutil"
)

func TestHealthEndpoint(t *testing.T) {
	db := testutil.SetupDB(t)
	srv := New(Config{DB: db})

	r := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()
	srv.Handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("status = %q, want ok", body["status"])
	}
}

func TestCORSHeaders(t *testing.T) {
	db := testutil.SetupDB(t)
	srv := New(Config{DB: db, AllowOrigin: "http://localhost:5173"})

	r := httptest.NewRequest(http.MethodOptions, "/api/boards", nil)
	w := httptest.NewRecorder()
	srv.Handler.ServeHTTP(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
	if acao := w.Header().Get("Access-Control-Allow-Origin"); acao != "http://localhost:5173" {
		t.Fatalf("ACAO = %q, want http://localhost:5173", acao)
	}
	if acam := w.Header().Get("Access-Control-Allow-Methods"); acam == "" {
		t.Fatal("expected Access-Control-Allow-Methods header")
	}
	if acah := w.Header().Get("Access-Control-Allow-Headers"); acah == "" {
		t.Fatal("expected Access-Control-Allow-Headers header")
	}
	if acac := w.Header().Get("Access-Control-Allow-Credentials"); acac != "true" {
		t.Fatalf("ACAC = %q, want true", acac)
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	// Test the recovery middleware directly
	panicking := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	handler := recovery(panicking)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Should not panic
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
