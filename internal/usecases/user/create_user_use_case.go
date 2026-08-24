package user

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jhmorais/cash-by-card/config"
	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
	repoToken "github.com/jhmorais/cash-by-card/internal/repositories/token"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
	"github.com/jhmorais/cash-by-card/utils"
)

type createUserUseCase struct {
	userRepository  repositories.UserRepository
	tokenRepository repoToken.PasswordResetTokenRepository
	emailSender     contracts.EmailSender
}

func NewCreateUserUseCase(
	userRepository repositories.UserRepository,
	tokenRepository repoToken.PasswordResetTokenRepository,
	emailSender contracts.EmailSender,
) contracts.CreateUserUseCase {
	return &createUserUseCase{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		emailSender:     emailSender,
	}
}

func (c *createUserUseCase) Execute(ctx context.Context, requesterEmail string, createUser *input.CreateUser) (*output.CreateUser, error) {
	if len(createUser.Name) > 250 {
		createUser.Name = createUser.Name[:250]
	}
	if createUser.Email == "" {
		return nil, fmt.Errorf("email é obrigatório")
	}
	if createUser.Name == "" {
		return nil, fmt.Errorf("nome é obrigatório")
	}

	requester, err := c.userRepository.FindUserByEmail(ctx, requesterEmail)
	if err != nil {
		return nil, err
	}
	if err := canAssignRole(requester, createUser.Role); err != nil {
		return nil, err
	}

	existing, err := c.userRepository.FindUserByEmail(ctx, createUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	if existing != nil && existing.Email != "" {
		return nil, fmt.Errorf("já existe usuário com o mesmo email")
	}

	userEntity := &entities.User{
		Name:      createUser.Name,
		Email:     createUser.Email,
		Role:      createUser.Role,
		CreatedAt: time.Now(),
	} // Password fica vazio (NULL) => pendente de primeiro acesso

	if err := c.userRepository.CreateUser(ctx, userEntity); err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	if err := c.sendFirstAccessEmail(ctx, userEntity); err != nil {
		log.Printf("usuario criado mas falhou envio de primeiro acesso: %v", err)
	}

	return &output.CreateUser{UserID: userEntity.ID}, nil
}

// sendFirstAccessEmail gera o token e envia o link de definicao de senha.
func (c *createUserUseCase) sendFirstAccessEmail(ctx context.Context, user *entities.User) error {
	plain, hash, err := utils.GenerateResetToken()
	if err != nil {
		return err
	}
	token := &entities.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(30 * time.Minute),
		CreatedAt: time.Now(),
	}
	if err := c.tokenRepository.CreateToken(ctx, token); err != nil {
		return err
	}
	link := strings.TrimSuffix(config.GetFrontURL(), "/") + "/primeiro-acesso?token=" + plain
	return c.emailSender.SendPasswordResetEmail(ctx, user.Email, link)
}
