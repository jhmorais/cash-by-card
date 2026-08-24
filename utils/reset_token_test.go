package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestGenerateResetToken(t *testing.T) {
	plain, hash, err := GenerateResetToken()
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if len(plain) != 64 {
		t.Fatalf("expected token with 64 hex chars, got %d", len(plain))
	}
	sum := sha256.Sum256([]byte(plain))
	if hash != hex.EncodeToString(sum[:]) {
		t.Fatal("expected hash to be the SHA-256 hex of the token")
	}
	plain2, _, _ := GenerateResetToken()
	if plain == plain2 {
		t.Fatal("expected two generated tokens to differ")
	}
}
