package sample

import (
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
)

func init() {
	gofakeit.Seed(time.Now().UnixNano())
}

func NewClientEntity(partnerID int, name string) *entities.Client {
	return NewClientEntityWithUser(partnerID, name)
}

func NewClientEntityWithUser(partnerID int, name string) *entities.Client {
	task := &entities.Client{
		Name:      name,
		PartnerID: &partnerID,
		CreatedAt: time.Now(),
	}
	return task
}
