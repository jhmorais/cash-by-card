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

func TestListUser_AdminVeApenasParceiros(t *testing.T) {
	uc := NewListUserUseCase(&mockUserRepoList{
		mockUserRepo: *requesterRepo(RoleAdmin, 1),
		listFunc: func(ctx context.Context) ([]*entities.User, error) {
			return []*entities.User{
				{ID: 2, Email: "org@x.com", Role: "organization", Password: "h"},
				{ID: 3, Email: "adm@x.com", Role: "admin", Password: "h"},
				{ID: 4, Email: "p1@x.com", Role: "partner"},              // password vazio → pendente
				{ID: 5, Email: "p2@x.com", Role: "partner", Password: "h"}, // ok
			}, nil
		},
	})
	out, err := uc.Execute(context.Background(), "a@x.com")
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if len(out.Users) != 2 {
		t.Fatalf("admin deve ver apenas parceiros, esperado 2, got %d", len(out.Users))
	}
	for _, u := range out.Users {
		if u.Role != "partner" {
			t.Fatalf("admin não deve ver role acima ou igual à sua, achou %s", u.Role)
		}
	}
	if !out.Users[0].PendingFirstAccess {
		t.Fatal("usuário com password vazio deve vir PendingFirstAccess=true")
	}
	if out.Users[1].PendingFirstAccess {
		t.Fatal("usuário com password definido deve vir PendingFirstAccess=false")
	}
}

func TestListUser_OrganizationVeTodos(t *testing.T) {
	uc := NewListUserUseCase(&mockUserRepoList{
		mockUserRepo: *requesterRepo(RoleOrganization, 1),
		listFunc: func(ctx context.Context) ([]*entities.User, error) {
			return []*entities.User{
				{ID: 2, Email: "org2@x.com", Role: "organization", Password: "h"},
				{ID: 3, Email: "adm@x.com", Role: "admin", Password: "h"},
				{ID: 4, Email: "p@x.com", Role: "partner"},
			}, nil
		},
	})
	out, err := uc.Execute(context.Background(), "o@x.com")
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if len(out.Users) != 3 {
		t.Fatalf("organization deve ver todos os usuários, esperado 3, got %d", len(out.Users))
	}
}

func TestListUser_RequesterInexistenteNaoParanica(t *testing.T) {
	uc := NewListUserUseCase(&mockUserRepo{findByEmail: func(ctx context.Context, email string) (*entities.User, error) {
		return nil, nil // gorm Find: usuário deletado com JWT ainda válido
	}})
	_, err := uc.Execute(context.Background(), "fantasma@x.com")
	if err == nil {
		t.Fatal("requester inexistente deve retornar erro, não pânico")
	}
}
