package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
)

const (
	ErrTypeError = "ERROR"
)

type ErrorModel struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func ValidateJwtTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := os.Getenv(JWT_SECRET_KEY)
		tokenValue := RemoveBearerPrefix(r.Header.Get("Authorization"))

		token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
				return []byte(secret), nil
			}

			return nil, errors.New("invalid token")
		})
		if err != nil {
			WriteErrModel(w, http.StatusNotFound,
				NewErrorResponse(fmt.Sprintf("failed to parse token, error: '%s'", err)))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			WriteErrModel(w, http.StatusNotFound,
				NewErrorResponse(fmt.Sprintf("failed to claims token, error: '%s'", errors.New("invalid token"))))
			return
		}

		userDomain := input.UserLogin{
			ID:    claims["id"].(int),
			Email: claims["email"].(string),
			Role:  claims["role"].(string),
		}
		fmt.Printf("User authenticated: %#v", userDomain)
		// w.Header().Get("jwt") // nome do campo do token, ULTIMO passo
		next.ServeHTTP(w, r)
	})
}

func WriteErrModel(w http.ResponseWriter, statusCode int, errModel *ErrorModel) {
	jsonStr, err := json.Marshal(errModel)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	fmt.Fprint(w, string(jsonStr))
}

func NewErrorResponse(msg string) *ErrorModel {
	return &ErrorModel{
		Message: msg,
		Type:    ErrTypeError,
	}
}

func RetrieveParam(r *http.Request, idParam string) (string, error) {
	encodedID := mux.Vars(r)[idParam]
	decodedID, err := url.QueryUnescape(encodedID)
	if err != nil {
		return "", err
	}
	return decodedID, nil
}

func ValidJSON(p interface{}) io.Reader {
	data, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return bytes.NewReader(data)
}
