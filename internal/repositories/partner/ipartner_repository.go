package repositories

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
)

type PartnerRepository interface {
	CreatePartner(ctx context.Context, entity *entities.Partner) error
	DeletePartner(ctx context.Context, entity *entities.Partner) error
	UpdatePartner(ctx context.Context, entity *entities.Partner) error
	FindPartnerByID(ctx context.Context, id int) (*entities.Partner, error)
	FindPartnerByName(ctx context.Context, name string) ([]*entities.Partner, error)
	FindPartnerByEmail(ctx context.Context, email string) ([]*entities.Partner, error)
	FindPartnerByCPF(ctx context.Context, cpf string) ([]*entities.Partner, error)
	FindPartnerByPartnerID(ctx context.Context, partnerID int, name string) ([]*entities.Partner, error)
	ListPartner(ctx context.Context) ([]*entities.Partner, error)
}
