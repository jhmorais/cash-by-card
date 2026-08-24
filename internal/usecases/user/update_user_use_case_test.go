package user

import (
	"context"
	"testing"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
)

type mockUserRepoUpdate struct {
	mockUserRepo
	byID    func(ctx context.Context, id int) (*entities.User, error)
	updatef func(ctx context.Context, entity *entities.User) error
}

func (m *mockUserRepoUpdate) FindUserByID(ctx context.Context, id int) (*entities.User, error) {
	return m.byID(ctx, id)
}
func (m *mockUserRepoUpdate) UpdateUser(ctx context.Context, entity *entities.User) error {
	return m.updatef(ctx, entity)
}

func TestUpdateUser_AdminNaoEditaAdmin(t *testing.T) {
	uc := NewUpdateUserUseCase(&mockUserRepoUpdate{
		mockUserRepo: *requesterRepo(RoleAdmin, 1),
		byID: func(ctx context.Context, id int) (*entities.User, error) {
			return &entities.User{ID: 5, Email: "t@x.com", Role: "admin", Name: "Velho"}, nil
		},
		updatef: func(ctx context.Context, e *entities.User) error { return nil },
	})
	_, err := uc.Execute(context.Background(), "a@x.com", &input.UpdateUser{ID: 5, Name: "N", Role: "admin"})
	if err == nil {
		t.Fatal("admin não pode editar admin")
	}
}

func TestUpdateUser_NaoPodeEditarASiMesmo(t *testing.T) {
	uc := NewUpdateUserUseCase(&mockUserRepoUpdate{
		mockUserRepo: *requesterRepo(RoleOrganization, 5),
		byID: func(ctx context.Context, id int) (*entities.User, error) {
			return &entities.User{ID: 5, Email: "mesmo@x.com", Role: "admin"}, nil
		},
		updatef: func(ctx context.Context, e *entities.User) error { return nil },
	})
	_, err := uc.Execute(context.Background(), "mesmo@x.com", &input.UpdateUser{ID: 5, Name: "N", Role: "admin"})
	if err == nil {
		t.Fatal("não pode editar a si mesmo")
	}
}

func TestUpdateUser_OrganizationEditaAdmin(t *testing.T) {
	var saved *entities.User
	uc := NewUpdateUserUseCase(&mockUserRepoUpdate{
		mockUserRepo: *requesterRepo(RoleOrganization, 1),
		byID: func(ctx context.Context, id int) (*entities.User, error) {
			return &entities.User{ID: 5, Email: "t@x.com", Role: "admin", Name: "Velho"}, nil
		},
		updatef: func(ctx context.Context, e *entities.User) error { saved = e; return nil },
	})
	out, err := uc.Execute(context.Background(), "o@x.com", &input.UpdateUser{ID: 5, Name: "Novo", Role: "admin"})
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if out.UserID != 5 || saved.Name != "Novo" {
		t.Fatalf("update incorreto: %+v", saved)
	}
}

func TestUpdateUser_OrganizationPromovePartnerParaAdmin(t *testing.T) {
	var saved *entities.User
	uc := NewUpdateUserUseCase(&mockUserRepoUpdate{
		mockUserRepo: *requesterRepo(RoleOrganization, 1),
		byID: func(ctx context.Context, id int) (*entities.User, error) {
			return &entities.User{ID: 7, Email: "p@x.com", Role: "partner", Name: "Velho"}, nil
		},
		updatef: func(ctx context.Context, e *entities.User) error { saved = e; return nil },
	})
	_, err := uc.Execute(context.Background(), "o@x.com", &input.UpdateUser{ID: 7, Name: "Novo", Role: "admin"})
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if saved.Role != "admin" || saved.Name != "Novo" {
		t.Fatalf("esperado promoção para admin, got %+v", saved)
	}
}

func TestUpdateUser_AdminEditaPartner(t *testing.T) {
	uc := NewUpdateUserUseCase(&mockUserRepoUpdate{
		mockUserRepo: *requesterRepo(RoleAdmin, 2),
		byID: func(ctx context.Context, id int) (*entities.User, error) {
			return &entities.User{ID: 8, Email: "p2@x.com", Role: "partner", Name: "Velho"}, nil
		},
		updatef: func(ctx context.Context, e *entities.User) error { return nil },
	})
	out, err := uc.Execute(context.Background(), "a@x.com", &input.UpdateUser{ID: 8, Name: "Novo", Role: "partner"})
	if err != nil || out.UserID != 8 {
		t.Fatalf("esperado sucesso, err='%v'", err)
	}
}
