package user

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	repoToken "github.com/jhmorais/cash-by-card/internal/repositories/token"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
	"github.com/jhmorais/cash-by-card/utils"
)

type resetPasswordUseCase struct {
	tokenRepository repoToken.PasswordResetTokenRepository
	userRepository  repositories.UserRepository
}

func NewResetPasswordUseCase(
	tokenRepository repoToken.PasswordResetTokenRepository,
	userRepository repositories.UserRepository,
) contracts.ResetPasswordUseCase {
	return &resetPasswordUseCase{tokenRepository: tokenRepository, userRepository: userRepository}
}

func (r *resetPasswordUseCase) Execute(ctx context.Context, resetPassword *input.ResetPassword) error {
	if resetPassword.Token == "" {
		return fmt.Errorf("token é obrigatório")
	}
	if err := utils.ValidatePasswordPolicy(resetPassword.NewPassword); err != nil {
		return err
	}

	hash := utils.HashResetToken(resetPassword.Token)
	token, err := r.tokenRepository.FindValidTokenByHash(ctx, hash)
	if err != nil {
		return fmt.Errorf("failed to validate token: %v", err)
	}
	if token == nil {
		return fmt.Errorf("token inválido ou expirado")
	}

	if err := r.userRepository.SetUserPassword(ctx, token.UserID, utils.EncryptPassword(resetPassword.NewPassword)); err != nil {
		return fmt.Errorf("failed to save password: %v", err)
	}

	if err := r.tokenRepository.MarkTokenUsed(ctx, token.ID); err != nil {
		return fmt.Errorf("failed to invalidate token: %v", err)
	}
	return nil
}
