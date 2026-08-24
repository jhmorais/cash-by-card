package repositories

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
)

type PasswordResetTokenRepository interface {
	CreateToken(ctx context.Context, entity *entities.PasswordResetToken) error
	// FindValidTokenByHash devolve nil, nil quando não há token válido (hash inexistente, expirado ou já usado).
	FindValidTokenByHash(ctx context.Context, tokenHash string) (*entities.PasswordResetToken, error)
	MarkTokenUsed(ctx context.Context, id int64) error
}
