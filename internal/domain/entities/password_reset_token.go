package entities

import "time"

// PasswordResetToken guarda o hash SHA-256 (hex) do token enviado por email.
// O token em texto puro nunca é persistido.
type PasswordResetToken struct {
	ID        int64      `gorm:"id" json:"id"`
	UserID    int        `json:"user_id"`
	TokenHash string     `json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at"`
	CreatedAt time.Time  `json:"created_at"`
}
