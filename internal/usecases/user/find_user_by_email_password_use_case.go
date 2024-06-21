package user

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
)

type findUserByEmailAndPasswordUseCase struct {
	userRepository repositories.UserRepository
}

func NewFindUserByEmailAndPasswordUseCase(userRepository repositories.UserRepository) contracts.FindUserByEmailAndPasswordUseCase {

	return &findUserByEmailAndPasswordUseCase{
		userRepository: userRepository,
	}
}

func (c *findUserByEmailAndPasswordUseCase) Execute(ctx context.Context, email string, password string) (*output.FindUser, error) {

	userEntity, err := c.userRepository.FindUserByEmailandPassword(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("erro to find user with email: '%s' at database: '%v'", email, err)
	}

	if userEntity == nil {
		return nil, fmt.Errorf("user not found")
	}

	output := &output.FindUser{
		ID:    userEntity.ID,
		Email: userEntity.Email,
		Name:  userEntity.Name,
		Role:  userEntity.Role,
	}

	return output, nil
}
