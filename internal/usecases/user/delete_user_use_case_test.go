package user

import (
	"context"
	"testing"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
)

type mockUserRepoDelete struct {
	mockUserRepo
	byID    func(ctx context.Context, id int) (*entities.User, error)
	deletef func(ctx context.Context, entity *entities.User) error
}

func (m *mockUserRepoDelete) FindUserByID(ctx context.Context, id int) (*entities.User, error) {
	return m.byID(ctx, id)
}
func (m *mockUserRepoDelete) DeleteUser(ctx context.Context, entity *entities.User) error {
	return m.deletef(ctx, entity)
}

type mockTokenRepoDelete struct {
	mockTokenRepo
	deleteByUserFunc func(ctx context.Context, userID int) error
}

func (m *mockTokenRepoDelete) DeleteTokensByUser(ctx context.Context, userID int) error {
	return m.deleteByUserFunc(ctx, userID)
}

func TestDeleteUser_NaoPodeExcluirASiMesmo(t *testing.T) {
	uc := NewDeleteUserUseCase(&mockUserRepoDelete{
		mockUserRepo: *requesterRepo(RoleOrganization, 1),
		byID: func(ctx context.Context, id int) (*entities.User, error) {
			return &entities.User{ID: 1, Email: "o@x.com", Role: RoleOrganization}, nil
		},
		deletef: func(ctx context.Context, e *entities.User) error { return nil },
	}, nil)
	_, err := uc.Execute(context.Background(), "o@x.com", 1)
	if err == nil {
		t.Fatal("não pode excluir a si mesmo")
	}
}

func TestDeleteUser_AdminNaoExcluiAdmin(t *testing.T) {
	uc := NewDeleteUserUseCase(&mockUserRepoDelete{
		mockUserRepo: *requesterRepo(RoleAdmin, 2),
		byID: func(ctx context.Context, id int) (*entities.User, error) {
			return &entities.User{ID: 5, Role: "admin", Email: "t@x.com"}, nil
		},
		deletef: func(ctx context.Context, e *entities.User) error { return nil },
	}, nil)
	_, err := uc.Execute(context.Background(), "a@x.com", 5)
	if err == nil {
		t.Fatal("admin não pode excluir admin")
	}
}

func TestDeleteUser_OrganizationExcluiPartnerApagandoTokensPrimeiro(t *testing.T) {
	deleted := false
	tokensDeletedFor := 0
	order := []string{}
	uc := NewDeleteUserUseCase(&mockUserRepoDelete{
		mockUserRepo: *requesterRepo(RoleOrganization, 1),
		byID: func(ctx context.Context, id int) (*entities.User, error) {
			return &entities.User{ID: 6, Role: "partner", Email: "p@x.com"}, nil
		},
		deletef: func(ctx context.Context, e *entities.User) error {
			deleted = true
			order = append(order, "user")
			return nil
		},
	}, &mockTokenRepoDelete{deleteByUserFunc: func(ctx context.Context, userID int) error {
		tokensDeletedFor = userID
		order = append(order, "tokens")
		return nil
	}})
	out, err := uc.Execute(context.Background(), "o@x.com", 6)
	if err != nil || !out.Success || !deleted {
		t.Fatalf("esperado sucesso, err='%v' success=%v deleted=%v", err, out.Success, deleted)
	}
	if tokensDeletedFor != 6 {
		t.Fatalf("esperado apagar tokens do user 6, got %d", tokensDeletedFor)
	}
	if len(order) != 2 || order[0] != "tokens" || order[1] != "user" {
		t.Fatalf("tokens devem ser apagados ANTES do user, ordem=%v", order)
	}
}
