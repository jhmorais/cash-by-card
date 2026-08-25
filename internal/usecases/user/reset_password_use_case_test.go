package user

import (
	"context"
	"testing"
	"time"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	repoToken "github.com/jhmorais/cash-by-card/internal/repositories/token"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
	"github.com/jhmorais/cash-by-card/utils"
)

type mockTokenRepoFull struct {
	repoToken.PasswordResetTokenRepository
	findFunc func(ctx context.Context, hash string) (*entities.PasswordResetToken, error)
	markFunc func(ctx context.Context, id int64) error
}

func (m *mockTokenRepoFull) FindValidTokenByHash(ctx context.Context, hash string) (*entities.PasswordResetToken, error) {
	return m.findFunc(ctx, hash)
}
func (m *mockTokenRepoFull) MarkTokenUsed(ctx context.Context, id int64) error { return m.markFunc(ctx, id) }

type mockUserRepoWithSet struct {
	repositories.UserRepository
	setFunc func(ctx context.Context, id int, hash string) error
}

func (m *mockUserRepoWithSet) SetUserPassword(ctx context.Context, id int, hash string) error {
	return m.setFunc(ctx, id, hash)
}

func TestResetPassword_TokenInvalido(t *testing.T) {
	uc := NewResetPasswordUseCase(
		&mockTokenRepoFull{findFunc: func(ctx context.Context, hash string) (*entities.PasswordResetToken, error) {
			return nil, nil
		}},
		&mockUserRepoWithSet{},
	)
	err := uc.Execute(context.Background(), &input.ResetPassword{Token: "x", NewPassword: "NovaSenha@1"})
	if err == nil {
		t.Fatal("token invalido deve retornar erro")
	}
}

func TestResetPassword_SenhaCurta(t *testing.T) {
	uc := NewResetPasswordUseCase(&mockTokenRepoFull{}, &mockUserRepoWithSet{})
	err := uc.Execute(context.Background(), &input.ResetPassword{Token: "x", NewPassword: "123"})
	if err == nil {
		t.Fatal("senha curta deve retornar erro")
	}
}

func TestResetPassword_SenhaFraca(t *testing.T) {
	uc := NewResetPasswordUseCase(&mockTokenRepoFull{}, &mockUserRepoWithSet{})
	err := uc.Execute(context.Background(), &input.ResetPassword{Token: "x", NewPassword: "novasenha"})
	if err == nil {
		t.Fatal("senha sem maiuscula, numero e caractere especial deve retornar erro")
	}
}

func TestResetPassword_Sucesso(t *testing.T) {
	token := &entities.PasswordResetToken{ID: 55, UserID: 9, ExpiresAt: time.Now().Add(time.Minute)}
	var markedID int64
	var savedHash string
	var savedID int
	uc := NewResetPasswordUseCase(
		&mockTokenRepoFull{
			findFunc: func(ctx context.Context, hash string) (*entities.PasswordResetToken, error) {
				return token, nil
			},
			markFunc: func(ctx context.Context, id int64) error { markedID = id; return nil },
		},
		&mockUserRepoWithSet{setFunc: func(ctx context.Context, id int, hash string) error {
			savedID = id
			savedHash = hash
			return nil
		}},
	)
	if err := uc.Execute(context.Background(), &input.ResetPassword{Token: "plain-token", NewPassword: "NovaSenha@1"}); err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if markedID != 55 {
		t.Fatalf("esperado marcar token 55 como usado, got %d", markedID)
	}
	if savedID != 9 || savedHash != utils.EncryptPassword("NovaSenha@1") {
		t.Fatalf("senha não salva corretamente: id=%d hash=%s", savedID, savedHash)
	}
}
