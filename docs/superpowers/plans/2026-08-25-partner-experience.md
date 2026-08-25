# Partner Experience Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Experiência do parceiro — rotas escopadas `/partner` (clientes com regra 24h, maquininhas read-only, relatório personalizado) + frontend com nav/home/páginas por role.

**Architecture:** Subrouter `/partner` (JWT + RoleMiddleware("partner")) reaproveitando repos existentes (`FindClientByPartnerID`, `FindLoanByPartnerID`); vínculo user↔parceiro por email; relatório agregado em memória no use case. Frontend: view por role na home, modo parceiro em clientes, apis novas.

**Tech Stack:** Go/gorm | React+MUI+vitest. **Spec:** `docs/superpowers/specs/2026-08-25-partner-experience-design.md`

**Convenções:** como o plano anterior (TDD vermelho→verde por use case, mocks embed+Func, commits por task, NUNCA commitar `.env`). Backend em `/home/gabigontijo/Documents/cash-by-card`, frontend em `/home/gabigontijo/Documents/cash-by-card-front`, branch `feature/user-access-password`.

---

### Task 1: Contratos + DTOs do parceiro

Create `internal/contracts/ipartner_clients_use_case.go`, `ipartner_report_use_case.go`:

```go
// ipartner_clients_use_case.go
package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/client"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
)

// PartnerClientsUseCase: CRUD de clientes escopado ao parceiro logado.
type PartnerClientsUseCase interface {
	ListClients(ctx context.Context, partnerUserEmail string) (*output.PartnerClients, error)
	CreateClient(ctx context.Context, partnerUserEmail string, createClient *input.CreateClient) (*output.PartnerCreateClient, error)
	UpdateClient(ctx context.Context, partnerUserEmail string, updateClient *input.UpdateClient) (*output.PartnerUpdateClient, error)
}

// ipartner_report_use_case.go
package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
)

type PartnerReportUseCase interface {
	Execute(ctx context.Context, partnerUserEmail string) (*output.PartnerReport, error)
}
```

Create `internal/ports/output/partner/` package:

```go
// clients.go
package output

import "github.com/jhmorais/cash-by-card/internal/domain/entities"

type PartnerClients struct {
	Clients []*entities.Client `json:"clients"`
}

type PartnerCreateClient struct {
	ClientID int `json:"clientId"`
}

type PartnerUpdateClient struct {
	ClientID int `json:"clientId"`
}
```

```go
// report.go
package output

import "time"

type PartnerReportSummary struct {
	TotalLoans      int     `json:"totalLoans"`
	TotalCommission float64 `json:"totalCommission"`
}

type PartnerMonth struct {
	Month      int     `json:"month"` // 1-12
	Loans      int     `json:"loans"`
	Commission float64 `json:"commission"`
}

type PartnerYear struct {
	Year   int            `json:"year"`
	Months []PartnerMonth `json:"months"` // sempre 12 posições
}

type PartnerMonthDetail struct {
	LoanID      int       `json:"loanId"`
	Commission  float64   `json:"commission"`
	CreatedAt   time.Time `json:"createdAt"`
	ClientName  string    `json:"clientName"`
}

type PartnerReport struct {
	Summary      PartnerReportSummary `json:"summary"`
	Annual       []PartnerYear        `json:"annual"`
	CurrentMonth []PartnerMonthDetail `json:"currentMonth"`
	GeneratedAt  time.Time            `json:"generatedAt"`
}
```

Verify: `go build ./...`. Commit: `feat: contratos e dtos da experiencia do parceiro`.

### Task 2: Use cases de clientes do parceiro (TDD)

Create `internal/usecases/partner/partner_clients_use_case.go` + `_test.go`. Package `partner` com mocks próprios (embed interfaces; NÃO reusar mocks do package user).

