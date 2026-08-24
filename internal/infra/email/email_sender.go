package email

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/jhmorais/cash-by-card/internal/contracts"
)

type smtpSender struct {
	host, port, user, password, from string
}

type logSender struct{}

// NewSenderFromEnv devolve SmtpSender se SMTP_HOST estiver configurado;
// senão devolve o fallback que apenas loga o link (dev).
func NewSenderFromEnv() contracts.EmailSender {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		return &logSender{}
	}
	return &smtpSender{
		host:     host,
		port:     envOr("SMTP_PORT", "587"),
		user:     os.Getenv("SMTP_USER"),
		password: os.Getenv("SMTP_PASSWORD"),
		from:     envOr("SMTP_FROM", os.Getenv("SMTP_USER")),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func (s *smtpSender) SendPasswordResetEmail(ctx context.Context, to, resetLink string) error {
	subject := "Cash By Card - definicao de senha"
	body := fmt.Sprintf(
		"Ola,\n\nUse o link abaixo para definir sua senha. Ele expira em 30 minutos e pode ser usado uma unica vez:\n\n%s\n\nSe voce nao solicitou, ignore este email.\n",
		resetLink,
	)
	msg := strings.Join([]string{
		"From: " + s.from,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	addr := s.host + ":" + s.port
	auth := smtp.PlainAuth("", s.user, s.password, s.host)
	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}

func (s *logSender) SendPasswordResetEmail(ctx context.Context, to, resetLink string) error {
	log.Printf("[DEV-EMAIL] para=%s link=%s", to, resetLink)
	return nil
}
