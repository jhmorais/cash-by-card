package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
)

var (
	JWT_SECRET_KEY = "JWT_SECRET_KEY"
)

func GenerateToken(user input.UserLogin) (string, error) {
	secret := os.Getenv(JWT_SECRET_KEY)

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