```go
// partner_clients_use_case.go
package partner

import (
	"context"
	"fmt"
	"time"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/client"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
	repoClient "github.com/jhmorais/cash-by-card/internal/repositories/client"
	repoPartner "github.com/jhmorais/cash-by-card/internal/repositories/partner"
)

const partnerEditWindow = 24 * time.Hour

type partnerClientsUseCase struct {
	clientRepository  repoClient.ClientRepository
	partnerRepository repoPartner.PartnerRepository
}

func NewPartnerClientsUseCase(clientRepository repoClient.ClientRepository, partnerRepository repoPartner.PartnerRepository) contracts.PartnerClientsUseCase {
	return &partnerClientsUseCase{clientRepository: clientRepository, partnerRepository: partnerRepository}
}

// resolvePartner devolve a entidade parceira pelo email do user; nil, nil quando não existe.
func (p *partnerClientsUseCase) resolvePartner(ctx context.Context, email string) (*entities.Partner, error) {
	partner, err := p.partnerRepository.FindPartnerByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if partner == nil || partner.Email == "" {
		return nil, nil
	}
	return partner, nil
}

func (p *partnerClientsUseCase) ListClients(ctx context.Context, partnerUserEmail string) (*output.PartnerClients, error) {
	partner, err := p.resolvePartner(ctx, partnerUserEmail)
	if err != nil {
		return nil, err
	}
	if partner == nil {
		return &output.PartnerClients{Clients: []*entities.Client{}}, nil
	}
	clients, err := p.clientRepository.FindClientByPartnerID(ctx, partner.ID, "")
	if err != nil {
		return nil, err
	}
	if clients == nil {
		clients = []*entities.Client{}
	}
	return &output.PartnerClients{Clients: clients}, nil
}

func (p *partnerClientsUseCase) CreateClient(ctx context.Context, partnerUserEmail string, createClient *input.CreateClient) (*output.PartnerCreateClient, error) {
	partner, err := p.resolvePartner(ctx, partnerUserEmail)
	if err != nil {
		return nil, err
	}
	if partner == nil {
		return nil, fmt.Errorf("nenhum parceiro vinculado a este usuário")
	}
	// vínculo automático: partner_id do body é ignorado
	partnerID := partner.ID
	entity := &entities.Client{
		Name:      createClient.Name,
		PixType:   createClient.PixType,
		PixKey:    createClient.PixKey,
		Phone:     createClient.Phone,
		CPF:       createClient.CPF,
		Documents: createClient.Documents,
		PartnerID: &partnerID,
		CreatedAt: time.Now(),
	}
	if err := p.clientRepository.CreateClient(ctx, entity); err != nil {
		return nil, err
	}
	return &output.PartnerCreateClient{ClientID: entity.ID}, nil
}

func (p *partnerClientsUseCase) UpdateClient(ctx context.Context, partnerUserEmail string, updateClient *input.UpdateClient) (*output.PartnerUpdateClient, error) {
	partner, err := p.resolvePartner(ctx, partnerUserEmail)
	if err != nil {
		return nil, err
	}
	if partner == nil {
		return nil, fmt.Errorf("nenhum parceiro vinculado a este usuário")
	}
	client, err := p.clientRepository.FindClientByID(ctx, updateClient.ID)
	if err != nil || client == nil || client.ID == 0 {
		return nil, fmt.Errorf("cliente não encontrado")
	}
	if client.PartnerID == nil || *client.PartnerID != partner.ID {
		return nil, fmt.Errorf("sem permissão para editar este cliente")
	}
	if time.Since(client.CreatedAt) > partnerEditWindow {
		return nil, fmt.Errorf("clientes só podem ser editados nas primeiras 24h após a criação; procure um administrador")
	}
	client.Name = updateClient.Name
	client.PixType = updateClient.PixType
	client.PixKey = updateClient.PixKey
	client.Phone = updateClient.Phone
	client.CPF = updateClient.CPF
	client.Documents = updateClient.Documents
	if err := p.clientRepository.UpdateClient(ctx, client); err != nil {
		return nil, err
	}
	return &output.PartnerUpdateClient{ClientID: client.ID}, nil
}
```

