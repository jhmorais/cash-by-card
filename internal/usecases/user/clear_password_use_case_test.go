package user

import (
	"context"
	"testing"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
)

type mockUserRepoClear struct {
	mockUserRepo
	byID   func(ctx context.Context, id int) (*entities.User, error)
	clearf func(ctx context.Context, id int) error
}

func (m *mockUserRepoClear) FindUserByID(ctx context.Context, id int) (*entities.User, error) {
	return m.byID(ctx, id)
}
func (m *mockUserRepoClear) ClearUserPassword(ctx context.Context, id int) error {
	return m.clearf(ctx, id)
}

func TestClearPassword_AdminSoParceiro(t *testing.T) {
	uc := NewClearPasswordUseCase(
		&mockUserRepoClear{
			mockUserRepo: *requesterRepo(RoleAdmin, 2),
			byID: func(ctx context.Context, id int) (*entities.User, error) {
				return &entities.User{ID: 5, Role: "admin", Email: "t@x.com"}, nil
			},
			clearf: func(ctx context.Context, id int) error { return nil },
		},
		noopTokenRepo(),
		&mockEmailSender{},
	)
	_, err := uc.Execute(context.Background(), "a@x.com", 5)
	if err == nil {
		t.Fatal("admin não pode limpar senha de admin")
	}
}

func TestClearPassword_NaoPodeLimparAPropria(t *testing.T) {
	uc := NewClearPasswordUseCase(
		&mockUserRepoClear{
			mockUserRepo: *requesterRepo(RoleOrganization, 1),
			byID: func(ctx context.Context, id int) (*entities.User, error) {
				return &entities.User{ID: 1, Email: "o@x.com", Role: RoleOrganization}, nil
			},
			clearf: func(ctx context.Context, id int) error { return nil },
		},
		noopTokenRepo(),
		&mockEmailSender{},
	)
	_, err := uc.Execute(context.Background(), "o@x.com", 1)
	if err == nil {
		t.Fatal("não pode limpar a própria senha por aqui")
	}
}

func TestClearPassword_OrganizationLimpaEEnvia(t *testing.T) {
	cleared := false
	sender := &mockEmailSender{}
	uc := NewClearPasswordUseCase(
		&mockUserRepoClear{
			mockUserRepo: *requesterRepo(RoleOrganization, 1),
			byID: func(ctx context.Context, id int) (*entities.User, error) {
				return &entities.User{ID: 6, Role: "partner", Email: "p@x.com", Name: "P"}, nil
			},
			clearf: func(ctx context.Context, id int) error { cleared = true; return nil },
		},
		noopTokenRepo(),
		sender,
	)
	out, err := uc.Execute(context.Background(), "o@x.com", 6)
	if err != nil || out.UserID != 6 || !cleared {
		t.Fatalf("esperado sucesso, err='%v'", err)
	}
	if sender.sentTo != "p@x.com" {
		t.Fatalf("esperado email para p@x.com, got %s", sender.sentTo)
	}
}
