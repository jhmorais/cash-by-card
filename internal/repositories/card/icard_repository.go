package repositories

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
)

type CardRepository interface {
	CreateCard(ctx context.Context, entity *entities.Card) error
	DeleteCard(ctx context.Context, entity *entities.Card) error
	UpdateCard(ctx context.Context, entity []*entities.Card) error
	FindCardByID(ctx context.Context, id int) (*entities.Card, error)
	FindCardByLoanID(ctx context.Context, LoanID int) ([]*entities.Card, error)
	ListCard(ctx context.Context) ([]*entities.Card, error)
}