Tests (escrever PRIMEIRO, ver falhar, depois implementar) cobrindo:
- `TestListClients_Escopo`: FindClientByPartnerID chamado com o ID da entidade resolvida por email; email sem entidade → lista vazia sem erro, repo não chamado
- `TestCreateClient_VinculoAutomatico`: cria com PartnerID = entidade dele MESMO se o body mandar outro partnerId (input.CreateClient não tem PartnerID usado)
- `TestUpdateClient_Regra24h`: client criado há 23h → edita; há 25h → erro contendo "24h"; client de outro parceiro → "sem permissão"; boundary: `partnerEditWindow - time.Minute` ok, `partnerEditWindow + time.Minute` bloqueia
- Partner não resolvido no create/update → erro "nenhum parceiro vinculado"

Mocks: `mockPartnerRepo{findByEmail func}`, `mockClientRepo{byPartnerID, create, byID, update funcs}` — embeddings das interfaces reais.

`go test ./internal/usecases/partner/ -v` verde; `go build ./...`; suite completa verde. Commit: `feat: use cases de clientes do parceiro com regra de 24h`.

### Task 3: Use case do relatório (TDD)

`internal/usecases/partner/partner_report_use_case.go`:

```go
package partner

import (
	"context"
	"time"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
	repoLoan "github.com/jhmorais/cash-by-card/internal/repositories/loan"
	repoPartner "github.com/jhmorais/cash-by-card/internal/repositories/partner"
)

type partnerReportUseCase struct {
	loanRepository    repoLoan.LoanRepository
	partnerRepository repoPartner.PartnerRepository
}

func NewPartnerReportUseCase(loanRepository repoLoan.LoanRepository, partnerRepository repoPartner.PartnerRepository) contracts.PartnerReportUseCase {
	return &partnerReportUseCase{loanRepository: loanRepository, partnerRepository: partnerRepository}
}

func (p *partnerReportUseCase) Execute(ctx context.Context, partnerUserEmail string) (*output.PartnerReport, error) {
	report := &output.PartnerReport{
		Summary:      output.PartnerReportSummary{},
		Annual:       []output.PartnerYear{},
		CurrentMonth: []output.PartnerMonthDetail{},
		GeneratedAt:  time.Now(),
	}
	partner, err := p.partnerRepository.FindPartnerByEmail(ctx, partnerUserEmail)
	if err != nil {
		return nil, err
	}
	if partner == nil || partner.Email == "" {
		return report, nil // sem entidade: relatório vazio
	}
	loans, err := p.loanRepository.FindLoanByPartnerID(ctx, partner.ID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	type key struct{ year, month int }
	agg := map[key]*output.PartnerMonth{}
	for _, l := range loans {
		created := l.CreatedAt
		report.Summary.TotalLoans++
		report.Summary.TotalCommission += l.PartnerAmount

		k := key{created.Year(), int(created.Month())}
		m, ok := agg[k]
		if !ok {
			m = &output.PartnerMonth{Month: k.month}
			agg[k] = m
		}
		m.Loans++
		m.Commission += l.PartnerAmount

		if created.Year() == now.Year() && int(created.Month()) == int(now.Month()) {
			clientName := ""
			if l.Client.Name != "" {
				clientName = l.Client.Name
			}
			report.CurrentMonth = append(report.CurrentMonth, output.PartnerMonthDetail{
				LoanID:     l.ID,
				Commission: l.PartnerAmount,
				CreatedAt:  created,
				ClientName: clientName,
			})
		}
	}

	// série anual completa: do menor ano com loan até o ano atual, 12 meses cada
	if len(agg) > 0 {
		minYear := now.Year()
		for k := range agg {
			if k.year < minYear {
				minYear = k.year
			}
		}
		for y := minYear; y <= now.Year(); y++ {
			py := output.PartnerYear{Year: y, Months: make([]output.PartnerMonth, 12)}
			for m := 1; m <= 12; m++ {
				py.Months[m-1] = output.PartnerMonth{Month: m}
				if agg[key{y, m}] != nil {
					py.Months[m-1] = *agg[key{y, m}]
				}
			}
			report.Annual = append(report.Annual, py)
		}
	}
	return report, nil
}
```

