package user

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
)

type listUserUseCase struct {
	userRepository repositories.UserRepository
}

func NewListUsertUseCase(userRepository repositories.UserRepository) contracts.ListUserUseCase {

	return &listUserUseCase{
		userRepository: userRepository,
	}
}

func (l *listUserUseCase) Execute(ctx context.Context) ([]*output.FindUser, error) {
	var err error

	result, err := l.userRepository.ListUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("error when list users on database: %v", err)
	}

	var listUser []*output.FindUser
	for _, entity := range result {
		user := &output.FindUser{
			ID:    entity.ID,
			Name:  entity.Name,
			Email: entity.Email,
			Role:  entity.Role,
		}
		listUser = append(listUser, user)
	}

	return listUser, nil
}
