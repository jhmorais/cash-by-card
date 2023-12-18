package repositories

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
)

type ClientRepository interface {
	CreateClient(ctx context.Context, entity *entities.Client) error
	DeleteClient(ctx context.Context, entity *entities.Client) error
	UpdateClient(ctx context.Context, entity *entities.Client) error
	FindClientByID(ctx context.Context, id int) (*entities.Client, error)
	FindClientByName(ctx context.Context, name string) ([]*entities.Client, error)
	FindClientByCPF(ctx context.Context, cpf string) ([]*entities.Client, error)
	FindClientByPartnerID(ctx context.Context, partnerID int, name string) ([]*entities.Client, error)
	ListClient(ctx context.Context) ([]*entities.Client, error)
}
