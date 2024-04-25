package login

// Vai chamar utils para fazer o hash da senha
//		https://github.com/HunCoding/meu-primeiro-crud-go/blob/main/src/model/user_domain_password.go

// Vai chamar use case para consultar usuario pelo email e a senha encriptada
// Vai chamar use case para gerar o token jwt para responder ao usuario
//		https://github.com/HunCoding/meu-primeiro-crud-go/blob/main/src/model/user_token_domain.go#L18

import (
	"context"
	"fmt"
	"time"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
	"github.com/jhmorais/cash-by-card/utils"
)

type loginUseCase struct {
	userRepository repositories.UserRepository
}

func NewLoginUseCase(userRepository repositories.UserRepository) contracts.CreateUserUseCase {

	return &loginUseCase{
		userRepository: userRepository,
	}
}

func (c *loginUseCase) Execute(ctx context.Context, loginUser *input.CreateUser) (*output.CreateUser, error) {

	if loginUser.Email == "" {
		return nil, fmt.Errorf("cannot login without email")
	}

	if loginUser.Role != "admin" && loginUser.Role != "regular" {
		return nil, fmt.Errorf("cannot login without valid role")
	}

	if len(loginUser.Password) < 6 {
		return nil, fmt.Errorf("cannot login without valid password")
	}

	hashUser := utils.EncryptPassword(loginUser.Password)

	user, err := c.userRepository.FindUserByEmailandPassword(ctx, loginUser.Email, hashUser)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	if len(user) < 1 {
		return nil, fmt.Errorf("failed, invalid inputs")
	}

	userEntity := &entities.User{
		Name:      createUser.Name,
		Email:     createUser.Email,
		Role:      createUser.Role,
		Password:  hashUser,
		CreatedAt: time.Now(),
	}

	err = c.userRepository.CreateUser(ctx, userEntity)
	if err != nil {
		return nil, fmt.Errorf("cannot save user at database: %v", err)
	}

	createUserOutput := &output.CreateUser{
		UserID: userEntity.ID,
	}

	return createUserOutput, nil
}
