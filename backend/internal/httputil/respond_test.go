package httputil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	JSON(w, http.StatusCreated, map[string]string{"hello": "world"})

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusCreated)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("content-type = %q, want application/json", ct)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["hello"] != "world" {
		t.Fatalf("body = %v, want hello=world", body)
	}
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	Error(w, http.StatusBadRequest, "something went wrong")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["error"] != "something went wrong" {
		t.Fatalf("error = %q, want %q", body["error"], "something went wrong")
	}
}

func TestDecode(t *testing.T) {
	body := `{"name":"alice","age":30}`
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	var v struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	if err := Decode(r, &v); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if v.Name != "alice" {
		t.Fatalf("name = %q, want alice", v.Name)
	}
	if v.Age != 30 {
		t.Fatalf("age = %d, want 30", v.Age)
	}
}
