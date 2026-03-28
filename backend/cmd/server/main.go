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
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = handlersOrigin()
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", origin)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handlersOrigin() string {
	return "http://localhost:5173"
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: Could not load .env file:", err)
	}

	if err := db.Connect(); err != nil {
		log.Fatal("DB connection failed", err)
	}

	handlers.InitGoogleOAuth()
	handlers.InitTwitterOAuth()
	handlers.InitLinkedInOAuth()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/login", handlers.LoginHandler)
	mux.HandleFunc("/api/logout", handlers.LogoutHandler)
	mux.HandleFunc("/api/signup", handlers.SignupHandler)
	mux.HandleFunc("/api/auth/google", handlers.GoogleLoginHandler)
	mux.HandleFunc("/api/auth/google/callback", handlers.GoogleCallbackHandler)
	mux.HandleFunc("/api/twitter", handlers.TwitterLoginHandler)
	mux.HandleFunc("/api/twitter/callback", handlers.TwitterCallbackHandler)
	mux.HandleFunc("/api/twitter/profile", handlers.GetTwitterProfileHandler)
	mux.HandleFunc("/api/linkedin", handlers.LinkedInLoginHandler)
	mux.HandleFunc("/api/linkedin/callback", handlers.LinkedInCallbackHandler)
	mux.HandleFunc("/api/platforms/disconnect", handlers.DisconnectPlatformHandler)
	mux.HandleFunc("/api/me/status", handlers.GetConnectionStatusHandler)
	mux.HandleFunc("/api/posts/publish", handlers.PublishPostHandler)

	corsHandler := corsMiddleware(mux)

	log.Println("Backend server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
