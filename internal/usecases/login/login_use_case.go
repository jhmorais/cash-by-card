package login

// Vai chamar utils para fazer o hash da senha
//		https://github.com/HunCoding/meu-primeiro-crud-go/blob/main/src/model/user_domain_password.go

// Vai chamar use case para consultar usuario pelo email e a senha encriptada
// Vai chamar use case para gerar o token jwt para responder ao usuario
//		https://github.com/HunCoding/meu-primeiro-crud-go/blob/main/src/model/user_token_domain.go#L18

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
	"github.com/jhmorais/cash-by-card/utils"
)

type loginUseCase struct {
	userRepository repositories.UserRepository
}

func NewLoginUseCase(userRepository repositories.UserRepository) contracts.LoginUseCase {

	return &loginUseCase{
		userRepository: userRepository,
	}
}

func (c *loginUseCase) Execute(ctx context.Context, loginUser *input.UserLogin) (string, error) {

	if loginUser.Email == "" {
		return "", fmt.Errorf("cannot login without email")
	}

	hashUser := utils.EncryptPassword(loginUser.Password)

	user, err := c.userRepository.FindUserByEmailandPassword(ctx, loginUser.Email, hashUser)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %v", err)
	}

	if user.Email == "" {
		return "", fmt.Errorf("failed, invalid inputs")
	}

	token, err := utils.GenerateToken(*loginUser)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}
