package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Rahul-Kumar-prog/linknest/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if user.Password == "" {
		writeError(w, http.StatusUnauthorized, "Use Google sign in for this account")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := issueAuthSession(w, user)
	if err != nil {
		log.Print("failed to generate jwt token: ", err)
		writeError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Login successful",
		"token":   token,
		"user": map[string]string{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	clearAuthSession(w)
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Logout successful",
	})
}
