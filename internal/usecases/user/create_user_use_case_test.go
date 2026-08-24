package user

import (
	"context"
	"testing"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
)

type mockUserRepoCreate struct {
	mockUserRepo
	createFunc func(ctx context.Context, entity *entities.User) error
}

func (m *mockUserRepoCreate) CreateUser(ctx context.Context, entity *entities.User) error {
	return m.createFunc(ctx, entity)
}

type mockTokenRepoCreate struct {
	mockTokenRepoFull
	createFunc func(ctx context.Context, entity *entities.PasswordResetToken) error
}

func (m *mockTokenRepoCreate) CreateToken(ctx context.Context, entity *entities.PasswordResetToken) error {
	if m.createFunc == nil {
		return nil
	}
	return m.createFunc(ctx, entity)
}

func noopTokenRepo() *mockTokenRepoCreate {
	return &mockTokenRepoCreate{
		mockTokenRepoFull: mockTokenRepoFull{
			findFunc: func(ctx context.Context, hash string) (*entities.PasswordResetToken, error) { return nil, nil },
			markFunc: func(ctx context.Context, id int64) error { return nil },
		},
	}
}

func TestCreateUser_AdminNaoCriaAdmin(t *testing.T) {
	uc := NewCreateUserUseCase(requesterRepo(RoleAdmin, 1), noopTokenRepo(), &mockEmailSender{})
	_, err := uc.Execute(context.Background(), "a@x.com", &input.CreateUser{Name: "N", Email: "n@x.com", Role: "admin"})
	if err == nil {
		t.Fatal("admin não pode criar admin")
	}
}

func TestCreateUser_PartnerNaoCria(t *testing.T) {
	uc := NewCreateUserUseCase(requesterRepo(RolePartner, 1), noopTokenRepo(), &mockEmailSender{})
	_, err := uc.Execute(context.Background(), "p@x.com", &input.CreateUser{Name: "N", Email: "n@x.com", Role: "partner"})
	if err == nil {
		t.Fatal("partner não pode criar usuários")
	}
}

func TestCreateUser_EmailDuplicado(t *testing.T) {
	uc := NewCreateUserUseCase(
		&mockUserRepo{findByEmail: func(ctx context.Context, email string) (*entities.User, error) {
			if email == "dup@x.com" {
				return &entities.User{ID: 9, Email: email}, nil // já existe
			}
			return &entities.User{ID: 1, Email: email, Role: RoleOrganization}, nil // requester
		}},
		noopTokenRepo(),
		&mockEmailSender{},
	)
	_, err := uc.Execute(context.Background(), "o@x.com", &input.CreateUser{Name: "N", Email: "dup@x.com", Role: "partner"})
	if err == nil {
		t.Fatal("email duplicado deve falhar")
	}
}

func TestCreateUser_OrganizationCriaAdminComPasswordVazioEEnviaEmail(t *testing.T) {
	var saved *entities.User
	sender := &mockEmailSender{}
	uc := NewCreateUserUseCase(
		&mockUserRepoCreate{
			mockUserRepo: mockUserRepo{findByEmail: func(ctx context.Context, email string) (*entities.User, error) {
				if email == "novo@x.com" {
					return &entities.User{}, nil // ainda não existe
				}
				return &entities.User{ID: 1, Email: email, Role: RoleOrganization}, nil // requester
			}},
			createFunc: func(ctx context.Context, e *entities.User) error { e.ID = 42; saved = e; return nil },
		},
		noopTokenRepo(),
		sender,
	)
	out, err := uc.Execute(context.Background(), "o@x.com", &input.CreateUser{Name: "Novo", Email: "novo@x.com", Role: "admin"})
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if out.UserID != 42 {
		t.Fatalf("esperado UserID 42, got %d", out.UserID)
	}
	if saved.Password != "" {
		t.Fatal("password deve nascer vazio (primeiro acesso)")
	}
	if sender.sentTo != "novo@x.com" {
		t.Fatalf("esperado email para novo@x.com, got %s", sender.sentTo)
	}
}
