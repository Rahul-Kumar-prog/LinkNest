package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "12345"
	dbname   = "linknest"
)

var DB *sql.DB

func Connect() error {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("DB open Error:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("DB ping Error:", err)
	}
	fmt.Println("✅ Database connected successfully")
	return nil
}
