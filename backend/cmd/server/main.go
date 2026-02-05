package main

import (
	"log"
	"net/http"

	"github.com/Rahul-Kumar-prog/linknest/internal/db"
	"github.com/Rahul-Kumar-prog/linknest/internal/handlers"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	if err := db.Connect(); err != nil {
		log.Fatal("DB connection failed", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/login", handlers.LoginHandler)
	mux.HandleFunc("/api/signup", handlers.SignupHandler)

	corsHandler := corsMiddleware(mux)

	log.Println("🚀 Backend server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
