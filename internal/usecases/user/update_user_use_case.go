package user

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
)

type updateUserUseCase struct {
	userRepository repositories.UserRepository
}

func NewUpdateUserUseCase(userRepository repositories.UserRepository) contracts.UpdateUserUseCase {
	return &updateUserUseCase{userRepository: userRepository}
}

func (u *updateUserUseCase) Execute(ctx context.Context, requesterEmail string, updateUser *input.UpdateUser) (*output.UpdateUser, error) {
	if updateUser.Name == "" {
		return nil, fmt.Errorf("nome é obrigatório")
	}
	requester, err := u.userRepository.FindUserByEmail(ctx, requesterEmail)
	if err != nil {
		return nil, err
	}
	target, err := u.userRepository.FindUserByID(ctx, updateUser.ID)
	if err != nil || target == nil || target.ID == 0 {
		return nil, fmt.Errorf("usuário não encontrado")
	}
	if err := canManageTarget(requester, target); err != nil {
		return nil, err
	}
	if updateUser.Role != "" && updateUser.Role != target.Role {
		if err := canAssignRole(requester, updateUser.Role); err != nil {
			return nil, err
		}
		target.Role = updateUser.Role
	}
	target.Name = updateUser.Name

	if err := u.userRepository.UpdateUser(ctx, target); err != nil {
		return nil, err
	}
	return &output.UpdateUser{UserID: target.ID}, nil
}
