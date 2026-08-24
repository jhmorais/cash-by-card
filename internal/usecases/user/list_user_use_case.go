package user

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
)

type listUserUseCase struct {
	userRepository repositories.UserRepository
}

func NewListUserUseCase(userRepository repositories.UserRepository) contracts.ListUserUseCase {
	return &listUserUseCase{userRepository: userRepository}
}

func (l *listUserUseCase) Execute(ctx context.Context, requesterEmail string) (*output.ListUser, error) {
	requester, err := l.userRepository.FindUserByEmail(ctx, requesterEmail)
	if err != nil {
		return nil, err
	}
	if err := canListUsers(requester); err != nil {
		return nil, err
	}

	users, err := l.userRepository.ListUser(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]output.UserItem, 0, len(users))
	for _, u := range users {
		items = append(items, output.UserItem{
			ID:                 u.ID,
			Name:               u.Name,
			Email:              u.Email,
			Role:               u.Role,
			PendingFirstAccess: u.Password == "",
		})
	}
	return &output.ListUser{Users: items}, nil
}
