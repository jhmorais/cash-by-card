package repositories

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
)

type UserRepository interface {
	CreateUser(ctx context.Context, entity *entities.User) error
	DeleteUser(ctx context.Context, entity *entities.User) error
	UpdateUser(ctx context.Context, entity *entities.User) error
	FindUserByID(ctx context.Context, id int) (*entities.User, error)
	FindUserByEmail(ctx context.Context, email string) ([]*entities.User, error)
	FindUserByEmailandPassword(ctx context.Context, email string, password string) (*entities.User, error)
	ListUser(ctx context.Context) ([]*entities.User, error)
}
