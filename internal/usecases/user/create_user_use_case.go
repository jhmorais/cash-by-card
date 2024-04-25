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
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
	"github.com/jhmorais/cash-by-card/utils"
)

type createUserUseCase struct {
	userRepository repositories.UserRepository
}

func NewCreateUserUseCase(userRepository repositories.UserRepository) contracts.CreateUserUseCase {

	return &createUserUseCase{
		userRepository: userRepository,
	}
}

func (c *createUserUseCase) Execute(ctx context.Context, createUser *input.CreateUser) (*output.CreateUser, error) {

	if len(createUser.Name) > 250 {
		createUser.Name = createUser.Name[:250]
	}

	if createUser.Email == "" {
		return nil, fmt.Errorf("cannot create a client without email")
	}

	if createUser.Role != "admin" && createUser.Role != "regular" {
		return nil, fmt.Errorf("cannot create a client without valid role")
	}

	if len(createUser.Password) < 6 {
		return nil, fmt.Errorf("cannot create a client without valid password")
	}

	user, err := c.userRepository.FindUserByEmail(ctx, createUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %v", err)
	}

	if len(user) > 0 {
		return nil, fmt.Errorf("failed, already exists user with the same email")
	}

	hashUser := utils.EncryptPassword(createUser.Password)

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
