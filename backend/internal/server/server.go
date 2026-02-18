package server

import (
	"database/sql"
	"net/http"
	"trello-clone/internal/auth"
	"trello-clone/internal/board"
	"trello-clone/internal/httputil"
)

type Config struct {
	DB            *sql.DB
	CookieDomain  string
	AllowOrigin   string
	BaseURL       string
	GoogleID      string
	GoogleSecret  string
	MicrosoftID   string
	MicrosoftSecret string
}

func New(cfg Config) *http.Server {
	authStore := &auth.Store{DB: cfg.DB}
	boardStore := &board.Store{DB: cfg.DB}

	authHandler := &auth.Handler{Store: authStore, CookieDomain: cfg.CookieDomain}
	boardHandler := &board.Handler{Store: boardStore}

	oauthCfg := auth.NewOAuthConfig(cfg.BaseURL, cfg.GoogleID, cfg.GoogleSecret, cfg.MicrosoftID, cfg.MicrosoftSecret)
	oauthHandler := &auth.OAuthHandler{
		Config:       oauthCfg,
		Store:        authStore,
		CookieDomain: cfg.CookieDomain,
		BaseURL:      cfg.AllowOrigin,
	}

	requireAuth := auth.RequireAuth(authStore)

	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Auth (public)
	mux.HandleFunc("POST /api/auth/signup", authHandler.Signup)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("GET /api/auth/oauth/{provider}", oauthHandler.Redirect)
	mux.HandleFunc("GET /api/auth/oauth/{provider}/callback", oauthHandler.Callback)

	// Auth (requires session)
	mux.Handle("POST /api/auth/logout", requireAuth(http.HandlerFunc(authHandler.Logout)))
	mux.Handle("GET /api/auth/me", requireAuth(http.HandlerFunc(authHandler.Me)))

	// Boards
	mux.Handle("GET /api/boards", requireAuth(http.HandlerFunc(boardHandler.ListBoards)))
	mux.Handle("POST /api/boards", requireAuth(http.HandlerFunc(boardHandler.CreateBoard)))
	mux.Handle("GET /api/boards/{id}", requireAuth(http.HandlerFunc(boardHandler.GetBoard)))
	mux.Handle("DELETE /api/boards/{id}", requireAuth(http.HandlerFunc(boardHandler.DeleteBoard)))

	// Columns
	mux.Handle("POST /api/boards/{boardID}/columns", requireAuth(http.HandlerFunc(boardHandler.CreateColumn)))
	mux.Handle("PATCH /api/columns/{id}", requireAuth(http.HandlerFunc(boardHandler.UpdateColumn)))
	mux.Handle("DELETE /api/columns/{id}", requireAuth(http.HandlerFunc(boardHandler.DeleteColumn)))

	// Cards
	mux.Handle("POST /api/columns/{columnID}/cards", requireAuth(http.HandlerFunc(boardHandler.CreateCard)))
	mux.Handle("PATCH /api/cards/{id}", requireAuth(http.HandlerFunc(boardHandler.UpdateCard)))
	mux.Handle("DELETE /api/cards/{id}", requireAuth(http.HandlerFunc(boardHandler.DeleteCard)))
	mux.Handle("POST /api/columns/{id}/move", requireAuth(http.HandlerFunc(boardHandler.MoveColumn)))
	mux.Handle("POST /api/cards/{id}/move", requireAuth(http.HandlerFunc(boardHandler.MoveCard)))

	// Apply middleware
	var handler http.Handler = mux
	handler = cors(cfg.AllowOrigin)(handler)
	handler = logging(handler)
	handler = recovery(handler)

	return &http.Server{Handler: handler}
}
