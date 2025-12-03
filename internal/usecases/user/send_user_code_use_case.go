package user

import (
	"context"
	"fmt"
	"time"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	"github.com/jhmorais/cash-by-card/internal/infra/mail"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
	"github.com/jhmorais/cash-by-card/utils"
)

type sendUserCodeUseCase struct {
	userRepository repositories.UserRepository
	mailer         *mail.SMTPMailer
}

func NewSendUserCodeUseCase(userRepository repositories.UserRepository, mailer *mail.SMTPMailer) contracts.SendUserCodeUseCase {

	return &sendUserCodeUseCase{
		userRepository: userRepository,
		mailer:         mailer,
	}
}

func (c *sendUserCodeUseCase) Execute(ctx context.Context, email string) error {

	user, err := c.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	code := utils.GenerateUUID()

	userEntity := &entities.User{
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Password:  user.Password,
		Code:      code,
		CreatedAt: user.CreatedAt,
		UpdatedAt: time.Now(),
	}

	err = c.userRepository.UpdateUser(ctx, userEntity)
	if err != nil {
		return fmt.Errorf("cannot save user at database: %v", err)
	}

	subject := "Seu código de verificação"
	body := fmt.Sprintf("Olá %s,\n\nSeu código é: %s\n\nEle expira em 10 minutos.\n", user.Name, code)

	if c.mailer != nil {
		if err := c.mailer.Send(ctx, email, subject, body); err != nil {
			return fmt.Errorf("failed to send email: %v", err)
		}
	}

	return nil
}
