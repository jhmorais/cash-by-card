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

func TestLogin_FluxosExistentes(t *testing.T) {
	t.Run("senha errada retorna erro generico", func(t *testing.T) {
		uc := NewLoginUseCase(&mockUserRepoLogin{
			byEmail: func(ctx context.Context, email string) (*entities.User, error) {
				return &entities.User{ID: 1, Email: email, Password: "hash-correto"}, nil
			},
			byEmailP: func(ctx context.Context, email, password string) (*entities.User, error) {
				return &entities.User{}, nil
			},
		})
		_, err := uc.Execute(context.Background(), &input.UserLogin{Email: "a@b.com", Password: "errada"})
		if err == nil || strings.Contains(err.Error(), "primeiro acesso") {
			t.Fatalf("esperado erro generico de login, got '%v'", err)
		}
	})

	t.Run("login valido gera token", func(t *testing.T) {
		t.Setenv("JWT_SECRET_KEY", "test-secret")
		uc := NewLoginUseCase(&mockUserRepoLogin{
			byEmail: func(ctx context.Context, email string) (*entities.User, error) {
				return &entities.User{ID: 1, Email: email, Password: "hash-correto"}, nil
			},
			byEmailP: func(ctx context.Context, email, password string) (*entities.User, error) {
				return &entities.User{ID: 1, Email: email, Password: "hash-correto"}, nil
			},
		})
		token, err := uc.Execute(context.Background(), &input.UserLogin{Email: "a@b.com", Password: "correta"})
		if err != nil {
			t.Fatalf("expected no error, got '%v'", err)
		}
		if token == "" {
			t.Fatal("esperado token nao vazio")
		}
	})

	t.Run("email inexistente retorna erro generico", func(t *testing.T) {
		uc := NewLoginUseCase(&mockUserRepoLogin{
			byEmail: func(ctx context.Context, email string) (*entities.User, error) {
				return &entities.User{}, nil // gorm Find: entidade vazia, sem erro
			},
			byEmailP: func(ctx context.Context, email, password string) (*entities.User, error) {
				return &entities.User{}, nil
			},
		})
		_, err := uc.Execute(context.Background(), &input.UserLogin{Email: "ninguem@x.com", Password: "qualquer"})
		if err == nil || strings.Contains(err.Error(), "primeiro acesso") {
			t.Fatalf("esperado erro generico, got '%v'", err)
		}
	})
}
