package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func Connect() error {
	var err error

	host := envOrDefault("DB_HOST", "localhost")
	port := envOrDefault("DB_PORT", "5432")
	user := envOrDefault("DB_USER", "postgres")
	password := envOrDefault("DB_PASS", "12345")
	dbname := envOrDefault("DB_NAME", "linknest")
	sslMode := envOrDefault("DB_SSLMODE", "disable")

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		dbname,
		sslMode,
	)

	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("DB open Error:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("DB ping Error:", err)
	}

	fmt.Println("Database connected successfully")

	if err := runMigrations(); err != nil {
		log.Printf("Migration warning: %v", err)
	}

	return nil
}

func runMigrations() error {
	migrations := []string{
		`CREATE SCHEMA IF NOT EXISTS app`,
		`CREATE EXTENSION IF NOT EXISTS pgcrypto`,
		`CREATE TABLE IF NOT EXISTS app.users (
			id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
			email TEXT NOT NULL UNIQUE,
			password TEXT,
			google_id TEXT,
			google_email TEXT,
			google_name TEXT,
			google_picture TEXT
		)`,
		`ALTER TABLE app.users ADD COLUMN IF NOT EXISTS google_id TEXT`,
		`ALTER TABLE app.users ADD COLUMN IF NOT EXISTS google_email TEXT`,
		`ALTER TABLE app.users ADD COLUMN IF NOT EXISTS google_name TEXT`,
		`ALTER TABLE app.users ADD COLUMN IF NOT EXISTS google_picture TEXT`,
		`ALTER TABLE app.users ADD COLUMN IF NOT EXISTS x_access_token TEXT`,
		`ALTER TABLE app.users ADD COLUMN IF NOT EXISTS x_refresh_token TEXT`,
		`ALTER TABLE app.users ADD COLUMN IF NOT EXISTS x_profile_picture TEXT`,
		`ALTER TABLE app.users ADD COLUMN IF NOT EXISTS x_user_id TEXT`,
		`ALTER TABLE app.users ADD COLUMN IF NOT EXISTS x_name TEXT`,
		`ALTER TABLE app.users ADD COLUMN IF NOT EXISTS x_username TEXT`,
		`ALTER TABLE app.users ADD COLUMN IF NOT EXISTS linkedin_access_token TEXT`,
		`ALTER TABLE app.users ADD COLUMN IF NOT EXISTS linkedin_refresh_token TEXT`,
		`ALTER TABLE app.users ADD COLUMN IF NOT EXISTS linkedin_profile_picture TEXT`,
		`ALTER TABLE app.users ADD COLUMN IF NOT EXISTS linkedin_user_id TEXT`,
		`ALTER TABLE app.users ADD COLUMN IF NOT EXISTS linkedin_name TEXT`,
	}

	for _, migration := range migrations {
		if _, err := DB.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %v", err)
		}
	}

	fmt.Println("Database migrations completed")
	return nil
}
