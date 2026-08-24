package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
)

// ForgotPasswordUseCase envia email de definicao de senha se o email existir.
// Nunca retorna erro por email inexistente (anti-enumeração).
type ForgotPasswordUseCase interface {
	Execute(ctx context.Context, forgotPassword *input.ForgotPassword) error
}
