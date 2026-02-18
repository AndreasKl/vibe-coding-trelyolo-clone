package auth

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Name         string `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
}

type Session struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
}

type Store struct {
	DB *sql.DB
}

func (s *Store) CreateUser(ctx context.Context, email, passwordHash, name string) (*User, error) {
	u := &User{}
	err := s.DB.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash, name) VALUES ($1, $2, $3)
		 RETURNING id, email, password_hash, name, created_at`,
		email, passwordHash, name,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Store) UserByEmail(ctx context.Context, email string) (*User, error) {
	u := &User{}
	err := s.DB.QueryRowContext(ctx,
		`SELECT id, email, password_hash, name, created_at FROM users WHERE email=$1`, email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Store) UserByID(ctx context.Context, id string) (*User, error) {
	u := &User{}
	err := s.DB.QueryRowContext(ctx,
		`SELECT id, email, password_hash, name, created_at FROM users WHERE id=$1`, id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Store) CreateSession(ctx context.Context, userID string) (*Session, error) {
	token, err := GenerateToken()
	if err != nil {
		return nil, err
	}
	sess := &Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}
	_, err = s.DB.ExecContext(ctx,
		`INSERT INTO sessions (token, user_id, expires_at) VALUES ($1, $2, $3)`,
		sess.Token, sess.UserID, sess.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *Store) SessionByToken(ctx context.Context, token string) (*Session, error) {
	sess := &Session{}
	err := s.DB.QueryRowContext(ctx,
		`SELECT token, user_id, expires_at FROM sessions WHERE token=$1 AND expires_at > now()`, token,
	).Scan(&sess.Token, &sess.UserID, &sess.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *Store) DeleteSession(ctx context.Context, token string) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM sessions WHERE token=$1`, token)
	return err
}

func (s *Store) FindOrCreateOAuthUser(ctx context.Context, provider, providerID, email, name string) (*User, error) {
	// Try to find by provider+providerID
	var userID string
	err := s.DB.QueryRowContext(ctx,
		`SELECT user_id FROM oauth_accounts WHERE provider=$1 AND provider_id=$2`,
		provider, providerID,
	).Scan(&userID)
	if err == nil {
		return s.UserByID(ctx, userID)
	}

	// Try to find by email
	u, err := s.UserByEmail(ctx, email)
	if err == nil {
		// Link OAuth account
		_, err = s.DB.ExecContext(ctx,
			`INSERT INTO oauth_accounts (user_id, provider, provider_id) VALUES ($1, $2, $3)`,
			u.ID, provider, providerID,
		)
		return u, err
	}

	// Create new user + OAuth account
	u, err = s.CreateUser(ctx, email, "", name)
	if err != nil {
		return nil, err
	}
	_, err = s.DB.ExecContext(ctx,
		`INSERT INTO oauth_accounts (user_id, provider, provider_id) VALUES ($1, $2, $3)`,
		u.ID, provider, providerID,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}
