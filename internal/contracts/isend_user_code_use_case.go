package contracts

import (
	"context"
)

type SendUserCodeUseCase interface {
	Execute(ctx context.Context, email string) error
}