NOTA: o uso de `l.Client.Name` exige que `FindLoanByPartnerID` faça Preload das associações — VERIFICAR em `internal/repositories/loan/loan_repository.go`; se não tiver `Preload(clause.Associations)`, adicionar (igual aos outros finders).

Tests (primeiro): loans espalhados por 2 anos → summary correto (soma PartnerAmount), annual com série completa (anos intermediários, 12 meses, zeros onde sem loan), currentMonth só com loans do ano/mês ATUAL (usar loans com `CreatedAt: time.Now()` e outros com datas antigas), clientName preenchido, email sem entidade → relatório vazio.

`go test ./internal/usecases/partner/ -v` verde. Commit: `feat: use case de relatorio do parceiro`.

### Task 4: Rotas, handlers e DI

`services/partner.go` (novo) — handlers `PartnerListClients`, `PartnerCreateClient`, `PartnerUpdateClient`, `PartnerReport`, `PartnerCardMachines` seguindo o padrão dos handlers existentes (requesterEmail via `utils.EmailFromContext(r.Context())`; erros de permissão/24h → 403 via `writeUserError` do services — é exportado no pacote; senão replicar a helper local). `PartnerCardMachines` chama o `ListCardMachinesUseCase` JÁ EXISTENTE no Handler.

`services/rest_server.go`: Handler struct += `PartnerClientsUseCase contracts.PartnerClientsUseCase` e `PartnerReportUseCase contracts.PartnerReportUseCase` (+ wiring no literal); rotas:

```go
	partnerRouter := router.PathPrefix("/partner").Subrouter()
	partnerRouter.Use(utils.ValidateJwtTokenMiddleware)
	partnerRouter.Use(utils.RoleMiddleware(userRepo.UserRepository, "partner"))
	partnerRouter.HandleFunc("/clients", handler.PartnerListClients).Methods(http.MethodGet)
	partnerRouter.HandleFunc("/clients", handler.PartnerCreateClient).Methods(http.MethodPost)
	partnerRouter.HandleFunc("/clients/{id}", handler.PartnerUpdateClient).Methods(http.MethodPut)
	partnerRouter.HandleFunc("/card-machines", handler.PartnerCardMachines).Methods(http.MethodGet)
	partnerRouter.HandleFunc("/report", handler.PartnerReport).Methods(http.MethodGet)
```

DI (`dependency_builder.go`): `Usecases.PartnerClientsUseCase: partner.NewPartnerClientsUseCase(builder.Repositories.ClientRepository, builder.Repositories.PartnerRepository)` e `PartnerReportUseCase: partner.NewPartnerReportUseCase(builder.Repositories.LoanRepository, builder.Repositories.PartnerRepository)` (import `usecases/partner`).

Verificar: `go build ./...`, `go test ./... -count=1` verde. Commit: `feat: rotas /partner com escopo por role partner`.

### Task 5: E2E backend (curl)

Reiniciar server (`make run` no terminal do usuário — ou temporariamente em background p/ verificação e depois liberar a porta). Cenários (seed: criar entidade parceira com email igual a um user partner via SQL direto + loans apontando partner_id dela):

1. login partner → `GET /partner/clients` → só clientes dele; email sem entidade → `{"clients":[]}`
2. `POST /partner/clients` → cria vinculado (conferir `partner_id` no banco); body com partnerId malicioso → ignorado
3. `PUT /partner/clients/{id}` recém-criado → 200; simulando 24h passado (UPDATE client SET created_at = NOW() - INTERVAL 25 HOUR) → erro com "24h"
4. `GET /partner/report` → summary/annual/currentMonth coerentes com o seed; partner sem loans → tudo vazio
5. partner em `/admin/clients` → 403 (regressão)
6. admin em `/partner/report` → 403 (só partner)

Sem commit (verificação). Registrar saídas no relatório da task.

### Task 6: Frontend — apis + tests

`src/apis/partner/index.js` (padrão apiFetch):

