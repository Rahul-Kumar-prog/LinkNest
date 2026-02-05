package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Rahul-Kumar-prog/linknest/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	addCorsHeaders(w, r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid body",
		})
		return
	}

	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid credentials",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid credentials",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"user": map[string]string{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}
