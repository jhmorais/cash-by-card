package user

import (
	"context"
	"testing"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
)

func TestGetUser_PorEmail(t *testing.T) {
	uc := NewGetUserUseCase(requesterRepo(RolePartner, 4))
	out, err := uc.Execute(context.Background(), "p@x.com")
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if out.Email != "p@x.com" || out.Role != "partner" {
		t.Fatalf("esperado p@x.com/partner, got %s/%s", out.Email, out.Role)
	}
}

func TestGetUser_NaoEncontrado(t *testing.T) {
	uc := NewGetUserUseCase(&mockUserRepo{findByEmail: func(ctx context.Context, email string) (*entities.User, error) {
		return &entities.User{}, nil
	}})
	_, err := uc.Execute(context.Background(), "ninguem@x.com")
	if err == nil {
		t.Fatal("usuário inexistente deve retornar erro")
	}
}
