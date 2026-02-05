package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Rahul-Kumar-prog/linknest/internal/db"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func addCorsHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	addCorsHeaders(w, r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	err := db.CreateUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, "User already exists or db error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
	})
}
