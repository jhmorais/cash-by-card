# User Access & Password Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implementar primeiro acesso, recuperação de senha por email com token de uso único, e página de administração de usuários com roles organization/admin/partner.

**Architecture:** Backend Go (clean-ish: contracts → use cases → repositories, DI no boot) adiciona tabela `password_reset_token`, interface `EmailSender` (SMTP com fallback para log), e use cases com checagem de permissão por requester. Frontend React+MUI adiciona páginas públicas/privadas e nav por role.

**Tech Stack:** Go 1.24, gorm/mysql, golang-jwt, net/smtp (stdlib, sem dependência nova) | React 18, MUI, axios-free fetch wrapper, vitest + @testing-library/react.

**Spec:** `docs/superpowers/specs/2026-08-24-user-access-password-design.md`

**Paths:** backend = `/home/gabigontijo/Documents/cash-by-card`; frontend = `/home/gabigontijo/Documents/cash-by-card-front`. Branch `feature/user-access-password` já criada em ambos (a partir da main). **Nunca commitar `.env`** (gitignored).

**Convenções do backend (seguir à risca):**
- Testes: `go test ./internal/usecases/... -run <Teste> -v` e suite completa `go test ./...`
- Contrato: arquivo `internal/contracts/i<nome>_use_case.go` com interface `Execute(...)`
- Use case: `internal/usecases/<dominio>/<nome>_use_case.go`, struct privada + `New<Nome>UseCase(deps) contracts.<Nome>UseCase`
- Mocks de repo: embed da interface + campo Func (padrão em `internal/usecases/loan/list_loan_use_case_test.go`)
- Handlers: `services/*.go`, método em `(h *Handler)`, sucesso com `json.Marshal` + `fmt.Fprint`, erro com `utils.WriteErrModel`
- Roles (strings exatas): `"organization"`, `"admin"`, `"partner"`

**Convenções do frontend:**
- Testes: `npm test` (vitest run); fetch mock global com `vi.stubGlobal('fetch', ...)` + `vi.stubEnv('VITE_API_BASE_URL', ...)`
- Página = wrapper em `src/pages/<nome>.jsx` + view em `src/sections/<nome>/view/<nome>-view.jsx`
- APIs: `src/apis/<nome>/index.js` usando `apiFetch` de `src/apis`

---

## File Structure (mapa de arquivos)

**Backend (criar/modificar):**
- Criar: `db/migrations/2026-08-24-user-access-password.sql`
- Modificar: `docker/mysql/dump.sql`, `internal/domain/entities/user.go` (novo entity PasswordResetToken em arquivo próprio: `internal/domain/entities/password_reset_token.go`)
- Criar: `internal/repositories/token/ipassword_reset_token_repository.go`, `internal/repositories/token/password_reset_token_repository.go`
- Modificar: `internal/repositories/user/iuser_repository.go`, `internal/repositories/user/user_repository.go`
- Criar: `internal/infra/email/email_sender.go`
- Criar: `internal/contracts/email.go` + `iforgot_password_use_case.go`, `ireset_password_use_case.go`, `ichange_password_use_case.go`, `iget_user_use_case.go`, `iupdate_user_use_case.go`, `idelete_user_use_case.go`, `iclear_password_use_case.go`
- Modificar: `internal/contracts/ilist_user_use_case.go`, `icreate_user_use_case.go`
- Criar inputs: `internal/ports/input/user/forgot_password.go`, `reset_password.go`, `change_password.go`, `update_user.go`
- Modificar: `internal/ports/input/user/create_user.go` (remove Password)
- Criar outputs: `internal/ports/output/user/list_user.go`, `get_user.go`, `update_user.go`, `delete_user.go`, `clear_password.go`
- Criar use cases: `internal/usecases/user/forgot_password_use_case.go` (+_test), `reset_password_use_case.go` (+_test), `change_password_use_case.go` (+_test), `get_user_use_case.go` (+_test), `update_user_use_case.go` (+_test), `delete_user_use_case.go` (+_test), `clear_password_use_case.go` (+_test)
- Modificar: `internal/usecases/user/create_user_use_case.go` (+_test), `list_user_use_case.go` (+_test), `internal/usecases/login/login_use_case.go` (+_test)
- Criar: `utils/reset_token.go` (+_test), modificar `utils/rest.go` (EmailFromContext + RoleMiddleware variadic)
- Modificar: `internal/infra/di/dependency_builder.go`, `services/user.go`, `services/rest_server.go`; criar `services/user_account.go`
- Modificar: `.env` (local, não commitar), `config/config.go` (GetFrontURL)

**Frontend (criar/modificar):**
- Criar: `src/apis/user/index.js` (+_test)
- Modificar: `src/hooks/authProvider.jsx`, `src/routes/sections.jsx`, `src/layouts/dashboard/config-navigation.jsx`, `src/layouts/dashboard/nav.jsx`, `src/sections/login/login-view.jsx`
- Criar: `src/pages/primeiro-acesso.jsx`, `src/pages/alterar-senha.jsx`, `src/pages/usuario.jsx`
- Criar: `src/sections/first-access/view/first-access-view.jsx` (+_test), `src/sections/change-password/view/change-password-view.jsx`, `src/sections/user/view/usuario-view.jsx`

---

# BACKEND

### Task 1: Migration SQL + entity + dump.sql

**Files:**
- Create: `db/migrations/2026-08-24-user-access-password.sql`
- Create: `internal/domain/entities/password_reset_token.go`
- Modify: `docker/mysql/dump.sql`
- Modify (banco local rodando): aplicar via docker exec

- [ ] **Step 1: Escrever o script de migration**

`db/migrations/2026-08-24-user-access-password.sql`:

```sql
-- 2026-08-24: primeiro acesso e administracao de usuarios
-- password NULL = usuario pendente de primeiro acesso
ALTER TABLE `user` MODIFY `password` VARCHAR(100) NULL;

CREATE TABLE `password_reset_token` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` INT NOT NULL,
  `token_hash` VARCHAR(64) NOT NULL,
  `expires_at` TIMESTAMP NOT NULL,
  `used_at` TIMESTAMP NULL DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`),
  KEY `password_reset_token_user_id` (`user_id`),
  CONSTRAINT `password_reset_token_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
```

- [ ] **Step 2: Atualizar o dump.sql**

Em `docker/mysql/dump.sql`, trocar a definição de `user` (coluna `password` passa a aceitar NULL):

```sql
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `email` varchar(100) NOT NULL,
  `name` varchar(100) NOT NULL,
  `password` varchar(100) DEFAULT NULL,
  `role` varchar(100) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_unique` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
```

E adicionar após a tabela `user`, antes de `SET FOREIGN_KEY_CHECKS = 1;`, a mesma `CREATE TABLE password_reset_token` do migration (Step 1).

- [ ] **Step 3: Criar a entity**

`internal/domain/entities/password_reset_token.go`:

```go
package entities

import "time"

// PasswordResetToken guarda o hash SHA-256 (hex) do token enviado por email.
// O token em texto puro nunca é persistido.
type PasswordResetToken struct {
	ID        int64      `gorm:"id" json:"id"`
	UserID    int        `json:"user_id"`
	TokenHash string     `json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at"`
	CreatedAt time.Time  `json:"created_at"`
}
```

(gorm com `SingularTable: true` mapeia para a tabela `password_reset_token`.)

- [ ] **Step 4: Aplicar a migration no banco local**

```bash
sg docker -c "docker exec -i mysqlcontainer mysql -uroot -ppassword database" < db/migrations/2026-08-24-user-access-password.sql
sg docker -c "docker exec mysqlcontainer mysql -uroot -ppassword database -e 'SHOW CREATE TABLE password_reset_token\G' " | head -5
```

Expected: tabela criada sem erro.

- [ ] **Step 5: Commit**

```bash
git add db/migrations/ docker/mysql/dump.sql internal/domain/entities/password_reset_token.go
git commit -m "feat: password_reset_token entity e migration de password nullable"
```

---

### Task 2: utils — token generator/hash e EmailFromContext (TDD)

**Files:**
- Create: `utils/reset_token.go`
- Test: `utils/reset_token_test.go`
- Modify: `utils/rest.go`

- [ ] **Step 1: Teste falhando**

`utils/reset_token_test.go`:

```go
package utils

import (
	"encoding/sha256"
	"encoding/hex"
	"testing"
)

func TestGenerateResetToken(t *testing.T) {
	plain, hash, err := GenerateResetToken()
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if len(plain) != 64 {
		t.Fatalf("expected token with 64 hex chars, got %d", len(plain))
	}
	sum := sha256.Sum256([]byte(plain))
	if hash != hex.EncodeToString(sum[:]) {
		t.Fatal("expected hash to be the SHA-256 hex of the token")
	}
	plain2, _, _ := GenerateResetToken()
	if plain == plain2 {
		t.Fatal("expected two generated tokens to differ")
	}
}
```

- [ ] **Step 2: Rodar e verificar falha**

Run: `go test ./utils/ -run TestGenerateResetToken -v`
Expected: FAIL `undefined: GenerateResetToken`

- [ ] **Step 3: Implementar**

`utils/reset_token.go`:

```go
package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateResetToken devolve (tokenEmTexto, hashSHA256Hex, err).
// O texto puro vai no email; o hash vai para o banco.
func GenerateResetToken() (string, string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}
	plain := hex.EncodeToString(buf)
	return plain, HashResetToken(plain), nil
}

func HashResetToken(plain string) string {
	sum := sha256Sum(plain)
	return hex.EncodeToString(sum)
}
```

Nota: `sha256Sum` não existe ainda — implementar direto:

```go
func HashResetToken(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}
```

(imports: `crypto/sha256`)

- [ ] **Step 4: Rodar e verificar passagem**

Run: `go test ./utils/ -v`
Expected: PASS

- [ ] **Step 5: Adicionar EmailFromContext em utils/rest.go**

No final de `utils/rest.go` (o `emailKey` já existe no arquivo):

```go
// EmailFromContext devolve o email autenticado pelo ValidateJwtTokenMiddleware.
func EmailFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(emailKey).(string); ok {
		return v
	}
	return ""
}
```

- [ ] **Step 6: Commit**

```bash
git add utils/
git commit -m "feat: gerador de token de reset e EmailFromContext"
```

---

### Task 3: Repositories — token e métodos de senha do user

**Files:**
- Create: `internal/repositories/token/ipassword_reset_token_repository.go`, `internal/repositories/token/password_reset_token_repository.go`
- Modify: `internal/repositories/user/iuser_repository.go`, `internal/repositories/user/user_repository.go`

(Sem unit test de repo — o projeto não tem harness para gorm; coberto pelos testes de use case com mock e pela verificação ao vivo da Task 13.)

- [ ] **Step 1: Interface e impl do token repo**

`internal/repositories/token/ipassword_reset_token_repository.go`:

```go
package repositories

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
)

type PasswordResetTokenRepository interface {
	CreateToken(ctx context.Context, entity *entities.PasswordResetToken) error
	// FindValidTokenByHash devolve nil, nil quando não há token válido (hash inexistente, expirado ou já usado).
	FindValidTokenByHash(ctx context.Context, tokenHash string) (*entities.PasswordResetToken, error)
	MarkTokenUsed(ctx context.Context, id int64) error
}
```

`internal/repositories/token/password_reset_token_repository.go`:

```go
package repositories

