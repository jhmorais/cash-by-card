package user

import (
	"context"
	"testing"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	"github.com/jhmorais/cash-by-card/utils"
)

func TestChangePassword_SenhaAtualErrada(t *testing.T) {
	uc := NewChangePasswordUseCase(&mockUserRepo{findByEmail: func(ctx context.Context, email string) (*entities.User, error) {
		return &entities.User{ID: 1, Email: email, Password: utils.EncryptPassword("correta")}, nil
	}})
	err := uc.Execute(context.Background(), "a@b.com", &input.ChangePassword{CurrentPassword: "errada", NewPassword: "NovaSenha@1"})
	if err == nil {
		t.Fatal("senha atual errada deve retornar erro")
	}
}

func TestChangePassword_SenhaFraca(t *testing.T) {
	uc := NewChangePasswordUseCase(&mockUserRepo{findByEmail: func(ctx context.Context, email string) (*entities.User, error) {
		return &entities.User{ID: 1, Email: email, Password: utils.EncryptPassword("correta")}, nil
	}})
	err := uc.Execute(context.Background(), "a@b.com", &input.ChangePassword{CurrentPassword: "correta", NewPassword: "novasenha"})
	if err == nil {
		t.Fatal("senha sem maiuscula, numero e caractere especial deve retornar erro")
	}
}

func TestChangePassword_Sucesso(t *testing.T) {
	var savedID int
	var savedHash string
	uc := NewChangePasswordUseCase(&mockUserRepoSetByEmail{
		mockUserRepo: mockUserRepo{
			findByEmail: func(ctx context.Context, email string) (*entities.User, error) {
				return &entities.User{ID: 3, Email: email, Password: utils.EncryptPassword("correta")}, nil
			},
		},
		setFunc: func(ctx context.Context, id int, hash string) error {
			savedID = id
			savedHash = hash
			return nil
		},
	})
	err := uc.Execute(context.Background(), "a@b.com", &input.ChangePassword{CurrentPassword: "correta", NewPassword: "NovaSenha@1"})
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if savedID != 3 || savedHash != utils.EncryptPassword("NovaSenha@1") {
		t.Fatalf("senha não salva: id=%d hash=%s", savedID, savedHash)
	}
}

type mockUserRepoSetByEmail struct {
	mockUserRepo
	setFunc func(ctx context.Context, id int, hash string) error
}

func (m *mockUserRepoSetByEmail) SetUserPassword(ctx context.Context, id int, hash string) error {
	return m.setFunc(ctx, id, hash)
}
