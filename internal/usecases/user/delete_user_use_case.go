package user

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
	repoToken "github.com/jhmorais/cash-by-card/internal/repositories/token"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
)

type deleteUserUseCase struct {
	userRepository  repositories.UserRepository
	tokenRepository repoToken.PasswordResetTokenRepository
}

func NewDeleteUserUseCase(
	userRepository repositories.UserRepository,
	tokenRepository repoToken.PasswordResetTokenRepository,
) contracts.DeleteUserUseCase {
	return &deleteUserUseCase{userRepository: userRepository, tokenRepository: tokenRepository}
}

func (d *deleteUserUseCase) Execute(ctx context.Context, requesterEmail string, id int) (*output.DeleteUser, error) {
	requester, err := d.userRepository.FindUserByEmail(ctx, requesterEmail)
	if err != nil {
		return nil, err
	}
	target, err := d.userRepository.FindUserByID(ctx, id)
	if err != nil || target == nil || target.ID == 0 {
		return nil, fmt.Errorf("usuário não encontrado")
	}
	if err := canManageTarget(requester, target); err != nil {
		return nil, err
	}
	// tokens ANTES do user: a FK password_reset_token.user_id impediria o DELETE (erro 1451)
	if d.tokenRepository != nil {
		if err := d.tokenRepository.DeleteTokensByUser(ctx, target.ID); err != nil {
			return nil, err
		}
	}
	if err := d.userRepository.DeleteUser(ctx, target); err != nil {
		return nil, err
	}
	return &output.DeleteUser{Success: true}, nil
}
