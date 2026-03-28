package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"

	"github.com/Rahul-Kumar-prog/linknest/internal/db"
	utils "github.com/Rahul-Kumar-prog/linknest/utils/jwt"
)

func frontendURL() string {
	url := os.Getenv("FRONTEND_URL")
	if url == "" {
		return "http://localhost:5173"
	}

	return url
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]string{
		"error": message,
	})
}

func issueAuthSession(w http.ResponseWriter, user *db.User) (string, error) {
	token, err := utils.GenerateJwtToken(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		MaxAge:   60 * 60 * 24,
		SameSite: http.SameSiteLaxMode,
	})

	return token, nil
}

func setCookie(w http.ResponseWriter, name, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		MaxAge:   maxAge,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearAuthSession(w http.ResponseWriter) {
	setCookie(w, "auth_token", "", -1)
}

func getAuthenticatedUser(r *http.Request) (*db.User, error) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return nil, err
	}

	claims, err := utils.ValidateJwtToken(cookie.Value)
	if err != nil {
		return nil, err
	}

	return db.GetUserByID(claims.UserID)
}

func redirectWithQuery(w http.ResponseWriter, r *http.Request, path string) {
	http.Redirect(w, r, frontendURL()+path, http.StatusTemporaryRedirect)
}

func randomURLSafeString(size int) (string, error) {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
