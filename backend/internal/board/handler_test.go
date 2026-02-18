package board_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"trello-clone/internal/server"
	"trello-clone/internal/testutil"
)

func signupAndGetCookie(t *testing.T, srv *http.Server) *http.Cookie {
	t.Helper()
	body := `{"email":"handler@example.com","password":"password123","name":"Handler User"}`
	r := httptest.NewRequest(http.MethodPost, "/api/auth/signup", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.Handler.ServeHTTP(w, r)
	if w.Code != http.StatusCreated {
		t.Fatalf("signup: status = %d, body = %s", w.Code, w.Body.String())
	}
	for _, c := range w.Result().Cookies() {
		if c.Name == "session" {
			return c
		}
	}
	t.Fatal("no session cookie after signup")
	return nil
}

func doRequest(t *testing.T, srv *http.Server, method, path string, body string, cookie *http.Cookie) *httptest.ResponseRecorder {
	t.Helper()
	var r *http.Request
	if body != "" {
		r = httptest.NewRequest(method, path, strings.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
	} else {
		r = httptest.NewRequest(method, path, nil)
	}
	if cookie != nil {
		r.AddCookie(cookie)
	}
	w := httptest.NewRecorder()
	srv.Handler.ServeHTTP(w, r)
	return w
}

func TestListBoardsHandler(t *testing.T) {
	db := testutil.SetupDB(t)
	srv := server.New(server.Config{DB: db})
	cookie := signupAndGetCookie(t, srv)

	w := doRequest(t, srv, http.MethodGet, "/api/boards", "", cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var boards []any
	json.Unmarshal(w.Body.Bytes(), &boards)
	if len(boards) != 0 {
		t.Fatalf("boards = %d, want 0", len(boards))
	}
}

func TestCreateBoardHandler(t *testing.T) {
	db := testutil.SetupDB(t)
	srv := server.New(server.Config{DB: db})
	cookie := signupAndGetCookie(t, srv)

	w := doRequest(t, srv, http.MethodPost, "/api/boards", `{"name":"Test Board"}`, cookie)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body = %s", w.Code, http.StatusCreated, w.Body.String())
	}

	var board map[string]any
	json.Unmarshal(w.Body.Bytes(), &board)
	if board["name"] != "Test Board" {
		t.Fatalf("name = %v, want Test Board", board["name"])
	}
}

func TestGetBoardHandler(t *testing.T) {
	db := testutil.SetupDB(t)
	srv := server.New(server.Config{DB: db})
	cookie := signupAndGetCookie(t, srv)

	// Create board
	cw := doRequest(t, srv, http.MethodPost, "/api/boards", `{"name":"Get Board"}`, cookie)
	var created map[string]any
	json.Unmarshal(cw.Body.Bytes(), &created)
	boardID := created["id"].(string)

	// Get board
	w := doRequest(t, srv, http.MethodGet, "/api/boards/"+boardID, "", cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var board map[string]any
	json.Unmarshal(w.Body.Bytes(), &board)
	columns, ok := board["columns"].([]any)
	if !ok {
		t.Fatal("expected columns array")
	}
	if len(columns) != 3 {
		t.Fatalf("columns = %d, want 3", len(columns))
	}
}

func TestDeleteBoardHandler(t *testing.T) {
	db := testutil.SetupDB(t)
	srv := server.New(server.Config{DB: db})
	cookie := signupAndGetCookie(t, srv)

	cw := doRequest(t, srv, http.MethodPost, "/api/boards", `{"name":"Delete Board"}`, cookie)
	var created map[string]any
	json.Unmarshal(cw.Body.Bytes(), &created)
	boardID := created["id"].(string)

	w := doRequest(t, srv, http.MethodDelete, "/api/boards/"+boardID, "", cookie)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestCreateColumnHandler(t *testing.T) {
	db := testutil.SetupDB(t)
	srv := server.New(server.Config{DB: db})
	cookie := signupAndGetCookie(t, srv)

	cw := doRequest(t, srv, http.MethodPost, "/api/boards", `{"name":"Col Board"}`, cookie)
	var board map[string]any
	json.Unmarshal(cw.Body.Bytes(), &board)
	boardID := board["id"].(string)

	w := doRequest(t, srv, http.MethodPost, fmt.Sprintf("/api/boards/%s/columns", boardID), `{"name":"New Column"}`, cookie)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body = %s", w.Code, http.StatusCreated, w.Body.String())
	}

	var col map[string]any
	json.Unmarshal(w.Body.Bytes(), &col)
	if col["name"] != "New Column" {
		t.Fatalf("name = %v, want New Column", col["name"])
	}
}

func TestUpdateColumnHandler(t *testing.T) {
	db := testutil.SetupDB(t)
	srv := server.New(server.Config{DB: db})
	cookie := signupAndGetCookie(t, srv)

	cw := doRequest(t, srv, http.MethodPost, "/api/boards", `{"name":"Upd Col Board"}`, cookie)
	var board map[string]any
	json.Unmarshal(cw.Body.Bytes(), &board)
	boardID := board["id"].(string)

	// Get board to find a column ID
	gw := doRequest(t, srv, http.MethodGet, "/api/boards/"+boardID, "", cookie)
	var fullBoard map[string]any
	json.Unmarshal(gw.Body.Bytes(), &fullBoard)
	columns := fullBoard["columns"].([]any)
	colID := columns[0].(map[string]any)["id"].(string)

	w := doRequest(t, srv, http.MethodPatch, "/api/columns/"+colID, `{"name":"Renamed"}`, cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", w.Code, http.StatusOK, w.Body.String())
	}

	var col map[string]any
	json.Unmarshal(w.Body.Bytes(), &col)
	if col["name"] != "Renamed" {
		t.Fatalf("name = %v, want Renamed", col["name"])
	}
}

func TestDeleteColumnHandler(t *testing.T) {
	db := testutil.SetupDB(t)
	srv := server.New(server.Config{DB: db})
	cookie := signupAndGetCookie(t, srv)

	cw := doRequest(t, srv, http.MethodPost, "/api/boards", `{"name":"Del Col Board"}`, cookie)
	var board map[string]any
	json.Unmarshal(cw.Body.Bytes(), &board)
	boardID := board["id"].(string)

	gw := doRequest(t, srv, http.MethodGet, "/api/boards/"+boardID, "", cookie)
	var fullBoard map[string]any
	json.Unmarshal(gw.Body.Bytes(), &fullBoard)
	columns := fullBoard["columns"].([]any)
	colID := columns[0].(map[string]any)["id"].(string)

	w := doRequest(t, srv, http.MethodDelete, "/api/columns/"+colID, "", cookie)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestCreateCardHandler(t *testing.T) {
	db := testutil.SetupDB(t)
	srv := server.New(server.Config{DB: db})
	cookie := signupAndGetCookie(t, srv)

	cw := doRequest(t, srv, http.MethodPost, "/api/boards", `{"name":"Card Board"}`, cookie)
	var board map[string]any
	json.Unmarshal(cw.Body.Bytes(), &board)
	boardID := board["id"].(string)

	gw := doRequest(t, srv, http.MethodGet, "/api/boards/"+boardID, "", cookie)
	var fullBoard map[string]any
	json.Unmarshal(gw.Body.Bytes(), &fullBoard)
	columns := fullBoard["columns"].([]any)
	colID := columns[0].(map[string]any)["id"].(string)

	w := doRequest(t, srv, http.MethodPost, fmt.Sprintf("/api/columns/%s/cards", colID), `{"title":"New Card"}`, cookie)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body = %s", w.Code, http.StatusCreated, w.Body.String())
	}

	var card map[string]any
	json.Unmarshal(w.Body.Bytes(), &card)
	if card["title"] != "New Card" {
		t.Fatalf("title = %v, want New Card", card["title"])
	}
}

func TestUpdateCardHandler(t *testing.T) {
	db := testutil.SetupDB(t)
	srv := server.New(server.Config{DB: db})
	cookie := signupAndGetCookie(t, srv)

	cw := doRequest(t, srv, http.MethodPost, "/api/boards", `{"name":"Upd Card Board"}`, cookie)
	var board map[string]any
	json.Unmarshal(cw.Body.Bytes(), &board)
	boardID := board["id"].(string)

	gw := doRequest(t, srv, http.MethodGet, "/api/boards/"+boardID, "", cookie)
	var fullBoard map[string]any
	json.Unmarshal(gw.Body.Bytes(), &fullBoard)
	columns := fullBoard["columns"].([]any)
	colID := columns[0].(map[string]any)["id"].(string)

	// Create card
	ccw := doRequest(t, srv, http.MethodPost, fmt.Sprintf("/api/columns/%s/cards", colID), `{"title":"Original"}`, cookie)
	var card map[string]any
	json.Unmarshal(ccw.Body.Bytes(), &card)
	cardID := card["id"].(string)

	w := doRequest(t, srv, http.MethodPatch, "/api/cards/"+cardID, `{"title":"Updated Title"}`, cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", w.Code, http.StatusOK, w.Body.String())
	}

	var updated map[string]any
	json.Unmarshal(w.Body.Bytes(), &updated)
	if updated["title"] != "Updated Title" {
		t.Fatalf("title = %v, want Updated Title", updated["title"])
	}
}

func TestDeleteCardHandler(t *testing.T) {
	db := testutil.SetupDB(t)
	srv := server.New(server.Config{DB: db})
	cookie := signupAndGetCookie(t, srv)

	cw := doRequest(t, srv, http.MethodPost, "/api/boards", `{"name":"Del Card Board"}`, cookie)
	var board map[string]any
	json.Unmarshal(cw.Body.Bytes(), &board)
	boardID := board["id"].(string)

	gw := doRequest(t, srv, http.MethodGet, "/api/boards/"+boardID, "", cookie)
	var fullBoard map[string]any
	json.Unmarshal(gw.Body.Bytes(), &fullBoard)
	columns := fullBoard["columns"].([]any)
	colID := columns[0].(map[string]any)["id"].(string)

	ccw := doRequest(t, srv, http.MethodPost, fmt.Sprintf("/api/columns/%s/cards", colID), `{"title":"To Delete"}`, cookie)
	var card map[string]any
	json.Unmarshal(ccw.Body.Bytes(), &card)
	cardID := card["id"].(string)

	w := doRequest(t, srv, http.MethodDelete, "/api/cards/"+cardID, "", cookie)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestMoveCardHandler(t *testing.T) {
	db := testutil.SetupDB(t)
	srv := server.New(server.Config{DB: db})
	cookie := signupAndGetCookie(t, srv)

	cw := doRequest(t, srv, http.MethodPost, "/api/boards", `{"name":"Move Card Board"}`, cookie)
	var board map[string]any
	json.Unmarshal(cw.Body.Bytes(), &board)
	boardID := board["id"].(string)

	gw := doRequest(t, srv, http.MethodGet, "/api/boards/"+boardID, "", cookie)
	var fullBoard map[string]any
	json.Unmarshal(gw.Body.Bytes(), &fullBoard)
	columns := fullBoard["columns"].([]any)
	srcColID := columns[0].(map[string]any)["id"].(string)
	dstColID := columns[1].(map[string]any)["id"].(string)

	// Create card in source column
	ccw := doRequest(t, srv, http.MethodPost, fmt.Sprintf("/api/columns/%s/cards", srcColID), `{"title":"Moving Card"}`, cookie)
	var card map[string]any
	json.Unmarshal(ccw.Body.Bytes(), &card)
	cardID := card["id"].(string)

	// Move to destination column
	moveBody := fmt.Sprintf(`{"column_id":"%s","position":0}`, dstColID)
	w := doRequest(t, srv, http.MethodPost, "/api/cards/"+cardID+"/move", moveBody, cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", w.Code, http.StatusOK, w.Body.String())
	}

	var moved map[string]any
	json.Unmarshal(w.Body.Bytes(), &moved)
	if moved["column_id"] != dstColID {
		t.Fatalf("column_id = %v, want %s", moved["column_id"], dstColID)
	}
}

func TestUnauthorized(t *testing.T) {
	db := testutil.SetupDB(t)
	srv := server.New(server.Config{DB: db})

	endpoints := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/boards"},
		{http.MethodPost, "/api/boards"},
		{http.MethodGet, "/api/boards/fake-id"},
		{http.MethodDelete, "/api/boards/fake-id"},
		{http.MethodPost, "/api/boards/fake-id/columns"},
		{http.MethodPatch, "/api/columns/fake-id"},
		{http.MethodDelete, "/api/columns/fake-id"},
		{http.MethodPost, "/api/columns/fake-id/cards"},
		{http.MethodPatch, "/api/cards/fake-id"},
		{http.MethodDelete, "/api/cards/fake-id"},
		{http.MethodPost, "/api/cards/fake-id/move"},
		{http.MethodPost, "/api/auth/logout"},
		{http.MethodGet, "/api/auth/me"},
	}

	for _, ep := range endpoints {
		t.Run(ep.method+" "+ep.path, func(t *testing.T) {
			w := doRequest(t, srv, ep.method, ep.path, "", nil)
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("%s %s: status = %d, want %d", ep.method, ep.path, w.Code, http.StatusUnauthorized)
			}
		})
	}
}
