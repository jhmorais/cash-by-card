package user

import (
	"context"
	"testing"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
)

type mockUserRepoList struct {
	mockUserRepo
	listFunc func(ctx context.Context) ([]*entities.User, error)
}

func (m *mockUserRepoList) ListUser(ctx context.Context) ([]*entities.User, error) {
	return m.listFunc(ctx)
}

func requesterRepo(role string, id int) *mockUserRepo {
	return &mockUserRepo{findByEmail: func(ctx context.Context, email string) (*entities.User, error) {
		return &entities.User{ID: id, Email: email, Role: role}, nil
	}}
}

func TestListUser_PartnerSemPermissao(t *testing.T) {
	uc := NewListUserUseCase(requesterRepo(RolePartner, 1))
	_, err := uc.Execute(context.Background(), "p@x.com")
	if err == nil {
		t.Fatal("partner não pode listar usuários")
	}
}

func TestListUser_AdminOk(t *testing.T) {
	uc := NewListUserUseCase(&mockUserRepoList{
		mockUserRepo: *requesterRepo(RoleAdmin, 1),
		listFunc: func(ctx context.Context) ([]*entities.User, error) {
			return []*entities.User{
				{ID: 2, Email: "x@x.com", Role: "partner"},              // password vazio → pendente
				{ID: 3, Email: "y@x.com", Role: "admin", Password: "h"}, // ok
			}, nil
		},
	})
	out, err := uc.Execute(context.Background(), "a@x.com")
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if len(out.Users) != 2 {
		t.Fatalf("esperado 2 usuários, got %d", len(out.Users))
	}
	if !out.Users[0].PendingFirstAccess {
		t.Fatal("usuário com password vazio deve vir PendingFirstAccess=true")
	}
	if out.Users[1].PendingFirstAccess {
		t.Fatal("usuário com password definido deve vir PendingFirstAccess=false")
	}
}