import (
	"context"
	"time"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	"gorm.io/gorm"
)

type passwordResetTokenRepository struct {
	db *gorm.DB
}

func NewPasswordResetTokenRepository(db *gorm.DB) PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

func (r *passwordResetTokenRepository) CreateToken(ctx context.Context, entity *entities.PasswordResetToken) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *passwordResetTokenRepository) FindValidTokenByHash(ctx context.Context, tokenHash string) (*entities.PasswordResetToken, error) {
	var entity *entities.PasswordResetToken
	err := r.db.WithContext(ctx).
		Where("token_hash = ? AND used_at IS NULL AND expires_at > ?", tokenHash, time.Now()).
		Limit(1).
		Find(&entity).Error
	if err != nil {
		return nil, err
	}
	if entity == nil || entity.ID == 0 {
		return nil, nil
	}
	return entity, nil
}

func (r *passwordResetTokenRepository) MarkTokenUsed(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&entities.PasswordResetToken{}).
		Where("id = ?", id).
		Update("used_at", time.Now()).Error
}
```

- [ ] **Step 2: Métodos de senha no UserRepository**

Adicionar à interface em `internal/repositories/user/iuser_repository.go`:

```go
	// SetUserPassword grava o hash da senha (map update — permite valor vazio).
	SetUserPassword(ctx context.Context, id int, hashedPassword string) error
	// ClearUserPassword seta password = NULL (usuário volta a primeiro acesso).
	ClearUserPassword(ctx context.Context, id int) error
```

Adicionar à impl em `internal/repositories/user/user_repository.go`:

```go
func (d *userRepository) SetUserPassword(ctx context.Context, id int, hashedPassword string) error {
	return d.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("id = ?", id).
		Update("password", hashedPassword).Error
}

func (d *userRepository) ClearUserPassword(ctx context.Context, id int) error {
	return d.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("id = ?", id).
		Update("password", gorm.Expr("NULL")).Error
}
```

- [ ] **Step 3: Build**

Run: `go build ./...`
Expected: sem erro

- [ ] **Step 4: Commit**

```bash
git add internal/repositories/
git commit -m "feat: repositorio de password_reset_token e metodos de senha do user"
```

---

### Task 4: EmailSender (contrato + SMTP + fallback log)

**Files:**
- Create: `internal/contracts/email.go`
- Create: `internal/infra/email/email_sender.go`
- Modify: `config/config.go`, `.env` (local)

- [ ] **Step 1: Contrato**

`internal/contracts/email.go`:

```go
package contracts

import "context"

// EmailSender envia emails transacionais. Implementação SMTP em internal/infra/email.
type EmailSender interface {
	SendPasswordResetEmail(ctx context.Context, to, resetLink string) error
}
```

- [ ] **Step 2: Implementações**

`internal/infra/email/email_sender.go`:

```go
package email

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/jhmorais/cash-by-card/config"
)

type smtpSender struct {
	host, port, user, password, from string
}

type logSender struct{}

// NewSenderFromEnv devolve SmtpSender se SMTP_HOST estiver configurado;
// senão devolve o fallback que apenas loga o link (dev).
func NewSenderFromEnv() contracts.EmailSender {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		return &logSender{}
	}
	return &smtpSender{
		host:     host,
		port:     envOr("SMTP_PORT", "587"),
		user:     os.Getenv("SMTP_USER"),
		password: os.Getenv("SMTP_PASSWORD"),
		from:     envOr("SMTP_FROM", os.Getenv("SMTP_USER")),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func (s *smtpSender) SendPasswordResetEmail(ctx context.Context, to, resetLink string) error {
	subject := "Cash By Card - definicao de senha"
	body := fmt.Sprintf(
		"Ola,\n\nUse o link abaixo para definir sua senha. Ele expira em 30 minutos e pode ser usado uma unica vez:\n\n%s\n\nSe voce nao solicitou, ignore este email.\n",
		resetLink,
	)
	msg := strings.Join([]string{
		"From: " + s.from,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	addr := s.host + ":" + s.port
	auth := smtp.PlainAuth("", s.user, s.password, s.host)
	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}

func (s *logSender) SendPasswordResetEmail(ctx context.Context, to, resetLink string) error {
	log.Printf("[DEV-EMAIL] para=%s link=%s", to, resetLink)
	return nil
}
```

(Adicionar import `"github.com/jhmorais/cash-by-card/internal/contracts"`.)

- [ ] **Step 3: GetFrontURL no config**

Adicionar a `config/config.go`:

```go
func GetFrontURL() string {
	url := viper.GetString("FRONT_URL")
	if url == "" {
		return "http://localhost:5173"
	}
	return url
}
```

- [ ] **Step 4: .env local (não commitar)**

Acrescentar ao `.env` da raiz:

```
FRONT_URL=http://localhost:5173
# SMTP (preencher quando tiver credenciais; vazio = loga o link no console)
SMTP_HOST=
SMTP_PORT=587
SMTP_USER=
SMTP_PASSWORD=
SMTP_FROM=
```

- [ ] **Step 5: Build + commit**

Run: `go build ./...` → sem erro

```bash
git add internal/contracts/email.go internal/infra/email/ config/config.go
git commit -m "feat: EmailSender com SMTP e fallback de log"
```

---

### Task 5: Inputs e outputs novos

**Files:**
- Create: `internal/ports/input/user/forgot_password.go`, `reset_password.go`, `change_password.go`, `update_user.go`
- Modify: `internal/ports/input/user/create_user.go`
- Create: `internal/ports/output/user/list_user.go`, `get_user.go`, `update_user.go`, `delete_user.go`, `clear_password.go`

- [ ] **Step 1: Inputs**

`internal/ports/input/user/forgot_password.go`:

```go
package input

type ForgotPassword struct {
	Email string `json:"email"`
}
```

`internal/ports/input/user/reset_password.go`:

```go
package input

type ResetPassword struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}
```

`internal/ports/input/user/change_password.go`:

```go
package input

type ChangePassword struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}
```

`internal/ports/input/user/update_user.go`:

```go
package input

type UpdateUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}
```

`internal/ports/input/user/create_user.go` (reescrever — sem Password; a senha agora nasce pelo fluxo de primeiro acesso):

```go
package input

type CreateUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
```

- [ ] **Step 2: Outputs**

`internal/ports/output/user/list_user.go`:

```go
package output

type UserItem struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	Role               string `json:"role"`
	PendingFirstAccess bool   `json:"pendingFirstAccess"`
}

type ListUser struct {
	Users []UserItem `json:"users"`
}
```

`internal/ports/output/user/get_user.go`:

```go
package output

type GetUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}
```

`internal/ports/output/user/update_user.go`:

```go
package output

type UpdateUser struct {
	UserID int `json:"userId"`
}
```

`internal/ports/output/user/delete_user.go`:

```go
package output

type DeleteUser struct {
	Success bool `json:"success"`
}
```

`internal/ports/output/user/clear_password.go`:

```go
package output

type ClearPassword struct {
	UserID int `json:"userId"`
}
```

- [ ] **Step 3: Build + commit**

Run: `go build ./...` → o build de `create_user_use_case.go` quebra (Password removido). Corrigir o use case NESTE passo apenas o suficiente para compilar: remover as linhas `if len(createUser.Password) < 6 {...}` e `hashUser := utils.EncryptPassword(createUser.Password)` e `Password: hashUser,` da entity (a role check também será reescrita na Task 9; por agora trocar a validação de role para aceitar os três valores):

```go
	if createUser.Role != "organization" && createUser.Role != "admin" && createUser.Role != "partner" {
		return nil, fmt.Errorf("cannot create a user without valid role")
	}
```

e a entity:

```go
	userEntity := &entities.User{
		Name:      createUser.Name,
		Email:     createUser.Email,
		Role:      createUser.Role,
		CreatedAt: time.Now(),
	}
```

Run: `go build ./...` → sem erro

```bash
git add internal/ports/
git commit -m "feat: inputs e outputs de primeiro acesso e administracao de usuarios"
```

---

### Task 6: Contratos novos

**Files:**
- Create: `internal/contracts/iforgot_password_use_case.go`, `ireset_password_use_case.go`, `ichange_password_use_case.go`, `iget_user_use_case.go`, `iupdate_user_use_case.go`, `idelete_user_use_case.go`, `iclear_password_use_case.go`
- Modify: `internal/contracts/ilist_user_use_case.go`, `icreate_user_use_case.go`

- [ ] **Step 1: Escrever os contratos**

`iforgot_password_use_case.go`:

```go
package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
)

// ForgotPasswordUseCase envia email de definicao de senha se o email existir.
// Nunca retorna erro por email inexistente (anti-enumeração).
type ForgotPasswordUseCase interface {
	Execute(ctx context.Context, forgotPassword *input.ForgotPassword) error
}
```

`ireset_password_use_case.go`:

```go
package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
)

type ResetPasswordUseCase interface {
	Execute(ctx context.Context, resetPassword *input.ResetPassword) error
}
```

`ichange_password_use_case.go`:

```go
package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
)

type ChangePasswordUseCase interface {
	Execute(ctx context.Context, email string, changePassword *input.ChangePassword) error
}
```

`iget_user_use_case.go`:

```go
package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
)

type GetUserUseCase interface {
	Execute(ctx context.Context, email string) (*output.GetUser, error)
}
```

`iupdate_user_use_case.go`:

```go
package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
)

type UpdateUserUseCase interface {
	Execute(ctx context.Context, requesterEmail string, updateUser *input.UpdateUser) (*output.UpdateUser, error)
}
```

`idelete_user_use_case.go`:

```go
package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
)

type DeleteUserUseCase interface {
	Execute(ctx context.Context, requesterEmail string, id int) (*output.DeleteUser, error)
}
```

`iclear_password_use_case.go`:

```go
package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
)

type ClearPasswordUseCase interface {
	Execute(ctx context.Context, requesterEmail string, id int) (*output.ClearPassword, error)
}
```

- [ ] **Step 2: Ajustar contratos existentes**

`icreate_user_use_case.go` (adiciona requesterEmail):

```go
package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
)

type CreateUserUseCase interface {
	Execute(ctx context.Context, requesterEmail string, createUser *input.CreateUser) (*output.CreateUser, error)
}
```

`ilist_user_use_case.go` (adiciona requesterEmail):

```go
package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
)

type ListUserUseCase interface {
	Execute(ctx context.Context, requesterEmail string) (*output.ListUser, error)
}
```

- [ ] **Step 3: Build**

Run: `go build ./...`
Expected: quebra nos use cases existentes (create_user, list_user) e na DI — serão reescritos nas Tasks 7-9. Se preferir manter o build verde a cada commit, faça Task 6+7+8+9 antes do commit único final destas. Caso contrário, commit apenas dos contratos novos com os use cases ainda compilando é impossível — **commit apenas no fim da Task 9**.

---

### Task 7: Use cases de senha — forgot, reset, change (TDD)

**Files:**
- Create: `internal/usecases/user/forgot_password_use_case.go` (+_test), `reset_password_use_case.go` (+_test), `change_password_use_case.go` (+_test)

- [ ] **Step 1: Testes do forgot**

`internal/usecases/user/forgot_password_use_case_test.go`:

```go
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
```

(`mockUserRepo`, `mockTokenRepo` e `mockEmailSender` ficam neste arquivo e são reusados pelos testes das Tasks 7 e 9 — todo o pacote compartilha.)

- [ ] **Step 2: Rodar e verificar falha**

Run: `go test ./internal/usecases/user/ -run TestForgotPassword -v`
Expected: FAIL `undefined: NewForgotPasswordUseCase`

- [ ] **Step 3: Implementar forgot**

`internal/usecases/user/forgot_password_use_case.go`:

```go
package user

