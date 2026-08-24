package user

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
)

type getUserUseCase struct {
	userRepository repositories.UserRepository
}

func NewGetUserUseCase(userRepository repositories.UserRepository) contracts.GetUserUseCase {
	return &getUserUseCase{userRepository: userRepository}
}

func (g *getUserUseCase) Execute(ctx context.Context, email string) (*output.GetUser, error) {
	user, err := g.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Email == "" {
		return nil, fmt.Errorf("usuário não encontrado")
	}
	return &output.GetUser{ID: user.ID, Email: user.Email, Name: user.Name, Role: user.Role}, nil
}
