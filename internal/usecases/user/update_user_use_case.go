package user

// Vai receber um input com os dados do usuario: nome, senha, email, role

// Vai chamar o use case para pesquisar se ja existe um usuario com o mesmo email
// Vai chamar o util para criar um hash da senha do input
//		https://github.com/HunCoding/meu-primeiro-crud-go/blob/main/src/model/user_domain_password.go

// Vai chamar um repository para salvar no banco o usuario

import (
	"context"
	"fmt"
	"time"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
	"github.com/jhmorais/cash-by-card/utils"
)

type updateUserUseCase struct {
	userRepository repositories.UserRepository
}

func NewUpdateUserUseCase(userRepository repositories.UserRepository) contracts.UpdateUserUseCase {

	return &updateUserUseCase{
		userRepository: userRepository,
	}
}

func (c *updateUserUseCase) Execute(ctx context.Context, updateUser *input.UpdateUser) error {

	if len(updateUser.Password) < 6 {
		return fmt.Errorf("cannot update a user without valid password")
	}

	user, err := c.userRepository.FindUserByEmail(ctx, updateUser.Email)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	minutesToAdd := 10
	if user.Code != updateUser.Code || time.Now().After(user.UpdatedAt.Add(time.Duration(minutesToAdd)*time.Minute)) {
		return fmt.Errorf("invalid or expired code")
	}

	hashUser := utils.EncryptPassword(updateUser.Password)

	userEntity := &entities.User{
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Password:  hashUser,
		Code:      "",
		CreatedAt: user.CreatedAt,
		UpdatedAt: time.Now(),
	}

	err = c.userRepository.UpdateUser(ctx, userEntity)
	if err != nil {
		return fmt.Errorf("cannot save user at database: %v", err)
	}

	return nil
}