import (
	"context"
	"log"
	"time"

	"github.com/jhmorais/cash-by-card/config"
	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	repoToken "github.com/jhmorais/cash-by-card/internal/repositories/token"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
	"github.com/jhmorais/cash-by-card/utils"
)

type forgotPasswordUseCase struct {
	userRepository  repositories.UserRepository
	tokenRepository repoToken.PasswordResetTokenRepository
	emailSender     contracts.EmailSender
}

func NewForgotPasswordUseCase(
	userRepository repositories.UserRepository,
	tokenRepository repoToken.PasswordResetTokenRepository,
	emailSender contracts.EmailSender,
) contracts.ForgotPasswordUseCase {
	return &forgotPasswordUseCase{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		emailSender:     emailSender,
	}
}

func (f *forgotPasswordUseCase) Execute(ctx context.Context, forgotPassword *input.ForgotPassword) error {
	if forgotPassword.Email == "" {
		return nil // resposta genérica também para input vazio
	}
	user, err := f.userRepository.FindUserByEmail(ctx, forgotPassword.Email)
	if err != nil || user == nil || user.Email == "" {
		return nil // nunca revelar existencia do email
	}

	plain, hash, err := utils.GenerateResetToken()
	if err != nil {
		return err
	}

	token := &entities.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(30 * time.Minute),
		CreatedAt: time.Now(),
	}
	if err := f.tokenRepository.CreateToken(ctx, token); err != nil {
		return err
	}

	link := config.GetFrontURL() + "/primeiro-acesso?token=" + plain
	if err := f.emailSender.SendPasswordResetEmail(ctx, user.Email, link); err != nil {
		log.Printf("falha ao enviar email de reset para %s: %v", user.Email, err)
	}
	return nil
}
```

- [ ] **Step 4: Rodar forgot e verificar passagem**

Run: `go test ./internal/usecases/user/ -run TestForgotPassword -v`
Expected: PASS nos dois testes

- [ ] **Step 5: Testes do reset**

`internal/usecases/user/reset_password_use_case_test.go`:

```go
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
	err := uc.Execute(context.Background(), &input.ResetPassword{Token: "x", NewPassword: "123456"})
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
	if err := uc.Execute(context.Background(), &input.ResetPassword{Token: "plain-token", NewPassword: "novasenha"}); err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if markedID != 55 {
		t.Fatalf("esperado marcar token 55 como usado, got %d", markedID)
	}
	if savedID != 9 || savedHash != utils.EncryptPassword("novasenha") {
		t.Fatalf("senha não salva corretamente: id=%d hash=%s", savedID, savedHash)
	}
}
```

(`mockTokenRepoFull` e `mockUserRepoWithSet` ficam neste arquivo; são reusados pela Task 9.)

- [ ] **Step 6: Rodar e verificar falha**

Run: `go test ./internal/usecases/user/ -run TestResetPassword -v`
Expected: FAIL `undefined: NewResetPasswordUseCase`

- [ ] **Step 7: Implementar reset**

`internal/usecases/user/reset_password_use_case.go`:

```go
package user

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	repoToken "github.com/jhmorais/cash-by-card/internal/repositories/token"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
	"github.com/jhmorais/cash-by-card/utils"
)

type resetPasswordUseCase struct {
	tokenRepository repoToken.PasswordResetTokenRepository
	userRepository  repositories.UserRepository
}

func NewResetPasswordUseCase(
	tokenRepository repoToken.PasswordResetTokenRepository,
	userRepository repositories.UserRepository,
) contracts.ResetPasswordUseCase {
	return &resetPasswordUseCase{tokenRepository: tokenRepository, userRepository: userRepository}
}

func (r *resetPasswordUseCase) Execute(ctx context.Context, resetPassword *input.ResetPassword) error {
	if resetPassword.Token == "" {
		return fmt.Errorf("token é obrigatório")
	}
	if len(resetPassword.NewPassword) < 6 {
		return fmt.Errorf("a nova senha deve ter pelo menos 6 caracteres")
	}

	hash := utils.HashResetToken(resetPassword.Token)
	token, err := r.tokenRepository.FindValidTokenByHash(ctx, hash)
	if err != nil {
		return fmt.Errorf("failed to validate token: %v", err)
	}
	if token == nil {
		return fmt.Errorf("token inválido ou expirado")
	}

	if err := r.userRepository.SetUserPassword(ctx, token.UserID, utils.EncryptPassword(resetPassword.NewPassword)); err != nil {
		return fmt.Errorf("failed to save password: %v", err)
	}

	if err := r.tokenRepository.MarkTokenUsed(ctx, token.ID); err != nil {
		return fmt.Errorf("failed to invalidate token: %v", err)
	}
	return nil
}
```

- [ ] **Step 8: Rodar reset e verificar passagem**

Run: `go test ./internal/usecases/user/ -run TestResetPassword -v`
Expected: PASS

- [ ] **Step 9: Teste + implementação do change password**

`internal/usecases/user/change_password_use_case_test.go`:

```go
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
	err := uc.Execute(context.Background(), "a@b.com", &input.ChangePassword{CurrentPassword: "errada", NewPassword: "novasenha"})
	if err == nil {
		t.Fatal("senha atual errada deve retornar erro")
	}
}

func TestChangePassword_Sucesso(t *testing.T) {
	userRepo := &mockUserRepoWithSet{setFunc: func(ctx context.Context, id int, hash string) error { return nil }}
	uc := NewChangePasswordUseCase(&mockUserRepo{findByEmail: func(ctx context.Context, email string) (*entities.User, error) {
		return &entities.User{ID: 3, Email: email, Password: utils.EncryptPassword("correta")}, nil
	}})
	err := uc.Execute(context.Background(), "a@b.com", &input.ChangePassword{CurrentPassword: "correta", NewPassword: "novasenha"})
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
}
```

`internal/usecases/user/change_password_use_case.go`:

```go
package user

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
	"github.com/jhmorais/cash-by-card/utils"
)

type changePasswordUseCase struct {
	userRepository repositories.UserRepository
}

func NewChangePasswordUseCase(userRepository repositories.UserRepository) contracts.ChangePasswordUseCase {
	return &changePasswordUseCase{userRepository: userRepository}
}

func (c *changePasswordUseCase) Execute(ctx context.Context, email string, changePassword *input.ChangePassword) error {
	if len(changePassword.NewPassword) < 6 {
		return fmt.Errorf("a nova senha deve ter pelo menos 6 caracteres")
	}
	user, err := c.userRepository.FindUserByEmail(ctx, email)
	if err != nil || user == nil || user.Email == "" {
		return fmt.Errorf("usuário não encontrado")
	}
	if user.Password == "" || user.Password != utils.EncryptPassword(changePassword.CurrentPassword) {
		return fmt.Errorf("senha atual incorreta")
	}
	return c.userRepository.SetUserPassword(ctx, user.ID, utils.EncryptPassword(changePassword.NewPassword))
}
```

- [ ] **Step 10: Rodar e commit**

Run: `go test ./internal/usecases/user/ -v` → todos PASS

```bash
git add internal/usecases/user/
git commit -m "feat: use cases de forgot/reset/change password"
```

---

### Task 8: Login — erro de primeiro acesso (TDD)

**Files:**
- Modify: `internal/usecases/login/login_use_case.go`
- Test: `internal/usecases/login/login_use_case_test.go`

- [ ] **Step 1: Teste falhando**

Criar `internal/usecases/login/login_use_case_test.go`:

```go
package login

import (
	"context"
	"strings"
	"testing"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
)

type mockUserRepoLogin struct {
	repositories.UserRepository
	byEmail  func(ctx context.Context, email string) (*entities.User, error)
	byEmailP func(ctx context.Context, email, password string) (*entities.User, error)
}

func (m *mockUserRepoLogin) FindUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	return m.byEmail(ctx, email)
}
func (m *mockUserRepoLogin) FindUserByEmailandPassword(ctx context.Context, email, password string) (*entities.User, error) {
	return m.byEmailP(ctx, email, password)
}

func TestLogin_UsuarioPendentePrimeiroAcesso(t *testing.T) {
	uc := NewLoginUseCase(&mockUserRepoLogin{
		byEmail: func(ctx context.Context, email string) (*entities.User, error) {
			return &entities.User{ID: 1, Email: email, Password: ""}, nil // password NULL no banco
		},
		byEmailP: func(ctx context.Context, email, password string) (*entities.User, error) {
			return &entities.User{}, nil
		},
	})
	_, err := uc.Execute(context.Background(), &input.UserLogin{Email: "a@b.com", Password: "qualquer"})
	if err == nil || !strings.Contains(err.Error(), "primeiro acesso") {
		t.Fatalf("esperado erro de primeiro acesso, got '%v'", err)
	}
}
```

- [ ] **Step 2: Rodar e verificar falha**

Run: `go test ./internal/usecases/login/ -v`
Expected: FAIL (erro atual é "falha, campos inválidos")

- [ ] **Step 3: Implementar**

Em `internal/usecases/login/login_use_case.go`, logo após o bloco `if loginUser.Email == "" {...}` e ANTES de `hashUser := utils.EncryptPassword(...)`:

```go
	user, err := c.userRepository.FindUserByEmail(ctx, loginUser.Email)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %v", err)
	}
	if user != nil && user.Email != "" && user.Password == "" {
		return "", fmt.Errorf("usuário pendente de primeiro acesso: solicite um novo link em \"Primeiro acesso\" na tela de login")
	}
```

E ajustar o resto da função que já usava `user, err :=` — trocar para reutilizar:

```go
	hashUser := utils.EncryptPassword(loginUser.Password)

	user, err = c.userRepository.FindUserByEmailandPassword(ctx, loginUser.Email, hashUser)
```

(o restante permanece igual)

- [ ] **Step 4: Rodar e verificar passagem**

Run: `go test ./internal/usecases/login/ -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/usecases/login/
git commit -m "feat: login informa usuario pendente de primeiro acesso"
```

---

### Task 9: Use cases de administração — list, create, update, delete, clear, get (TDD)

**Files:**
- Create: `internal/usecases/user/get_user_use_case.go` (+_test), `update_user_use_case.go` (+_test), `delete_user_use_case.go` (+_test), `clear_password_use_case.go` (+_test)
- Modify: `internal/usecases/user/create_user_use_case.go` (+_test), `list_user_use_case.go` (+_test)

**Regra de permissão compartilhada** (usar em todos os use cases deste task — colocar no arquivo `internal/usecases/user/permissions.go`):

```go
package user

import (
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
)

const (
	RoleOrganization = "organization"
	RoleAdmin        = "admin"
	RolePartner      = "partner"
)

