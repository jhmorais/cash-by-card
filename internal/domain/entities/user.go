package entities

import "time"

type User struct {
	ID        int    `gorm:"id" json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Password  string `json:"-"`
	Role      string `json:"role"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
