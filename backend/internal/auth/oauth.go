package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	microsoftEndpoint "golang.org/x/oauth2/microsoft"
)

type OAuthConfig struct {
	Google    *oauth2.Config
	Microsoft *oauth2.Config
}

func NewOAuthConfig(baseURL, googleID, googleSecret, msID, msSecret string) *OAuthConfig {
	cfg := &OAuthConfig{}
	if googleID != "" {
		cfg.Google = &oauth2.Config{
			ClientID:     googleID,
			ClientSecret: googleSecret,
			RedirectURL:  baseURL + "/api/auth/oauth/google/callback",
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		}
	}
	if msID != "" {
		cfg.Microsoft = &oauth2.Config{
			ClientID:     msID,
			ClientSecret: msSecret,
			RedirectURL:  baseURL + "/api/auth/oauth/microsoft/callback",
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     microsoftEndpoint.AzureADEndpoint("common"),
		}
	}
	return cfg
}

type OAuthHandler struct {
	Config       *OAuthConfig
	Store        *Store
	CookieDomain string
	BaseURL      string
}

func (h *OAuthHandler) providerConfig(provider string) *oauth2.Config {
	switch provider {
	case "google":
		return h.Config.Google
	case "microsoft":
		return h.Config.Microsoft
	}
	return nil
}

func (h *OAuthHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	cfg := h.providerConfig(provider)
	if cfg == nil {
		http.Error(w, "unsupported provider", http.StatusBadRequest)
		return
	}

	state, _ := generateState()
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   600,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, cfg.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	cfg := h.providerConfig(provider)
	if cfg == nil {
		http.Error(w, "unsupported provider", http.StatusBadRequest)
		return
	}

	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	token, err := cfg.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "oauth exchange failed", http.StatusBadRequest)
		return
	}

	info, err := fetchUserInfo(r.Context(), provider, cfg, token)
	if err != nil {
		http.Error(w, "failed to get user info", http.StatusInternalServerError)
		return
	}

	u, err := h.Store.FindOrCreateOAuthUser(r.Context(), provider, info.ID, info.Email, info.Name)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	sess, err := h.Store.CreateSession(r.Context(), u.ID)
	if err != nil {
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sess.Token,
		Path:     "/",
		Domain:   h.CookieDomain,
		MaxAge:   30 * 24 * 60 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, h.BaseURL, http.StatusTemporaryRedirect)
}

type userInfo struct {
	ID    string
	Email string
	Name  string
}

func fetchUserInfo(ctx context.Context, provider string, cfg *oauth2.Config, token *oauth2.Token) (*userInfo, error) {
	client := cfg.Client(ctx, token)

	var url string
	switch provider {
	case "google":
		url = "https://www.googleapis.com/oauth2/v2/userinfo"
	case "microsoft":
		url = "https://graph.microsoft.com/v1.0/me"
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	info := &userInfo{}
	switch provider {
	case "google":
		info.ID = fmt.Sprint(data["id"])
		info.Email, _ = data["email"].(string)
		info.Name, _ = data["name"].(string)
	case "microsoft":
		info.ID = fmt.Sprint(data["id"])
		info.Email, _ = data["mail"].(string)
		if info.Email == "" {
			info.Email, _ = data["userPrincipalName"].(string)
		}
		info.Name, _ = data["displayName"].(string)
	}

	if info.Email == "" {
		return nil, fmt.Errorf("no email from %s", provider)
	}
	return info, nil
}

func generateState() (string, error) {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b), nil
}