// canListUsers: organization e admin acessam a administração de usuários.
func canListUsers(requester *entities.User) error {
	if requester.Role != RoleOrganization && requester.Role != RoleAdmin {
		return fmt.Errorf("sem permissão para acessar a administração de usuários")
	}
	return nil
}

// canManageTarget: organization gerencia admin e partner; admin gerencia apenas partner.
// Ninguém gerencia a si mesmo via ações administrativas.
func canManageTarget(requester *entities.User, target *entities.User) error {
	if err := canListUsers(requester); err != nil {
		return err
	}
	if requester.ID == target.ID {
		return fmt.Errorf("não é possível executar esta ação sobre o próprio usuário")
	}
	if requester.Role == RoleAdmin && target.Role != RolePartner {
		return fmt.Errorf("admin só pode gerenciar usuários partner")
	}
	return nil
}

// canAssignRole: organization cria/edita admin e partner; admin apenas partner.
func canAssignRole(requester *entities.User, role string) error {
	if err := canListUsers(requester); err != nil {
		return err
	}
	if role != RoleAdmin && role != RolePartner {
		return fmt.Errorf("role inválida para cadastro: use admin ou partner")
	}
	if requester.Role == RoleAdmin && role != RolePartner {
		return fmt.Errorf("admin só pode cadastrar usuários partner")
	}
	return nil
}
```

- [ ] **Step 1: Testes de permissão e list**

`internal/usecases/user/list_user_use_case_test.go` (reescrever o arquivo):

```go
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
	uc := NewListUserUseCase(requesterRepo(RoleAdmin, 1))
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
}
```

Para o segundo teste passar, o `listFunc` deve devolver `[]*entities.User{{ID: 2, Email: "x@x.com", Role: "partner"}, {ID: 3, Email: "y@x.com", Role: "admin", Password: "hash"}}` — construa o mock com:

```go
	uc := NewListUserUseCase(&mockUserRepoList{
		mockUserRepo: *requesterRepo(RoleAdmin, 1),
		listFunc: func(ctx context.Context) ([]*entities.User, error) {
			return []*entities.User{
				{ID: 2, Email: "x@x.com", Role: "partner"},               // password vazio → pendente
				{ID: 3, Email: "y@x.com", Role: "admin", Password: "h"},  // ok
			}, nil
		},
	})
```

- [ ] **Step 2: Implementar list (reescrever `list_user_use_case.go`)**

```go
package user

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
)

type listUserUseCase struct {
	userRepository repositories.UserRepository
}

func NewListUserUseCase(userRepository repositories.UserRepository) contracts.ListUserUseCase {
	return &listUserUseCase{userRepository: userRepository}
}

func (l *listUserUseCase) Execute(ctx context.Context, requesterEmail string) (*output.ListUser, error) {
	requester, err := l.userRepository.FindUserByEmail(ctx, requesterEmail)
	if err != nil {
		return nil, err
	}
	if err := canListUsers(requester); err != nil {
		return nil, err
	}

	users, err := l.userRepository.ListUser(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]output.UserItem, 0, len(users))
	for _, u := range users {
		items = append(items, output.UserItem{
			ID:                 u.ID,
			Name:               u.Name,
			Email:              u.Email,
			Role:               u.Role,
			PendingFirstAccess: u.Password == "",
		})
	}
	return &output.ListUser{Users: items}, nil
}
```

Nota: `ListUser` do repository hoje usa `Select("id","email","name","role")` — remover o `.Select(...)` (manter `Preload` + `Order`) para o `Password` vir preenchido (campo é `json:"-"`, não vaza no JSON).

- [ ] **Step 3: Rodar list**

Run: `go test ./internal/usecases/user/ -run TestListUser -v`
Expected: PASS

- [ ] **Step 4: Testes do create (reescrever `create_user_use_case_test.go`)**

Assinatura única do construtor: `NewCreateUserUseCase(userRepository repositories.UserRepository, tokenRepository repoToken.PasswordResetTokenRepository, emailSender contracts.EmailSender)`. O mock `mockUserRepoCreate` embute `mockUserRepo` — usá-lo como userRepository nos testes que precisam capturar o Create.

```go
package user

import (
	"context"
	"testing"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	repoToken "github.com/jhmorais/cash-by-card/internal/repositories/token"
)

type mockUserRepoCreate struct {
	mockUserRepo
	createFunc func(ctx context.Context, entity *entities.User) error
}

func (m *mockUserRepoCreate) CreateUser(ctx context.Context, entity *entities.User) error {
	return m.createFunc(ctx, entity)
}

func noopTokenRepo() *mockTokenRepoFull {
	return &mockTokenRepoFull{
		findFunc: func(ctx context.Context, hash string) (*entities.PasswordResetToken, error) { return nil, nil },
		markFunc: func(ctx context.Context, id int64) error { return nil },
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
			return &entities.User{ID: 9, Email: email}, nil // já existe
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
				return &entities.User{}, nil
			}},
			createFunc: func(ctx context.Context, e *entities.User) error { saved = e; return nil },
		},
		noopTokenRepo(),
		sender,
	)
	out, err := uc.Execute(context.Background(), "o@x.com", &input.CreateUser{Name: "Novo", Email: "novo@x.com", Role: "admin"})
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if out.UserID == 0 {
		t.Fatal("esperado UserID preenchido")
	}
	if saved.Password != "" {
		t.Fatal("password deve nascer vazio (primeiro acesso)")
	}
	if sender.sentTo != "novo@x.com" {
		t.Fatalf("esperado email para novo@x.com, got %s", sender.sentTo)
	}
}
```

Nota: no teste de sucesso o `createFunc` precisa atribuir `e.ID = 42` antes de retornar para o `out.UserID` vir preenchido:

```go
			createFunc: func(ctx context.Context, e *entities.User) error { e.ID = 42; saved = e; return nil },
```

(e o assertion `out.UserID != 42`).

- [ ] **Step 5: Implementar create (reescrever `create_user_use_case.go`)**

```go
package user

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jhmorais/cash-by-card/config"
	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
	repoToken "github.com/jhmorais/cash-by-card/internal/repositories/token"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
	"github.com/jhmorais/cash-by-card/utils"
)

type createUserUseCase struct {
	userRepository  repositories.UserRepository
	tokenRepository repoToken.PasswordResetTokenRepository
	emailSender     contracts.EmailSender
}

func NewCreateUserUseCase(
	userRepository repositories.UserRepository,
	tokenRepository repoToken.PasswordResetTokenRepository,
	emailSender contracts.EmailSender,
) contracts.CreateUserUseCase {
	return &createUserUseCase{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		emailSender:     emailSender,
	}
}

func (c *createUserUseCase) Execute(ctx context.Context, requesterEmail string, createUser *input.CreateUser) (*output.CreateUser, error) {
	if len(createUser.Name) > 250 {
		createUser.Name = createUser.Name[:250]
	}
	if createUser.Email == "" {
		return nil, fmt.Errorf("email é obrigatório")
	}
	if createUser.Name == "" {
		return nil, fmt.Errorf("nome é obrigatório")
	}

	requester, err := c.userRepository.FindUserByEmail(ctx, requesterEmail)
	if err != nil {
		return nil, err
	}
	if err := canAssignRole(requester, createUser.Role); err != nil {
		return nil, err
	}

	existing, err := c.userRepository.FindUserByEmail(ctx, createUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	if existing != nil && existing.Email != "" {
		return nil, fmt.Errorf("já existe usuário com o mesmo email")
	}

	userEntity := &entities.User{
		Name:      createUser.Name,
		Email:     createUser.Email,
		Role:      createUser.Role,
		CreatedAt: time.Now(),
	} // Password fica vazio (NULL) => pendente de primeiro acesso

	if err := c.userRepository.CreateUser(ctx, userEntity); err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	if err := c.sendFirstAccessEmail(ctx, userEntity); err != nil {
		log.Printf("usuario criado mas falhou envio de primeiro acesso: %v", err)
	}

	return &output.CreateUser{UserID: userEntity.ID}, nil
}

// sendFirstAccessEmail gera o token e envia o link de definicao de senha.
func (c *createUserUseCase) sendFirstAccessEmail(ctx context.Context, user *entities.User) error {
	plain, hash, err := utils.GenerateResetToken()
	if err != nil {
		return err
	}
	token := &entities.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(30 * time.Minute),
		CreatedAt: time.Now(),
	}
	if err := c.tokenRepository.CreateToken(ctx, token); err != nil {
		return err
	}
	link := config.GetFrontURL() + "/primeiro-acesso?token=" + plain
	return c.emailSender.SendPasswordResetEmail(ctx, user.Email, link)
}
```

- [ ] **Step 6: Rodar create**

Run: `go test ./internal/usecases/user/ -run TestCreateUser -v`
Expected: PASS nos 4

- [ ] **Step 7: Testes + impl de update**

`internal/usecases/user/update_user_use_case_test.go`:

```go
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
	uc := NewUpdateUserUseCase(requesterRepo(RoleAdmin, 1))
	_, err := uc.Execute(context.Background(), "a@x.com", &input.UpdateUser{ID: 5, Name: "N", Role: "admin"})
	if err == nil {
		t.Fatal("admin não pode editar admin")
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
```

`internal/usecases/user/update_user_use_case.go`:

```go
package user

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
)

type updateUserUseCase struct {
	userRepository repositories.UserRepository
}

func NewUpdateUserUseCase(userRepository repositories.UserRepository) contracts.UpdateUserUseCase {
	return &updateUserUseCase{userRepository: userRepository}
}

func (u *updateUserUseCase) Execute(ctx context.Context, requesterEmail string, updateUser *input.UpdateUser) (*output.UpdateUser, error) {
	if updateUser.Name == "" {
		return nil, fmt.Errorf("nome é obrigatório")
	}
	requester, err := u.userRepository.FindUserByEmail(ctx, requesterEmail)
	if err != nil {
		return nil, err
	}
	target, err := u.userRepository.FindUserByID(ctx, updateUser.ID)
	if err != nil || target == nil || target.ID == 0 {
		return nil, fmt.Errorf("usuário não encontrado")
	}
	if err := canManageTarget(requester, target); err != nil {
		return nil, err
	}
	if updateUser.Role != "" && updateUser.Role != target.Role {
		if err := canAssignRole(requester, updateUser.Role); err != nil {
			return nil, err
		}
		target.Role = updateUser.Role
	}
	target.Name = updateUser.Name

	if err := u.userRepository.UpdateUser(ctx, target); err != nil {
		return nil, err
	}
	return &output.UpdateUser{UserID: target.ID}, nil
}
```

Run: `go test ./internal/usecases/user/ -run TestUpdateUser -v` → PASS

- [ ] **Step 8: Testes + impl de delete**

`internal/usecases/user/delete_user_use_case_test.go`:

```go
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

func TestDeleteUser_NaoPodeExcluirASiMesmo(t *testing.T) {
	uc := NewDeleteUserUseCase(requesterRepo(RoleOrganization, 1))
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
	})
	_, err := uc.Execute(context.Background(), "a@x.com", 5)
	if err == nil {
		t.Fatal("admin não pode excluir admin")
	}
}