```js
import { apiFetch } from '..';

export const partnerClients = async () => {
  const res = await apiFetch('partner/clients', { method: 'get' });
  return res.json();
};

export const partnerCreateClient = async (client) => {
  const res = await apiFetch('partner/clients', { method: 'post', body: JSON.stringify(client) });
  return res.json();
};

export const partnerUpdateClient = async (id, client) => {
  const res = await apiFetch(`partner/clients/${id}`, { method: 'put', body: JSON.stringify(client) });
  return res.json();
};

export const partnerCardMachines = async () => {
  const res = await apiFetch('partner/card-machines', { method: 'get' });
  return res.json();
};

export const partnerReport = async () => {
  const res = await apiFetch('partner/report', { method: 'get' });
  return res.json();
};
```

`src/apis/partner/index.test.js` no padrão dos apis (stub fetch/env/token; 5 testes: URLs+methods+bearer). `npm test` verde; eslint limpo. Commit: `feat: apis do parceiro`.

### Task 7: Frontend — PartnerReportView + home por role

`src/sections/overview/view/partner-report-view.jsx` (Container + h4 "Meus empréstimos"): cards de resumo (Total de empréstimos / Comissão total, formatando valor com `Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' })`); tabela anual (Ano | Mês | Empréstimos | Comissão — mês como nome pt-BR); tabela mês atual (Data | Cliente | Comissão — data `date-fns` format dd/MM/yyyy); TODOS os blocos com estado vazio "Sem empréstimos feitos para o período"; loading simples. `src/pages/app.jsx`: `const { role } = useAuth();` → role === 'partner' ? `<PartnerReportView />` : `<DashboardView />` (título Helmet por role). `src/sections/overview/view/partner-report-view.test.jsx`: fetch stub retornando report fixture → resumo renderizado; fixture vazia → mensagens de estado vazio. `npm test` + build. Commit: `feat: relatorio personalizado do parceiro na home`.

### Task 8: Frontend — clientes modo parceiro, simulação e nav

- `src/sections/client/view/client-view.jsx` (LER PRIMEIRO): aceitar modo via `useAuth().role === 'partner'` → usar `partnerClients/partnerCreateClient/partnerUpdateClient` nas chamadas; OCULTAR botões editar (renderizar apenas quando `Date.now() - new Date(c.created_at? c.CreatedAt) < 24h` — campo json `createdAt`) e excluir (nunca para partner). O backend continua a autoridade.
- `src/sections/simulation/view/simulation-view.jsx`: trocar `allCardMachines` por `partnerCardMachines` quando role partner (buscar role via useAuth).
- `src/layouts/dashboard/nav.jsx`: filtro — partner vê apenas `/` (relatórios), `/cliente`, `/simulacoes`, `/alterar-senha`; esconder `/emprestimo`, `/parceiro`, `/maquininha` para partner (além do `/usuarios` já existente).
- `npm test` + `npm run build`. Commit: `feat: modo parceiro em clientes, simulacao e navegacao`.

### Task 9: E2E completo + review final

Backend+frontend rodando; cenário manual no browser com o user partner (seed da Task 5): home mostra relatório; clientes só dele; adicionar/editar(<24h)/sem excluir; simulações carrega maquininhas; nav correta; admin/organization inalterados. `go test ./... && npm test` finais; review final (subagent) do diff completo da feature parceiro contra a spec; atualizar `.claude/skills/run-dev/SKILL.md` se houver novo gotcha.

---

## Self-review do plano

- Cobertura da spec: rotas ✓ (6, incluindo card-machines), 24h ✓, vínculo email ✓, report 3 blocos + série completa + vazios ✓, home por role ✓, nav ✓, simulação ✓, testes TDD ✓, E2E ✓
- Tipos: outputs definidos na Task 1 usados nas Tasks 2-3; construtores batem com DI da Task 4
- Sem placeholders: código completo nas tasks backend; frontend com código nas apis e instruções concretas+checklist nas views (seguindo o padrão já executado no plano anterior)
