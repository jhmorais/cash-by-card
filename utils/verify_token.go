package utils

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
)

func VerifyToken(tokenValue string) (*output.FindUser, error) {
	secret := os.Getenv(JWT_SECRET_KEY)

	token, err := jwt.Parse(RemoveBearerPrefix(tokenValue), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(secret), nil
		}

		return nil, errors.New("bad request - invalid token")
	})
	if err != nil {
		return nil, errors.New("unauthorized - invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("unauthorized - invalid token")
	}

	return &output.FindUser{
		ID:    claims["id"].(int),
		Email: claims["email"].(string),
		Role:  claims["role"].(string),
	}, errors.New("invalid token")
}
