package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/Rahul-Kumar-prog/linknest/internal/db"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAuthConfig *oauth2.Config

type GoogleUserInfo struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func InitGoogleOAuth() {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if clientID == "" || clientSecret == "" || redirectURL == "" {
		log.Println("Google OAuth not configured; skipping Google auth initialization")
		return
	}

	googleOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	if googleOAuthConfig == nil {
		writeError(w, http.StatusServiceUnavailable, "Google auth is not configured")
		return
	}

	url := googleOAuthConfig.AuthCodeURL("google-login", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if googleOAuthConfig == nil {
		writeError(w, http.StatusServiceUnavailable, "Google auth is not configured")
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		redirectWithQuery(w, r, "/?auth=error&provider=google")
		return
	}

	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("google token exchange failed: %v", err)
		redirectWithQuery(w, r, "/?auth=error&provider=google")
		return
	}

	client := googleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")
	if err != nil {
		log.Printf("google userinfo failed: %v", err)
		redirectWithQuery(w, r, "/?auth=error&provider=google")
		return
	}
	defer resp.Body.Close()

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Printf("google userinfo decode failed: %v", err)
		redirectWithQuery(w, r, "/?auth=error&provider=google")
		return
	}

	user, err := db.UpsertGoogleUser(userInfo.Sub, userInfo.Email, userInfo.Name, userInfo.Picture)
	if err != nil {
		log.Printf("google upsert user failed: %v", err)
		redirectWithQuery(w, r, "/?auth=error&provider=google")
		return
	}

	if _, err := issueAuthSession(w, user); err != nil {
		log.Printf("google issue session failed: %v", err)
		redirectWithQuery(w, r, "/?auth=error&provider=google")
		return
	}

	redirectWithQuery(w, r, "/?auth=success&provider=google")
}
