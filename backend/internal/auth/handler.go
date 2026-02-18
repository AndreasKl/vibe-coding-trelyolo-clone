package auth

import (
	"net/http"
	"time"
	"trello-clone/internal/httputil"
)

type Handler struct {
	Store        *Store
	CookieDomain string
}

func (h *Handler) setSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		Domain:   h.CookieDomain,
		MaxAge:   30 * 24 * 60 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handler) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		Domain:   h.CookieDomain,
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Email == "" || req.Password == "" {
		httputil.Error(w, http.StatusBadRequest, "email and password required")
		return
	}
	if len(req.Password) < 8 {
		httputil.Error(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	name := req.Name
	if name == "" {
		name = req.Email
	}

	u, err := h.Store.CreateUser(r.Context(), req.Email, hash, name)
	if err != nil {
		httputil.Error(w, http.StatusConflict, "email already registered")
		return
	}

	sess, err := h.Store.CreateSession(r.Context(), u.ID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.setSessionCookie(w, sess.Token)
	httputil.JSON(w, http.StatusCreated, u)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := httputil.Decode(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	u, err := h.Store.UserByEmail(r.Context(), req.Email)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if !CheckPassword(u.PasswordHash, req.Password) {
		httputil.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	sess, err := h.Store.CreateSession(r.Context(), u.ID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.setSessionCookie(w, sess.Token)
	httputil.JSON(w, http.StatusOK, u)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		h.Store.DeleteSession(r.Context(), cookie.Value)
	}
	h.clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	u := UserFromContext(r.Context())
	if u == nil {
		httputil.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}
	httputil.JSON(w, http.StatusOK, u)
}

// cleanup deletes expired sessions — called periodically
func (h *Handler) CleanupSessions() {
	h.Store.DB.Exec("DELETE FROM sessions WHERE expires_at < $1", time.Now())
}
