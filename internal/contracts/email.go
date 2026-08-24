package contracts

import "context"

// EmailSender envia emails transacionais. Implementação SMTP em internal/infra/email.
type EmailSender interface {
	SendPasswordResetEmail(ctx context.Context, to, resetLink string) error
}
