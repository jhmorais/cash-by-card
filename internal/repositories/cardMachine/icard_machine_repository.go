package repositories

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
)

type CardMachineRepository interface {
	CreateCardMachine(ctx context.Context, entity *entities.CardMachine) error
	DeleteCardMachine(ctx context.Context, entity *entities.CardMachine) error
	UpdateCardMachine(ctx context.Context, entity *entities.CardMachine) error
	FindCardMachineByID(ctx context.Context, id int) (*entities.CardMachine, error)
	ListCardMachine(ctx context.Context) ([]*entities.CardMachine, error)
}
