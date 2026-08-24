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
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
	repoToken "github.com/jhmorais/cash-by-card/internal/repositories/token"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
	"github.com/jhmorais/cash-by-card/utils"
)

type clearPasswordUseCase struct {
	userRepository  repositories.UserRepository
	tokenRepository repoToken.PasswordResetTokenRepository
	emailSender     contracts.EmailSender
}

func NewClearPasswordUseCase(
	userRepository repositories.UserRepository,
	tokenRepository repoToken.PasswordResetTokenRepository,
	emailSender contracts.EmailSender,
) contracts.ClearPasswordUseCase {
	return &clearPasswordUseCase{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		emailSender:     emailSender,
	}
}

func (c *clearPasswordUseCase) Execute(ctx context.Context, requesterEmail string, id int) (*output.ClearPassword, error) {
	requester, err := c.userRepository.FindUserByEmail(ctx, requesterEmail)
	if err != nil {
		return nil, err
	}
	target, err := c.userRepository.FindUserByID(ctx, id)
	if err != nil || target == nil || target.ID == 0 {
		return nil, fmt.Errorf("usuário não encontrado")
	}
	if err := canManageTarget(requester, target); err != nil {
		return nil, err
	}

	// token ANTES de limpar a senha: se a criação do token falhar, a senha fica
	// intacta e o usuário não fica trancado fora da conta.
	plain, hash, err := utils.GenerateResetToken()
	if err != nil {
		return nil, err
	}
	token := &entities.PasswordResetToken{
		UserID:    target.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(30 * time.Minute),
		CreatedAt: time.Now(),
	}
	if err := c.tokenRepository.CreateToken(ctx, token); err != nil {
		return nil, err
	}

	if err := c.userRepository.ClearUserPassword(ctx, target.ID); err != nil {
		return nil, err
	}

	link := strings.TrimSuffix(config.GetFrontURL(), "/") + "/primeiro-acesso?token=" + plain
	if err := c.emailSender.SendPasswordResetEmail(ctx, target.Email, link); err != nil {
		log.Printf("senha limpa mas falhou envio do email: %v", err)
	}

	return &output.ClearPassword{UserID: target.ID}, nil
}
