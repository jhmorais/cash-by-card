# Design: Experiência do parceiro (acesso restrito + relatório personalizado)

- **Data**: 2026-08-25
- **Branch**: `feature/user-access-password` (continuação — mesma linha de trabalho de acesso/roles)
- **Status**: aprovado pelo usuário

## Contexto

Hoje o role `partner` só acessa `/auth/*` e `/account/*` (alterar senha). Este design dá ao
parceiro uma experiência completa e escopada: ver e adicionar seus clientes, usar simulações,
e um relatório personalizado dos empréstimos em que ele foi parceiro — sem acesso a nada além
do que lhe é relacionado.

## Decisões (validadas com o usuário)

| Tema | Decisão |
|---|---|
| Vínculo user↔parceiro | Por email: `partner.email = user.email`. Sem entidade parceira com o email → parceiro sem dados (estados vazios, sem erro) |
| Clientes do parceiro | Apenas os vinculados a ele (`client.partner_id`) |
| Editar cliente (partner) | Só nas primeiras **24h** após `created_at`; depois, apenas admin/organization |
| Excluir cliente | Partner **nunca** (rota nem existe para ele) |
| Superfície de API | Subrouter novo `/partner` com `ValidateJwtTokenMiddleware` + `RoleMiddleware("partner")` — rotas `/admin` intocadas |
| Comissão | `loan.partner_amount` (valor destinado ao parceiro no empréstimo) |

## Rotas novas (backend, TDD)

| Rota | Regra |
|---|---|
| `GET /partner/clients` | Lista clientes com `partner_id` = entidade do parceiro logado |
| `POST /partner/clients` | Cria cliente **já vinculado a ele** — `partner_id` do body é ignorado/sobrescrito |
| `PUT /partner/clients/{id}` | Edita se: cliente é dele E `created_at` há menos de 24h. Senão → 403 com mensagem "clientes só podem ser editados nas primeiras 24h; procure um administrador" |
| `DELETE /partner/clients/{id}` | Não existe |
| `GET /partner/card-machines` | Lista read-only de maquininhas (insumo da simulação) |
| `GET /partner/report` | Relatório completo (abaixo) |

## Relatório do parceiro — `GET /partner/report`

Resposta única, três blocos, tudo filtrado por `loan.partner_id` = entidade dele:

```json
{
  "summary": { "totalLoans": 0, "totalCommission": 0 },
  "annual": [ { "year": 2026, "months": [ { "month": 1, "loans": 0, "commission": 0 } ] } ],
  "currentMonth": [ { "loanId": 1, "commission": 0, "createdAt": "...", "clientName": "..." } ],
  "generatedAt": "..."
}
```

- **summary**: total de empréstimos em que foi parceiro + soma de `partner_amount`
- **annual**: por ano (desde o primeiro loan dele) → por mês (1-12): nº de empréstimos + comissão
  somada do mês. Meses sem atividade vêm com zeros (série completa)
- **currentMonth**: mês corrente, um item por empréstimo: `partnerAmount`, `created_at` do loan,
  nome do cliente (via `client_id`)
- **Sem dados**: blocos vazios (`summary` zerado, `annual` vazio, `currentMonth` vazio) — o
  frontend mostra "Sem empréstimos feitos para o período"

## Resolução do parceiro logado (helper)

`ResolvePartnerByEmail(ctx, email)`: `FindPartnerByEmail`; devolve `nil, nil` quando não há
entidade. Rotas `/partner` com parceiro não resolvido: clientes → lista vazia + create falha
com erro explícito ("nenhum parceiro vinculado a este usuário"); report → resposta vazia.

## Frontend

- **Nav por role** (nav.jsx): partner vê apenas *relatórios* (`/`), *clientes*, *simulações*,
  *alterar senha*. Empréstimos/parceiros/maquininhas somem para partner
- **Home `/`**: por role — partner → nova `PartnerReportView`; org/admin → dashboard atual
  (mesma rota, troca de view no `IndexPage`)
- **PartnerReportView**: cards de resumo (total de empréstimos + comissão total), série anual
  (tabela ano/mês: nº loans + comissão), detalhamento do mês atual (tabela: comissão, data,
  cliente) e estados vazios com a mensagem acordada
- **Clientes em modo parceiro**: reusa a página existente com modo `partner`: sem botão
  excluir; editar aparece apenas quando `created_at` < 24h (o backend é a autoridade — o
  botão oculto é UX, o 403 é garantia). Chamadas trocadas para `/partner/clients`
- **Simulações**: quando `auth.role === 'partner'`, busca maquininhas em
  `GET /partner/card-machines`; demais roles seguem em `/admin/card-machines`
- Novas apis em `src/apis/partner/` seguindo o padrão `apiFetch`

## Segurança

- Toda regra de escopo no backend ( partner só vê o dele via SQL WHERE; regra 24h no use case)
- Frontend oculta botões apenas por UX — nunca é a barreira
- Partner continua sem acesso a qualquer `/admin/*` (403 pelo RoleMiddleware existente)

## Testes

- **Backend (TDD)**: escopo por email; create com vínculo automático (body ignorado); regra 24h
  (boundary: 23h59 edita, 24h01 bloqueia); report (somas, agrupamento anual com zeros, mês
  atual com dados de teste); parceiro não resolvido
- **Frontend**: vitest nas views/apis novas; nav partner; clients em modo parceiro
- **E2E live**: curl com partner real (seed com loans)

## Fora de escopo (futuro)

- FK real `user.partner_id` (substituiria o vínculo por email)
- Paginação na listagem de clientes do parceiro
- Exportação (PDF/CSV) do relatório
