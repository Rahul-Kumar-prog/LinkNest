package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/Rahul-Kumar-prog/linknest/internal/db"
	"golang.org/x/oauth2"
)

var xOAuthConfig *oauth2.Config

type XUserResponse struct {
	Data XUserData `json:"data"`
}

type xTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

type XUserData struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Username        string `json:"username"`
	ProfileImageURL string `json:"profile_image_url"`
}

func InitTwitterOAuth() {
	clientID := os.Getenv("X_CLIENT_ID")
	clientSecret := os.Getenv("X_CLIENT_SECRET")
	redirectURL := os.Getenv("X_REDIRECT_URL")
	if clientID == "" || clientSecret == "" || redirectURL == "" {
		log.Println("X OAuth not configured; skipping X auth initialization")
		return
	}

	xOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"tweet.read", "tweet.write", "users.read", "offline.access"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://twitter.com/i/oauth2/authorize",
			TokenURL:  "https://api.x.com/2/oauth2/token",
			AuthStyle: oauth2.AuthStyleInHeader,
		},
	}
}

func TwitterLoginHandler(w http.ResponseWriter, r *http.Request) {
	if xOAuthConfig == nil {
		writeError(w, http.StatusServiceUnavailable, "X auth is not configured")
		return
	}

	user, err := getAuthenticatedUser(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Login required before connecting X")
		return
	}

	state, err := randomURLSafeString(24)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to initialize X login")
		return
	}

	verifier, err := randomURLSafeString(48)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to initialize X login")
		return
	}

	challengeBytes := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(challengeBytes[:])

	setCookie(w, "x_oauth_state", state, 600)
	setCookie(w, "x_oauth_verifier", verifier, 600)
	setCookie(w, "x_oauth_user_id", user.ID, 600)

	url := xOAuthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func TwitterCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if xOAuthConfig == nil {
		writeError(w, http.StatusServiceUnavailable, "X auth is not configured")
		return
	}

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	stateCookie, stateErr := r.Cookie("x_oauth_state")
	verifierCookie, verifierErr := r.Cookie("x_oauth_verifier")
	userCookie, userErr := r.Cookie("x_oauth_user_id")
	if code == "" || state == "" || stateErr != nil || verifierErr != nil || userErr != nil {
		redirectWithQuery(w, r, "/?oauth=error&platform=x")
		return
	}

	if stateCookie.Value != state {
		redirectWithQuery(w, r, "/?oauth=error&platform=x")
		return
	}

	token, err := exchangeXToken(code, verifierCookie.Value)
	if err != nil {
		log.Printf("x token exchange failed: %v", err)
		redirectWithQuery(w, r, "/?oauth=error&platform=x")
		return
	}

	req, err := http.NewRequest(http.MethodGet, "https://api.x.com/2/users/me?user.fields=profile_image_url,name,username", nil)
	if err != nil {
		log.Printf("failed to build X user info request: %v", err)
		redirectWithQuery(w, r, "/?oauth=error&platform=x")
		return
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("failed to fetch X user info: %v", err)
		redirectWithQuery(w, r, "/?oauth=error&platform=x")
		return
	}
	defer resp.Body.Close()

	var xUserResp XUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&xUserResp); err != nil {
		log.Printf("failed to decode X user response: %v", err)
		redirectWithQuery(w, r, "/?oauth=error&platform=x")
		return
	}

	if err := db.UpdateUserTwitterInfo(
		userCookie.Value,
		token.AccessToken,
		token.RefreshToken,
		xUserResp.Data.ProfileImageURL,
		xUserResp.Data.ID,
		xUserResp.Data.Name,
		xUserResp.Data.Username,
	); err != nil {
		log.Printf("failed to update X info: %v", err)
		redirectWithQuery(w, r, "/?oauth=error&platform=x")
		return
	}

	setCookie(w, "x_oauth_state", "", -1)
	setCookie(w, "x_oauth_verifier", "", -1)
	setCookie(w, "x_oauth_user_id", "", -1)

	redirectWithQuery(w, r, "/?oauth=success&platform=x")
}

func exchangeXToken(code, verifier string) (*oauth2.Token, error) {
	form := url.Values{}
	form.Set("code", code)
	form.Set("grant_type", "authorization_code")
	form.Set("redirect_uri", xOAuthConfig.RedirectURL)
	form.Set("code_verifier", verifier)

	req, err := http.NewRequest(
		http.MethodPost,
		xOAuthConfig.Endpoint.TokenURL,
		bytes.NewBufferString(form.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(xOAuthConfig.ClientID, xOAuthConfig.ClientSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, loggableError(body, resp.Status)
	}

	var tokenResp xTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &oauth2.Token{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		TokenType:    tokenResp.TokenType,
	}, nil
}

func loggableError(body []byte, fallback string) error {
	message := string(body)
	if message == "" {
		message = fallback
	}

	return &xAPIError{message: message}
}

type xAPIError struct {
	message string
}

func (e *xAPIError) Error() string {
	return e.message
}

func GetTwitterProfileHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getAuthenticatedUser(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	if user.XProfilePicture == "" {
		writeJSON(w, http.StatusOK, map[string]any{
			"connected": false,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"connected": true,
		"user": map[string]any{
			"id":              user.XUserID,
			"profile_picture": user.XProfilePicture,
		},
	})
}