func TestDeleteUser_OrganizationExcluiPartner(t *testing.T) {
	deleted := false
	uc := NewDeleteUserUseCase(&mockUserRepoDelete{
		mockUserRepo: *requesterRepo(RoleOrganization, 1),
		byID: func(ctx context.Context, id int) (*entities.User, error) {
			return &entities.User{ID: 6, Role: "partner", Email: "p@x.com"}, nil
		},
		deletef: func(ctx context.Context, e *entities.User) error { deleted = true; return nil },
	})
	out, err := uc.Execute(context.Background(), "o@x.com", 6)
	if err != nil || !out.Success || !deleted {
		t.Fatalf("esperado sucesso, err='%v' success=%v deleted=%v", err, out.Success, deleted)
	}
}
```

`internal/usecases/user/delete_user_use_case.go`:

```go
package user

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
)

type deleteUserUseCase struct {
	userRepository repositories.UserRepository
}

func NewDeleteUserUseCase(userRepository repositories.UserRepository) contracts.DeleteUserUseCase {
	return &deleteUserUseCase{userRepository: userRepository}
}

func (d *deleteUserUseCase) Execute(ctx context.Context, requesterEmail string, id int) (*output.DeleteUser, error) {
	requester, err := d.userRepository.FindUserByEmail(ctx, requesterEmail)
	if err != nil {
		return nil, err
	}
	target, err := d.userRepository.FindUserByID(ctx, id)
	if err != nil || target == nil || target.ID == 0 {
		return nil, fmt.Errorf("usuário não encontrado")
	}
	if err := canManageTarget(requester, target); err != nil {
		return nil, err
	}
	if err := d.userRepository.DeleteUser(ctx, target); err != nil {
		return nil, err
	}
	return &output.DeleteUser{Success: true}, nil
}
```

Run: `go test ./internal/usecases/user/ -run TestDeleteUser -v` → PASS

- [ ] **Step 9: Testes + impl de clear password**

`internal/usecases/user/clear_password_use_case_test.go`:

```go
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
		nil, nil,
	)
	_, err := uc.Execute(context.Background(), "a@x.com", 5)
	if err == nil {
		t.Fatal("admin não pode limpar senha de admin")
	}
}

func TestClearPassword_OrganizationLimpaEEnvia(t *testing.T) {
	cleared := false
	sender := &mockEmailSender{}
	tokenRepo := &mockTokenRepoFull{
		findFunc: func(ctx context.Context, hash string) (*entities.PasswordResetToken, error) { return nil, nil },
		markFunc: func(ctx context.Context, id int64) error { return nil },
	}
	uc := NewClearPasswordUseCase(
		&mockUserRepoClear{
			mockUserRepo: *requesterRepo(RoleOrganization, 1),
			byID: func(ctx context.Context, id int) (*entities.User, error) {
				return &entities.User{ID: 6, Role: "partner", Email: "p@x.com", Name: "P"}, nil
			},
			clearf: func(ctx context.Context, id int) error { cleared = true; return nil },
		},
		tokenRepo,
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
```

`internal/usecases/user/clear_password_use_case.go`:

```go
package user

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jhmorais/cash-by-card/config"
	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
	repoToken "github.com/jhmorais/cash-by-card/internal/repositories/token"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
	"github.com/jhmorais/cash-by-card/utils"
)

type clearPasswordUseCase struct {
	userRepository  repositories.UserRepository
	tokenRepository repoToken.PasswordResetTokenRepository
	emailSender     contracts.EmailSender
}

func NewClearPasswordUseCase(
	userRepository repositories.UserRepository,
	tokenRepository repoToken.PasswordResetTokenRepository,
	emailSender contracts.EmailSender,
) contracts.ClearPasswordUseCase {
	return &clearPasswordUseCase{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		emailSender:     emailSender,
	}
}

func (c *clearPasswordUseCase) Execute(ctx context.Context, requesterEmail string, id int) (*output.ClearPassword, error) {
	requester, err := c.userRepository.FindUserByEmail(ctx, requesterEmail)
	if err != nil {
		return nil, err
	}
	target, err := c.userRepository.FindUserByID(ctx, id)
	if err != nil || target == nil || target.ID == 0 {
		return nil, fmt.Errorf("usuário não encontrado")
	}
	if err := canManageTarget(requester, target); err != nil {
		return nil, err
	}

	if err := c.userRepository.ClearUserPassword(ctx, target.ID); err != nil {
		return nil, err
	}

	plain, hash, err := utils.GenerateResetToken()
	if err != nil {
		return nil, err
	}
	token := &entities.PasswordResetToken{
		UserID:    target.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(30 * time.Minute),
		CreatedAt: time.Now(),
	}
	if err := c.tokenRepository.CreateToken(ctx, token); err != nil {
		return nil, err
	}
	link := config.GetFrontURL() + "/primeiro-acesso?token=" + plain
	if err := c.emailSender.SendPasswordResetEmail(ctx, target.Email, link); err != nil {
		log.Printf("senha limpa mas falhou envio do email: %v", err)
	}

	return &output.ClearPassword{UserID: target.ID}, nil
}
```

Run: `go test ./internal/usecases/user/ -run TestClearPassword -v` → PASS

- [ ] **Step 10: get user (me)**

`internal/usecases/user/get_user_use_case_test.go`:

```go
package user

import (
	"context"
	"testing"
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
```

`internal/usecases/user/get_user_use_case.go`:

```go
package user

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/user"
)

type getUserUseCase struct {
	userRepository repositories.UserRepository
}

func NewGetUserUseCase(userRepository repositories.UserRepository) contracts.GetUserUseCase {
	return &getUserUseCase{userRepository: userRepository}
}

func (g *getUserUseCase) Execute(ctx context.Context, email string) (*output.GetUser, error) {
	user, err := g.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Email == "" {
		return nil, fmt.Errorf("usuário não encontrado")
	}
	return &output.GetUser{ID: user.ID, Email: user.Email, Name: user.Name, Role: user.Role}, nil
}
```

Run: `go test ./internal/usecases/user/ -v` → todos PASS

- [ ] **Step 11: Commit (Tasks 6+7+8+9 juntos — contratos só compilam com os use cases)**

```bash
git add internal/contracts/ internal/usecases/ internal/repositories/user/user_repository.go
git commit -m "feat: use cases de administracao de usuarios com matriz de permissoes"
```

---

### Task 10: RoleMiddleware variadic

**Files:**
- Modify: `utils/rest.go` (RoleMiddleware), `services/rest_server.go` (call site)

- [ ] **Step 1: Trocar a assinatura**

Em `utils/rest.go`, `RoleMiddleware`:

```go
func RoleMiddleware(userRepo repositories.UserRepository, roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			email := r.Context().Value(emailKey).(string)
			ctx := context.Background()

			user, err := userRepo.FindUserByEmail(ctx, email)
			if err != nil {
				WriteErrModel(w, http.StatusUnauthorized,
					NewErrorResponse("user not found"))
				return
			}

			for _, role := range roles {
				if user.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			WriteErrModel(w, http.StatusForbidden,
				NewErrorResponse("forbidden"))
		})
	}
}
```

- [ ] **Step 2: Atualizar o call site em services/rest_server.go**

```go
	adminRouter.Use(utils.RoleMiddleware(userRepo.UserRepository, "organization", "admin"))
```

- [ ] **Step 3: Build + test + commit**

Run: `go build ./... && go test ./...`
Expected: sem erro

```bash
git add utils/rest.go services/rest_server.go
git commit -m "feat: RoleMiddleware aceita multiplos roles (organization, admin)"
```

---

### Task 11: Handlers e rotas

**Files:**
- Create: `services/user_account.go`
- Modify: `services/user.go`, `services/rest_server.go`, `internal/infra/di/dependency_builder.go`

- [ ] **Step 1: DI — novos use cases**

Em `internal/infra/di/dependency_builder.go`:

Usecases struct — adicionar campos:

```go
	ForgotPasswordUseCase  contracts.ForgotPasswordUseCase
	ResetPasswordUseCase   contracts.ResetPasswordUseCase
	ChangePasswordUseCase  contracts.ChangePasswordUseCase
	GetUserUseCase         contracts.GetUserUseCase
	UpdateUserUseCase      contracts.UpdateUserUseCase
	DeleteUserUseCase      contracts.DeleteUserUseCase
	ClearPasswordUseCase   contracts.ClearPasswordUseCase
```

Repositories struct — adicionar:

```go
	PasswordResetTokenRepository repoToken.PasswordResetTokenRepository
```

(import `repoToken "github.com/jhmorais/cash-by-card/internal/repositories/token"` e `"github.com/jhmorais/cash-by-card/internal/infra/email"`)

No corpo de `NewBuild()`, junto às demais construções:

```go
	tokenRepository := repoToken.NewPasswordResetTokenRepository(builder.DB)
	emailSender := email.NewSenderFromEnv()
```

E as atribuições de use case (seguindo o padrão existente):

```go
	ForgotPasswordUseCase:  user.NewForgotPasswordUseCase(builder.Repositories.UserRepository, tokenRepository, emailSender),
	ResetPasswordUseCase:   user.NewResetPasswordUseCase(tokenRepository, builder.Repositories.UserRepository),
	ChangePasswordUseCase:  user.NewChangePasswordUseCase(builder.Repositories.UserRepository),
	GetUserUseCase:         user.NewGetUserUseCase(builder.Repositories.UserRepository),
	UpdateUserUseCase:      user.NewUpdateUserUseCase(builder.Repositories.UserRepository),
	DeleteUserUseCase:      user.NewDeleteUserUseCase(builder.Repositories.UserRepository),
	ClearPasswordUseCase:   user.NewClearPasswordUseCase(builder.Repositories.UserRepository, tokenRepository, emailSender),
	CreateUserUseCase:      user.NewCreateUserUseCase(builder.Repositories.UserRepository, tokenRepository, emailSender),
	ListUserUseCase:        user.NewListUserUseCase(builder.Repositories.UserRepository),
