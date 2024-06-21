package utils

import "strings"

func RemoveBearerPrefix(token string) string {
	return strings.TrimPrefix(token, "Bearer ")
}
