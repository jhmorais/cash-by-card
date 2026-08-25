package user

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/jhmorais/cash-by-card/config"
	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	repoToken "github.com/jhmorais/cash-by-card/internal/repositories/token"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
	"github.com/jhmorais/cash-by-card/utils"
)

type forgotPasswordUseCase struct {
	userRepository  repositories.UserRepository
	tokenRepository repoToken.PasswordResetTokenRepository
	emailSender     contracts.EmailSender
}

func NewForgotPasswordUseCase(
	userRepository repositories.UserRepository,
	tokenRepository repoToken.PasswordResetTokenRepository,
	emailSender contracts.EmailSender,
) contracts.ForgotPasswordUseCase {
	return &forgotPasswordUseCase{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		emailSender:     emailSender,
	}
}

func (f *forgotPasswordUseCase) Execute(ctx context.Context, forgotPassword *input.ForgotPassword) error {
	if forgotPassword.Email == "" {
		return nil // resposta genérica também para input vazio
	}
	user, err := f.userRepository.FindUserByEmail(ctx, forgotPassword.Email)
	if err != nil || user == nil || user.Email == "" {
		return nil // nunca revelar existencia do email
	}

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
	if err := f.tokenRepository.CreateToken(ctx, token); err != nil {
		return err
	}

	// usuário sem senha => primeiro acesso; com senha => recuperação
	path := "/primeiro-acesso"
	if user.Password != "" {
		path = "/recuperar-senha"
	}
	link := strings.TrimSuffix(config.GetFrontURL(), "/") + path + "?token=" + plain
	if err := f.emailSender.SendPasswordResetEmail(ctx, user.Email, link); err != nil {
		log.Printf("falha ao enviar email de reset para %s: %v", user.Email, err)
	}
	return nil
}
