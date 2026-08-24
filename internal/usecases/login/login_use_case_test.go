package login

import (
	"context"
	"strings"
	"testing"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
)

type mockUserRepoLogin struct {
	repositories.UserRepository
	byEmail  func(ctx context.Context, email string) (*entities.User, error)
	byEmailP func(ctx context.Context, email, password string) (*entities.User, error)
}

func (m *mockUserRepoLogin) FindUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	return m.byEmail(ctx, email)
}
func (m *mockUserRepoLogin) FindUserByEmailandPassword(ctx context.Context, email, password string) (*entities.User, error) {
	return m.byEmailP(ctx, email, password)
}

func TestLogin_UsuarioPendentePrimeiroAcesso(t *testing.T) {
	uc := NewLoginUseCase(&mockUserRepoLogin{
		byEmail: func(ctx context.Context, email string) (*entities.User, error) {
			return &entities.User{ID: 1, Email: email, Password: ""}, nil // password NULL no banco
		},
		byEmailP: func(ctx context.Context, email, password string) (*entities.User, error) {
			return &entities.User{}, nil
		},
	})
	_, err := uc.Execute(context.Background(), &input.UserLogin{Email: "a@b.com", Password: "qualquer"})
	if err == nil || !strings.Contains(err.Error(), "primeiro acesso") {
		t.Fatalf("esperado erro de primeiro acesso, got '%v'", err)
	}
}
