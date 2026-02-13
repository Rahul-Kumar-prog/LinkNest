package main

import (
	"log"
	"net/http"

	"github.com/Rahul-Kumar-prog/linknest/internal/db"
	"github.com/Rahul-Kumar-prog/linknest/internal/handlers"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000") // or 5173

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: Could not load .env file:", err)
	}
	if err := db.Connect(); err != nil {
		log.Fatal("DB connection failed", err)
	}
	println("DB Connected Successfully")

	mux := http.NewServeMux()

	mux.HandleFunc("/api/login", handlers.LoginHandler)
	mux.HandleFunc("/api/signup", handlers.SignupHandler)
	mux.HandleFunc("/api/twitter", handlers.TwitterHandler)

	corsHandler := corsMiddleware(mux)

	log.Println("🚀 Backend server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
