package db

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string
	Email    string
	Password string
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	err := DB.QueryRow(
		`SELECT id, email, password FROM app.users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func CreateUser(email, password string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	result, err := DB.Exec(
		`INSERT INTO app.users (email, password)
		 VALUES ($1, $2)`,
		email,
		string(hashedPassword),
	)

	// 🔴 THIS CHECK WAS MISSING
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 1 {
		fmt.Println("✅ User created successfully:", email)
	}

	return nil
}
