package user

import (
	"context"
	"strings"
	"testing"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	repoToken "github.com/jhmorais/cash-by-card/internal/repositories/token"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
)

// ---- mocks compartilhados por todos os testes do pacote user ----

type mockUserRepo struct {
	repositories.UserRepository
	findByEmail func(ctx context.Context, email string) (*entities.User, error)
}

func (m *mockUserRepo) FindUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	return m.findByEmail(ctx, email)
}

type mockTokenRepo struct {
	repoToken.PasswordResetTokenRepository
	createFunc func(ctx context.Context, entity *entities.PasswordResetToken) error
}

func (m *mockTokenRepo) CreateToken(ctx context.Context, entity *entities.PasswordResetToken) error {
	return m.createFunc(ctx, entity)
}

type mockEmailSender struct {
	sentTo   string
	sentLink string
}

func (m *mockEmailSender) SendPasswordResetEmail(ctx context.Context, to, link string) error {
	m.sentTo = to
	m.sentLink = link
	return nil
}

// ---- testes ----

func TestForgotPassword_EmailInexistenteNaoGeraToken(t *testing.T) {
	tokenRepo := &mockTokenRepo{createFunc: func(ctx context.Context, e *entities.PasswordResetToken) error {
		t.Fatal("não deve criar token para email inexistente")
		return nil
	}}
	uc := NewForgotPasswordUseCase(
		&mockUserRepo{findByEmail: func(ctx context.Context, email string) (*entities.User, error) {
			return &entities.User{}, nil
		}},
		tokenRepo,
		&mockEmailSender{},
	)
	err := uc.Execute(context.Background(), &input.ForgotPassword{Email: "ninguem@x.com"})
	if err != nil {
		t.Fatalf("email inexistente não deve retornar erro, got '%v'", err)
	}
}

func TestForgotPassword_EmailValidoGeraTokenEEnvia(t *testing.T) {
	var saved *entities.PasswordResetToken
	sender := &mockEmailSender{}
	tokenRepo := &mockTokenRepo{createFunc: func(ctx context.Context, e *entities.PasswordResetToken) error {
		saved = e
		return nil
	}}
	uc := NewForgotPasswordUseCase(
		&mockUserRepo{findByEmail: func(ctx context.Context, email string) (*entities.User, error) {
			return &entities.User{ID: 7, Email: email}, nil
		}},
		tokenRepo,
		sender,
	)
	if err := uc.Execute(context.Background(), &input.ForgotPassword{Email: "a@b.com"}); err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if saved == nil || saved.UserID != 7 {
		t.Fatalf("esperado token salvo para user 7, got %+v", saved)
	}
	if len(saved.TokenHash) != 64 || saved.ExpiresAt.IsZero() {
		t.Fatalf("token salvo inválido: %+v", saved)
	}
	if sender.sentTo != "a@b.com" || !strings.Contains(sender.sentLink, "token=") {
		t.Fatalf("email não enviado corretamente: to=%s link=%s", sender.sentTo, sender.sentLink)
	}
	if strings.Contains(sender.sentLink, saved.TokenHash) {
		t.Fatal("o link não deve conter o hash; deve conter o token em texto puro")
	}
}
