package user

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
	"github.com/jhmorais/cash-by-card/utils"
)

type changePasswordUseCase struct {
	userRepository repositories.UserRepository
}

func NewChangePasswordUseCase(userRepository repositories.UserRepository) contracts.ChangePasswordUseCase {
	return &changePasswordUseCase{userRepository: userRepository}
}

func (c *changePasswordUseCase) Execute(ctx context.Context, email string, changePassword *input.ChangePassword) error {
	if err := utils.ValidatePasswordPolicy(changePassword.NewPassword); err != nil {
		return err
	}
	user, err := c.userRepository.FindUserByEmail(ctx, email)
	if err != nil || user == nil || user.Email == "" {
		return fmt.Errorf("usuário não encontrado")
	}
	if user.Password == "" || user.Password != utils.EncryptPassword(changePassword.CurrentPassword) {
		return fmt.Errorf("senha atual incorreta")
	}
	return c.userRepository.SetUserPassword(ctx, user.ID, utils.EncryptPassword(changePassword.NewPassword))
}
