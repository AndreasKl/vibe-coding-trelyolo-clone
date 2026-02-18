package auth

import (
	"context"
	"net/http"
	"trello-clone/internal/httputil"
)

type contextKey string

const userKey contextKey = "user"

func UserFromContext(ctx context.Context) *User {
	u, _ := ctx.Value(userKey).(*User)
	return u
}

func RequireAuth(store *Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session")
			if err != nil {
				httputil.Error(w, http.StatusUnauthorized, "not authenticated")
				return
			}
			sess, err := store.SessionByToken(r.Context(), cookie.Value)
			if err != nil {
				httputil.Error(w, http.StatusUnauthorized, "invalid session")
				return
			}
			u, err := store.UserByID(r.Context(), sess.UserID)
			if err != nil {
				httputil.Error(w, http.StatusUnauthorized, "user not found")
				return
			}
			ctx := context.WithValue(r.Context(), userKey, u)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