```

(Ajuste nomes de campos conforme o struct real — os use cases de user já são construídos em `NewBuild`; substituir as linhas antigas de CreateUserUseCase/ListUserUseCase pelas novas com dependências.)

- [ ] **Step 2: Handlers de conta**

`services/user_account.go`:

```go
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	"github.com/jhmorais/cash-by-card/utils"
)

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading request body"))
		return
	}
	forgot := input.ForgotPassword{}
	if err := json.Unmarshal(body, &forgot); err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("failed to parse request body"))
		return
	}
	if err := h.ForgotPasswordUseCase.Execute(ctx, &forgot); err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse(fmt.Sprintf("failed to process forgot password, error: '%s'", err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"message":"Se o email estiver cadastrado, você receberá as instruções por email","type":"SUCCESS"}`)
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading request body"))
		return
	}
	reset := input.ResetPassword{}
	if err := json.Unmarshal(body, &reset); err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("failed to parse request body"))
		return
	}
	if err := h.ResetPasswordUseCase.Execute(ctx, &reset); err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest,
			utils.NewErrorResponse(fmt.Sprintf("failed to reset password, error: '%s'", err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"message":"senha definida com sucesso","type":"SUCCESS"}`)
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading request body"))
		return
	}
	change := input.ChangePassword{}
	if err := json.Unmarshal(body, &change); err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("failed to parse request body"))
		return
	}
	email := utils.EmailFromContext(r.Context())
	if err := h.ChangePasswordUseCase.Execute(ctx, email, &change); err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest,
			utils.NewErrorResponse(fmt.Sprintf("failed to change password, error: '%s'", err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"message":"senha alterada com sucesso","type":"SUCCESS"}`)
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	email := utils.EmailFromContext(r.Context())
	user, err := h.GetUserUseCase.Execute(context.Background(), email)
	if err != nil {
		utils.WriteErrModel(w, http.StatusNotFound,
			utils.NewErrorResponse(fmt.Sprintf("failed to get user, error: '%s'", err.Error())))
		return
	}
	jsonResponse, err := json.Marshal(user)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError, utils.NewErrorResponse("Failed to marshal user response"))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}
```

- [ ] **Step 3: Reescrever handlers de administração em services/user.go**

Substituir `ListUsers` e `CreateUser` existentes e adicionar `UpdateUser`, `DeleteUser`, `ClearUserPassword` (mesmo padrão dos demais — Lê body/param, chama use case com `utils.EmailFromContext(r.Context())`):

```go
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	requesterEmail := utils.EmailFromContext(r.Context())
	response, err := h.ListUserUseCase.Execute(context.Background(), requesterEmail)
	if err != nil {
		utils.WriteErrModel(w, http.StatusForbidden,
			utils.NewErrorResponse(fmt.Sprintf("failed to list users, error: '%s'", err.Error())))
		return
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError, utils.NewErrorResponse("Failed to marshal user response"))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading request body"))
		return
	}
	createUser := input.CreateUser{}
	if err := json.Unmarshal(body, &createUser); err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("failed to parse request body"))
		return
	}
	requesterEmail := utils.EmailFromContext(r.Context())
	response, err := h.CreateUserUseCase.Execute(context.Background(), requesterEmail, &createUser)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest,
			utils.NewErrorResponse(fmt.Sprintf("failed to create user, error: '%s'", err.Error())))
		return
	}
	jsonResponse, _ := json.Marshal(response)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(jsonResponse))
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idParam, err := utils.RetrieveParam(r, "id")
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading id"))
		return
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error casting id"))
		return
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading request body"))
		return
	}
	updateUser := input.UpdateUser{ID: id}
	if err := json.Unmarshal(body, &updateUser); err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("failed to parse request body"))
		return
	}
	updateUser.ID = id
	requesterEmail := utils.EmailFromContext(r.Context())
	response, err := h.UpdateUserUseCase.Execute(context.Background(), requesterEmail, &updateUser)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest,
			utils.NewErrorResponse(fmt.Sprintf("failed to update user, error: '%s'", err.Error())))
		return
	}
	jsonResponse, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idParam, err := utils.RetrieveParam(r, "id")
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading id"))
		return
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error casting id"))
		return
	}
	requesterEmail := utils.EmailFromContext(r.Context())
	response, err := h.DeleteUserUseCase.Execute(context.Background(), requesterEmail, id)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest,
			utils.NewErrorResponse(fmt.Sprintf("failed to delete user, error: '%s'", err.Error())))
		return
	}
	jsonResponse, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func (h *Handler) ClearUserPassword(w http.ResponseWriter, r *http.Request) {
	idParam, err := utils.RetrieveParam(r, "id")
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading id"))
		return
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error casting id"))
		return
	}
	requesterEmail := utils.EmailFromContext(r.Context())
	response, err := h.ClearPasswordUseCase.Execute(context.Background(), requesterEmail, id)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest,
			utils.NewErrorResponse(fmt.Sprintf("failed to clear password, error: '%s'", err.Error())))
		return
	}
	jsonResponse, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}
```

Nota: `utils.RetrieveParam(r, "id")` devolve `(string, error)` — por isso os dois passos com `strconv.Atoi`, como nos handlers existentes de loans. Imports do arquivo: `context`, `encoding/json`, `fmt`, `io`, `net/http`, `strconv`, `input "github.com/jhmorais/cash-by-card/internal/ports/input/user"`, `"github.com/jhmorais/cash-by-card/utils"`.

- [ ] **Step 4: Rotas em services/rest_server.go**

Handler struct — adicionar campos:

```go
	ForgotPasswordUseCase  contracts.ForgotPasswordUseCase
	ResetPasswordUseCase   contracts.ResetPasswordUseCase
	ChangePasswordUseCase  contracts.ChangePasswordUseCase
	GetUserUseCase         contracts.GetUserUseCase
	UpdateUserUseCase      contracts.UpdateUserUseCase
	DeleteUserUseCase      contracts.DeleteUserUseCase
	ClearPasswordUseCase   contracts.ClearPasswordUseCase
```

Construção do handler — preencher com `useCases.<Campo>`.

Rotas:

```go
	authRouter.HandleFunc("/login", handler.LoginUser).Methods(http.MethodPost)
	authRouter.HandleFunc("/forgot-password", handler.ForgotPassword).Methods(http.MethodPost)
	authRouter.HandleFunc("/reset-password", handler.ResetPassword).Methods(http.MethodPost)

	accountRouter := router.PathPrefix("/account").Subrouter()
	accountRouter.Use(utils.ValidateJwtTokenMiddleware)
	accountRouter.HandleFunc("/me", handler.GetMe).Methods(http.MethodGet)
	accountRouter.HandleFunc("/change-password", handler.ChangePassword).Methods(http.MethodPost)

	adminRouter.HandleFunc("/users", handler.ListUsers).Methods(http.MethodGet)
	adminRouter.HandleFunc("/users", handler.CreateUser).Methods(http.MethodPost)
	adminRouter.HandleFunc("/users/{id}", handler.UpdateUser).Methods(http.MethodPut)
	adminRouter.HandleFunc("/users/{id}", handler.DeleteUser).Methods(http.MethodDelete)
	adminRouter.HandleFunc("/users/{id}/clear-password", handler.ClearUserPassword).Methods(http.MethodPost)
```

- [ ] **Step 5: Build + test + commit**

Run: `go build ./... && go test ./...`
Expected: sem erro

```bash
git add services/ internal/infra/di/
git commit -m "feat: rotas de conta e administracao de usuarios"
```

---

### Task 12: Atualizar seed do banco local

**Files:** nenhum do repo (dados locais)

- [ ] **Step 1: Promover o usuário de teste a organization**

```bash
PWD_HASH=$(echo -n "admin123" | md5sum | cut -d' ' -f1)
sg docker -c "docker exec mysqlcontainer mysql -uroot -ppassword database -e \"UPDATE user SET role='organization' WHERE email='gontijogabi93@gmail.com';\""
sg docker -c "docker exec mysqlcontainer mysql -uroot -ppassword database -e 'SELECT id,email,role FROM user;'"
```

Expected: `gontijogabi93@gmail.com | organization`

- [ ] **Step 2: Sem commit (dados locais)**

---

### Task 13: Verificação ao vivo do backend (E2E curl)

**Files:** nenhum (validação)

- [ ] **Step 1: Reiniciar o server (go run não recarrega!)**

```bash
lsof -ti:3000 -sTCP:LISTEN | xargs -r kill; sleep 1
cd /home/gabigontijo/Documents/cash-by-card && go run cmd/restserver/main.go &> /tmp/cashbycard.log &
```

- [ ] **Step 2: Cenário completo**

```bash
# login organization
JWT=$(curl -s -X POST localhost:3000/auth/login -H 'Content-Type: application/json' \
  -d '{"email":"gontijogabi93@gmail.com","password":"admin123"}' | sed 's/^"//; s/"$//')

# me
curl -s localhost:3000/account/me -H "Authorization: Bearer $JWT"
# expected: {"id":..,"email":"gontijogabi93@gmail.com","name":"Gabi","role":"organization"}

# cria partner -> token no console (/tmp/cashbycard.log, linha [DEV-EMAIL])
curl -s -X POST localhost:3000/admin/users -H "Authorization: Bearer $JWT" -H 'Content-Type: application/json' \
  -d '{"name":"Parceiro Um","email":"parceiro1@teste.com","role":"partner"}'
grep DEV-EMAIL /tmp/cashbycard.log | tail -1
TOKEN=<extrair o token do link do log>

# login ANTES do reset -> erro de primeiro acesso
curl -s -X POST localhost:3000/auth/login -H 'Content-Type: application/json' \
  -d '{"email":"parceiro1@teste.com","password":"x"}'
# expected: mensagem com "primeiro acesso"

# reset com o token
curl -s -X POST localhost:3000/auth/reset-password -H 'Content-Type: application/json' \
  -d "{\"token\":\"$TOKEN\",\"newPassword\":\"senha123\"}"
# expected: sucesso

# reuso do MESMO token -> deve falhar (single-use)
curl -s -X POST localhost:3000/auth/reset-password -H 'Content-Type: application/json' \
  -d "{\"token\":\"$TOKEN\",\"newPassword\":\"outra123\"}"
# expected: token inválido ou expirado

# login com a nova senha
curl -s -X POST localhost:3000/auth/login -H 'Content-Type: application/json' \
  -d '{"email":"parceiro1@teste.com","password":"senha123"}'
# expected: JWT

# partner tenta listar users -> 403/erro
JWT_P=$(curl -s -X POST localhost:3000/auth/login -H 'Content-Type: application/json' \
  -d '{"email":"parceiro1@teste.com","password":"senha123"}' | sed 's/^"//; s/"$//')
curl -s localhost:3000/admin/users -H "Authorization: Bearer $JWT_P"
# expected: sem permissão

# change-password
curl -s -X POST localhost:3000/account/change-password -H "Authorization: Bearer $JWT_P" \
  -H 'Content-Type: application/json' -d '{"currentPassword":"senha123","newPassword":"nova456"}'

# clear-password pela organization + reuso do fluxo
curl -s -X POST localhost:3000/admin/users/<id-do-partner>/clear-password -H "Authorization: Bearer $JWT"
# expected: userId; novo DEV-EMAIL no log; password NULL no banco

# forgot-password de email inexistente -> mesma resposta genérica
curl -s -X POST localhost:3000/auth/forgot-password -H 'Content-Type: application/json' \
  -d '{"email":"naoexiste@x.com"}'
```

Todos os cenários devem se comportar como nos comentários `expected:`.

---

# FRONTEND

(`/home/gabigontijo/Documents/cash-by-card-front`; `npm test` roda vitest)

### Task 14: API de user + testes

**Files:**
- Create: `src/apis/user/index.js`
- Test: `src/apis/user/index.test.js`

- [ ] **Step 1: Implementar (rodar o teste depois — aqui o teste é de contrato da API, escrever junto)**

`src/apis/user/index.js`:

```js
import { apiFetch } from '..';

export const listUsers = async () => {
  const res = await apiFetch('admin/users', { method: 'get' });
  return res.json();
};

export const createUser = async (user) => {
  const res = await apiFetch('admin/users', {
    method: 'post',
    body: JSON.stringify(user),
  });
  return res.json();
};

export const updateUser = async (id, user) => {
  const res = await apiFetch(`admin/users/${id}`, {
    method: 'put',
    body: JSON.stringify(user),
  });
  return res.json();
};

export const deleteUser = async (id) => {
  const res = await apiFetch(`admin/users/${id}`, { method: 'delete' });
  return res.json();
};

export const clearUserPassword = async (id) => {
  const res = await apiFetch(`admin/users/${id}/clear-password`, { method: 'post' });
  return res.json();
};

export const getMe = async () => {
  const res = await apiFetch('account/me', { method: 'get' });
  return res.json();
};

export const forgotPassword = async (email) => {
  const res = await apiFetch('auth/forgot-password', {
    method: 'post',
    body: JSON.stringify({ email }),
  });
  return res.json();
};

export const resetPassword = async (token, newPassword) => {
  const res = await apiFetch('auth/reset-password', {
    method: 'post',
    body: JSON.stringify({ token, newPassword }),
  });
  return res.json();
};

export const changePassword = async (currentPassword, newPassword) => {
  const res = await apiFetch('account/change-password', {
    method: 'post',
    body: JSON.stringify({ currentPassword, newPassword }),
  });
  return res.json();
};
```

- [ ] **Step 2: Teste**

`src/apis/user/index.test.js` (padrão de `src/apis/loan/index.test.js`):

```js
import { it, vi, expect, describe, afterEach, beforeEach } from 'vitest';

import { listUsers, createUser, clearUserPassword, forgotPassword, resetPassword } from '.';

describe('user apis', () => {
  let fetchMock;

  beforeEach(() => {
    vi.stubEnv('VITE_API_BASE_URL', 'http://localhost:3000');
    window.localStorage.setItem('token', 'fake-token');
    fetchMock = vi.fn(() =>
      Promise.resolve({ ok: true, json: () => Promise.resolve({ users: [] }) })
    );
    vi.stubGlobal('fetch', fetchMock);
  });

  afterEach(() => {
    vi.unstubAllGlobals();
    vi.unstubAllEnvs();
    window.localStorage.clear();
  });

  it('listUsers chama GET admin/users com bearer', async () => {
    await listUsers();
    const [url, options] = fetchMock.mock.calls[0];
    expect(url).toBe('http://localhost:3000/admin/users');
    expect(options.method).toBe('get');
    expect(options.headers.Authorization).toBe('Bearer fake-token');
  });

  it('createUser envia POST com body', async () => {
    await createUser({ name: 'N', email: 'n@x.com', role: 'partner' });
    const [url, options] = fetchMock.mock.calls[0];
    expect(options.method).toBe('post');
    expect(JSON.parse(options.body)).toEqual({ name: 'N', email: 'n@x.com', role: 'partner' });
  });

  it('clearUserPassword usa a rota de clear-password', async () => {
    await clearUserPassword(7);
    const [url] = fetchMock.mock.calls[0];
    expect(url).toBe('http://localhost:3000/admin/users/7/clear-password');
  });

  it('forgotPassword não envia token', async () => {
    window.localStorage.removeItem('token');
    await forgotPassword('a@b.com');
    const [url, options] = fetchMock.mock.calls[0];
    expect(url).toBe('http://localhost:3000/auth/forgot-password');
    expect(options.headers.Authorization).toBeUndefined();
  });

  it('resetPassword envia token e nova senha', async () => {
    await resetPassword('tok', 'nova123');
    const [, options] = fetchMock.mock.calls[0];
    expect(JSON.parse(options.body)).toEqual({ token: 'tok', newPassword: 'nova123' });
  });
});
```

- [ ] **Step 3: Rodar**

Run: `npm test`
Expected: PASS (todos os suites)

- [ ] **Step 4: Commit**

```bash
git add src/apis/user/
git commit -m "feat: apis de conta e administracao de usuarios"
```

---

### Task 15: AuthProvider com role

**Files:**
- Modify: `src/hooks/authProvider.jsx`

- [ ] **Step 1: Alterações**

No `loginAction`, após gravar o token, buscar o role:

```js
import { login } from 'src/apis/login';
import { getMe } from 'src/apis/user';
```

```js
const [role, setRole] = useState(localStorage.getItem('role') || '');

const loginAction = async (email, password) => {
  const response = await login(email, password);
  if (response) {
    setUser(email);
    setToken(response);
    localStorage.setItem('token', response);
    localStorage.setItem('user', email);
    try {
      const me = await getMe();
      setRole(me.role);
      localStorage.setItem('role', me.role);
    } catch (e) {
      setRole('');
      localStorage.removeItem('role');
    }
    navigate('/');
  } else {
    throw new Error(response.message || 'Erro ao fazer o login');
  }
};

const logOut = () => {
  setUser(null);
  setToken('');
  setRole('');
  localStorage.removeItem('token');
  localStorage.removeItem('user');
  localStorage.removeItem('role');
  navigate('/login');
};
```

No `useEffect` de restauração e no `authValue`:

```js
    const storedRole = localStorage.getItem('role');
    if (storedRole) setRole(storedRole);

  const authValue = useMemo(() => ({ token, user, role, loginAction, logOut }), [token, user, role]);
```

- [ ] **Step 2: Commit**

```bash
git add src/hooks/authProvider.jsx
git commit -m "feat: AuthProvider expoe role do usuario logado"
```

---

### Task 16: Página de primeiro acesso

**Files:**
- Create: `src/pages/primeiro-acesso.jsx`, `src/sections/first-access/view/first-access-view.jsx`, `src/sections/first-access/view/index.js`
- Test: `src/sections/first-access/view/first-access-view.test.jsx`

- [ ] **Step 1: View**

`src/sections/first-access/view/first-access-view.jsx`:

```jsx
import { useState } from 'react';
import { Link as RouterLink, useNavigate, useSearchParams } from 'react-router-dom';
// components MUI no padrão de import do projeto (src/sections/login/login-view.jsx):
import Stack from '@mui/material/Stack';
import TextField from '@mui/material/TextField';
import Alert from '@mui/material/Alert';
import Link from '@mui/material/Link';
import LoadingButton from '@mui/lab/LoadingButton';

import { resetPassword } from 'src/apis/user';

export default function FirstAccessView() {
  const [params] = useSearchParams();
  const navigate = useNavigate();
  const [password, setPassword] = useState('');
  const [confirm, setConfirm] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [isLoading, setLoading] = useState(false);

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError('');
    if (password.length < 6) {
      setError('A senha deve ter pelo menos 6 caracteres');
      return;
    }
    if (password !== confirm) {
      setError('As senhas não conferem');
      return;
    }
    setLoading(true);
    try {
      await resetPassword(params.get('token') || '', password);
      setSuccess('Senha definida com sucesso! Você já pode fazer login.');
      setTimeout(() => navigate('/login'), 2000);
    } catch (err) {
      setError(err.message || 'Não foi possível definir a senha');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Stack spacing={3} component="form" onSubmit={handleSubmit}>
      {error && <Alert severity="error">{error}</Alert>}
      {success && <Alert severity="success">{success}</Alert>}
      {!params.get('token') && <Alert severity="warning">Token ausente. Solicite um novo pelo link "Primeiro acesso" na tela de login.</Alert>}
      <TextField
        name="password"
        label="Nova senha"
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <TextField
        name="confirm"
        label="Confirmar nova senha"
        type="password"
        value={confirm}
        onChange={(e) => setConfirm(e.target.value)}
      />
      <LoadingButton loading={isLoading} fullWidth size="large" type="submit" variant="contained">
        Definir senha
      </LoadingButton>
      <Link component={RouterLink} to="/login" variant="subtitle2" underline="hover">
        Voltar para o login
      </Link>
    </Stack>
  );
}
```

`src/sections/first-access/view/index.js`:

```js
export { default } from './first-access-view';
```

- [ ] **Step 2: Página**

`src/pages/primeiro-acesso.jsx` (mesmo wrapper do login — o layout centrado do LoginView vem da rota; usar Card igual ao login page):

```jsx
import { Helmet } from 'react-helmet-async';

import FirstAccessView from 'src/sections/first-access/view';

export default function FirstAccessPage() {
  return (
    <>
      <Helmet>
        <title> Primeiro acesso | Cash By Card </title>
      </Helmet>
      <FirstAccessView />
    </>
  );
}
```

(Se o LoginView usa um Container/Card externo, replicar a mesma moldura — conferir `src/sections/login/login-view.jsx` e copiar o wrapper.)

- [ ] **Step 3: Teste**

`src/sections/first-access/view/first-access-view.test.jsx`:

```jsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { it, vi, expect, describe, afterEach, beforeEach } from 'vitest';

import FirstAccessView from './first-access-view';

describe('FirstAccessView', () => {
  beforeEach(() => {
    vi.stubEnv('VITE_API_BASE_URL', 'http://localhost:3000');
    vi.stubGlobal('fetch', vi.fn(() =>
      Promise.resolve({ ok: true, json: () => Promise.resolve({ message: 'ok' }) })
    ));
  });
  afterEach(() => {
    vi.unstubAllGlobals();
    vi.unstubAllEnvs();
  });

  it('mostra aviso quando nao ha token na url', () => {
    render(<MemoryRouter><FirstAccessView /></MemoryRouter>);
    expect(screen.getByText(/Token ausente/i)).toBeTruthy();
  });

  it('recusa senhas que nao conferem', async () => {
    render(
      <MemoryRouter initialEntries={['/primeiro-acesso?token=abc']}>
        <FirstAccessView />
      </MemoryRouter>
    );
    fireEvent.change(screen.getByLabelText('Nova senha'), { target: { value: 'senha123' } });
    fireEvent.change(screen.getByLabelText('Confirmar nova senha'), { target: { value: 'diferente' } });
    fireEvent.click(screen.getByRole('button', { name: /Definir senha/i }));
    expect(await screen.findByText('As senhas não conferem')).toBeTruthy();
  });
});
```

- [ ] **Step 4: Rodar + commit**

Run: `npm test` → PASS

```bash
git add src/pages/primeiro-acesso.jsx src/sections/first-access/
git commit -m "feat: pagina de primeiro acesso com token"
```

---

### Task 17: Página de alterar senha

**Files:**
- Create: `src/pages/alterar-senha.jsx`, `src/sections/change-password/view/change-password-view.jsx`, `src/sections/change-password/view/index.js`

- [ ] **Step 1: View (mesma estrutura da Task 16, com 3 campos e submit)**

`src/sections/change-password/view/change-password-view.jsx`:

```jsx
import { useState } from 'react';
import Stack from '@mui/material/Stack';
import TextField from '@mui/material/TextField';
import Alert from '@mui/material/Alert';
import LoadingButton from '@mui/lab/LoadingButton';

import { changePassword } from 'src/apis/user';

export default function ChangePasswordView() {
  const [current, setCurrent] = useState('');
  const [password, setPassword] = useState('');
  const [confirm, setConfirm] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [isLoading, setLoading] = useState(false);

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError('');
    setSuccess('');
    if (password.length < 6) return setError('A nova senha deve ter pelo menos 6 caracteres');
    if (password !== confirm) return setError('As senhas não conferem');
    setLoading(true);
    try {
      await changePassword(current, password);
      setSuccess('Senha alterada com sucesso!');
      setCurrent('');
      setPassword('');
      setConfirm('');
    } catch (err) {
      setError(err.message || 'Não foi possível alterar a senha');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Stack spacing={3} component="form" onSubmit={handleSubmit}>
      {error && <Alert severity="error">{error}</Alert>}
      {success && <Alert severity="success">{success}</Alert>}
      <TextField name="current" label="Senha atual" type="password" value={current} onChange={(e) => setCurrent(e.target.value)} />
      <TextField name="password" label="Nova senha" type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
      <TextField name="confirm" label="Confirmar nova senha" type="password" value={confirm} onChange={(e) => setConfirm(e.target.value)} />
      <LoadingButton loading={isLoading} fullWidth size="large" type="submit" variant="contained">
        Alterar senha
      </LoadingButton>
    </Stack>
  );
}
```

(+ `index.js` re-export, + página `src/pages/alterar-senha.jsx` com o mesmo wrapper da Task 16, título "Alterar senha | Cash By Card")

- [ ] **Step 2: Commit**

```bash
git add src/pages/alterar-senha.jsx src/sections/change-password/
git commit -m "feat: pagina de alteracao de senha"
```

---

### Task 18: Página de administração de usuários

**Files:**
- Create: `src/pages/usuario.jsx`, `src/sections/user/view/usuario-view.jsx`, `src/sections/user/view/index.js`

- [ ] **Step 1: View completa**

`src/sections/user/view/usuario-view.jsx`:

```jsx
import { useEffect, useState } from 'react';
import { useAuth } from 'src/hooks/authProvider';
import Button from '@mui/material/Button';
import Stack from '@mui/material/Stack';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import MenuItem from '@mui/material/MenuItem';
import Alert from '@mui/material/Alert';

import { listUsers, createUser, updateUser, deleteUser, clearUserPassword } from 'src/apis/user';

const ROLE_LABEL = { organization: 'Organization', admin: 'Admin', partner: 'Parceiro' };

export default function UsuarioView() {
  const { role: myRole, user: myEmail } = useAuth();
  const [users, setUsers] = useState([]);
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editing, setEditing] = useState(null); // user em edição ou null para novo
  const [form, setForm] = useState({ name: '', email: '', role: 'partner' });

  const roleOptions = myRole === 'organization' ? ['admin', 'partner'] : ['partner'];

  const load = async () => {
    try {
      setError('');
      const data = await listUsers();
      setUsers(data.users || []);
    } catch (err) {
      setError(err.message);
    }
  };

  useEffect(() => {
    load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const openCreate = () => {
    setEditing(null);
    setForm({ name: '', email: '', role: roleOptions[0] });
    setDialogOpen(true);
  };

  const openEdit = (u) => {
    setEditing(u);
    setForm({ name: u.name, email: u.email, role: u.role });
    setDialogOpen(true);
  };

  const submit = async () => {
    try {
      setError('');
      if (editing) {
        await updateUser(editing.id, { name: form.name, role: form.role });
      } else {
        await createUser(form);
      }
      setDialogOpen(false);
      await load();
    } catch (err) {
      setError(err.message);
    }
  };

  const handleDelete = async (u) => {
    if (!window.confirm(`Excluir o usuário ${u.email}?`)) return;
    try {
      setError('');
      await deleteUser(u.id);
      await load();
    } catch (err) {
      setError(err.message);
    }
  };

  const handleClear = async (u) => {
    if (!window.confirm(`Limpar a senha de ${u.email}? Ele receberá um email para definir uma nova.`)) return;
    try {
      setError('');
      await clearUserPassword(u.id);
      setMessage(`Senha de ${u.email} limpa — email de redefinição enviado.`);
      await load();
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <Stack spacing={3}>
      {error && <Alert severity="error">{error}</Alert>}
      {message && <Alert severity="success">{message}</Alert>}
      <Stack direction="row" justifyContent="flex-end">
        <Button variant="contained" onClick={openCreate}>
          Adicionar usuário
        </Button>
      </Stack>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>Nome</TableCell>
            <TableCell>Email</TableCell>
            <TableCell>Role</TableCell>
            <TableCell>Status</TableCell>
            <TableCell>Ações</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {users.map((u) => (
            <TableRow key={u.id}>
              <TableCell>{u.name}</TableCell>
              <TableCell>{u.email}</TableCell>
              <TableCell>{ROLE_LABEL[u.role] || u.role}</TableCell>
              <TableCell>{u.pendingFirstAccess ? 'Pendente de primeiro acesso' : 'Senha definida'}</TableCell>
              <TableCell>
                <Stack direction="row" spacing={1}>
                  <Button size="small" onClick={() => openEdit(u)} disabled={u.email === myEmail}>
                    Editar
                  </Button>
                  <Button size="small" color="warning" onClick={() => handleClear(u)} disabled={u.email === myEmail}>
                    Limpar senha
                  </Button>
                  <Button size="small" color="error" onClick={() => handleDelete(u)} disabled={u.email === myEmail}>
                    Excluir
                  </Button>
                </Stack>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>

      <Dialog open={dialogOpen} onClose={() => setDialogOpen(false)}>
        <DialogTitle>{editing ? 'Editar usuário' : 'Adicionar usuário'}</DialogTitle>
        <DialogContent>
          <Stack spacing={2} sx={{ mt: 1, minWidth: 360 }}>
            <TextField label="Nome" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
            <TextField
              label="Email"
              type="email"
              value={form.email}
              disabled={!!editing}
              onChange={(e) => setForm({ ...form, email: e.target.value })}
            />
            <TextField
              select
              label="Role"
              value={form.role}
              onChange={(e) => setForm({ ...form, role: e.target.value })}
            >
              {roleOptions.map((r) => (
                <MenuItem key={r} value={r}>
                  {ROLE_LABEL[r]}
                </MenuItem>
              ))}
            </TextField>
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDialogOpen(false)}>Cancelar</Button>
          <Button variant="contained" onClick={submit}>
            Salvar
          </Button>
        </DialogActions>
      </Dialog>
    </Stack>
  );
}
```

(+ `index.js` re-export, + `src/pages/usuario.jsx` wrapper com título "Usuários | Cash By Card")

- [ ] **Step 2: Commit**

```bash
git add src/pages/usuario.jsx src/sections/user/
git commit -m "feat: pagina de administracao de usuarios"
```

---

### Task 19: Rotas, navegação e links de login

**Files:**
- Modify: `src/routes/sections.jsx`, `src/layouts/dashboard/config-navigation.jsx`, `src/layouts/dashboard/nav.jsx`, `src/sections/login/login-view.jsx`

- [ ] **Step 1: Rotas**

Em `src/routes/sections.jsx`:

```js
export const FirstAccessPage = lazy(() => import('src/pages/primeiro-acesso'));
export const ChangePasswordPage = lazy(() => import('src/pages/alterar-senha'));
export const UsuarioPage = lazy(() => import('src/pages/usuario'));
```

No array de rotas — dentro dos children do PrivateRoute:

```js
        { path: 'alterar-senha', element: <ChangePasswordPage /> },
        { path: 'usuarios', element: <UsuarioPage /> },
```

Fora (pública):

```js
    {
      path: 'primeiro-acesso',
      element: <FirstAccessPage />,
    },
```

- [ ] **Step 2: Nav config**

`src/layouts/dashboard/config-navigation.jsx` — adicionar ao array:

```js
  {
    title: 'usuários',
    path: '/usuarios',
    icon: icon('ic_lock'),
  },
  {
    title: 'alterar senha',
    path: '/alterar-senha',
    icon: icon('ic_lock'),
  },
```

- [ ] **Step 3: nav.jsx — filtrar por role**

Em `src/layouts/dashboard/nav.jsx` (usa `useAuth` na linha ~44 e mapeia `navConfig` na ~116). Antes do map:

```js
  const items = navConfig.filter((item) => {
    if (item.path === '/usuarios') return ['organization', 'admin'].includes(auth.role);
    return true;
  });
```

Trocar `{navConfig.map((item) => (` por `{items.map((item) => (`.

- [ ] **Step 4: login-view — links "Primeiro acesso" e "Esqueci minha senha"**

Em `src/sections/login/login-view.jsx`:

Estados novos (junto aos existentes):

```js
  const [resetMode, setResetMode] = useState(false);
  const [resetEmail, setResetEmail] = useState('');
  const [resetSent, setResetSent] = useState('');
```

Import:

```js
import { forgotPassword } from 'src/apis/user';
```

Substituir o bloco do Link (linhas ~92-96, `<Stack direction="row" ...><Link variant="subtitle2" underline="hover">...`) por:

```jsx
      <Stack direction="row" alignItems="center" justifyContent="flex-end" spacing={2} sx={{ my: 3 }}>
        <Link variant="subtitle2" underline="hover" sx={{ cursor: 'pointer' }} onClick={() => { setResetMode(true); setResetSent(''); }}>
          Primeiro acesso
        </Link>
        <Link variant="subtitle2" underline="hover" sx={{ cursor: 'pointer' }} onClick={() => { setResetMode(true); setResetSent(''); }}>
          Esqueci minha senha
        </Link>
      </Stack>

      {resetMode && (
        <Stack spacing={2} sx={{ mb: 3 }}>
          {resetSent && <Alert severity="success">{resetSent}</Alert>}
          <TextField
            name="resetEmail"
            label="Email cadastrado"
            type="email"
            value={resetEmail}
            onChange={(e) => setResetEmail(e.target.value)}
          />
          <LoadingButton
            variant="contained"
            onClick={async () => {
              try {
                await forgotPassword(resetEmail);
                setResetSent('Se o email estiver cadastrado, você receberá as instruções por email.');
              } catch (err) {
                setAlertError(err.message);
              }
            }}
          >
            Enviar link de definição de senha
          </LoadingButton>
        </Stack>
      )}
```

(Import `Alert from '@mui/material/Alert'` se ainda não existir.)

- [ ] **Step 5: Rodar testes + commit**

Run: `npm test` → PASS

```bash
git add src/routes/ src/layouts/ src/sections/login/
git commit -m "feat: rotas, nav por role e links de primeiro acesso no login"
```

---

### Task 20: Verificação E2E com frontend rodando

- [ ] **Step 1: Subir backend e frontend**

```bash
# backend
cd /home/gabigontijo/Documents/cash-by-card && go run cmd/restserver/main.go &> /tmp/cashbycard.log &
# frontend
cd /home/gabigontijo/Documents/cash-by-card-front && npm run dev &> /tmp/cbc-front.log &
```

- [ ] **Step 2: Cenário manual (browser em http://localhost:5173)**

1. Login com `gontijogabi93@gmail.com` / `admin123` → entra; nav mostra "usuários" e "alterar senha"
2. `/usuarios` → adicionar partner "parceiro2@teste.com" → aparece "Pendente de primeiro acesso"; link DEV-EMAIL no log do backend
3. Copiar token do log → `http://localhost:5173/primeiro-acesso?token=...` → definir senha → login com ela
4. Login como esse partner → nav SEM "usuários" → `/alterar-senha` funciona
5. Partner acessa `/usuarios` direto na URL → API retorna "sem permissão" (erro exibido)
6. Como organization: limpar senha do partner → email novo no log → partner antigo não loga mais com a senha velha ("pendente de primeiro acesso")
7. Tela de login → "Esqueci minha senha" com email inexistente → mensagem genérica igual

- [ ] **Step 3: Suite completa + status final**

```bash
cd /home/gabigontijo/Documents/cash-by-card && go test ./... && cd ../cash-by-card-front && npm test
git status
```

Expected: tudo PASS, working trees limpos (todos os commits feitos por task).
