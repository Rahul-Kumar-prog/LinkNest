package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwtToken(UserID, UserEmail string) (string, error) {
	secret := "Rahul@1191"
	log.Println(secret)
	if secret == "" {
		return "", fmt.Errorf("JWT_KEY not set")
	}

	key := []byte(secret)

	claims := jwt.MapClaims{
		"user_id": UserID,
		"email":   UserEmail,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	s, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return s, nil
}
