# Design: Primeiro acesso, recuperação de senha e administração de usuários

- **Data**: 2026-08-24
- **Branch**: `feature/user-access-password` (backend `cash-by-card` e frontend `cash-by-card-front`, ambos a partir da `main`)
- **Status**: aprovado pelo usuário

## Contexto

Hoje o sistema tem uma única tabela `user` com `role = 'admin'`, login via `POST /auth/login`
(senha MD5) e nenhuma gestão de usuários — os registros são criados direto no banco. Este
design introduz três níveis de acesso (**organization**, **admin**, **partner**), uma página
de administração de usuários, e o fluxo de primeiro acesso/recuperação de senha por email
com token de uso único.

O `admin` atual do banco já corresponde ao nível intermediário — nenhuma migração de roles.

## Decisões (validadas com o usuário)

| Tema | Decisão |
|---|---|
| Envio de email | SMTP via interface `EmailSender` (config no `.env`); sem SMTP configurado, o link é logado no console do backend (dev) |
| Roles | `organization`, `admin`, `partner` (substituem a noção de manager; `admin` existente permanece) |
| Branch | Nova branch a partir da `main` (PR da `feature/loans-filters-pagination` já mergeado) |
| Schema em produção | SQL manual: script versionado em `db/migrations/`, aplicado uma vez por SSH/phpMyAdmin no servidor |
| Token | Tabela dedicada `password_reset_token`, hash SHA-256, expiração 30 min, uso único |

## Modelo de dados

`db/migrations/2026-08-24-user-access-password.sql` (novo) e `docker/mysql/dump.sql` atualizado:

```sql
ALTER TABLE `user` MODIFY `password` VARCHAR(100) NULL;  -- NULL = pendente de primeiro acesso

CREATE TABLE `password_reset_token` (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  token_hash VARCHAR(64) NOT NULL,       -- SHA-256 hex do token enviado por email
  expires_at TIMESTAMP NOT NULL,         -- 30 minutos a partir da criação
  used_at TIMESTAMP NULL,                -- preenchido no uso; single-use
  created_at TIMESTAMP NOT NULL,
  FOREIGN KEY (user_id) REFERENCES `user` (`id`)
);
```

## Fluxos

### Criar usuário (org/admin na página de administração)
1. `POST /admin/users {name, email, role}` → valida permissão e duplicidade de email
2. Grava com `password = NULL` (pendente de primeiro acesso)
3. Gera token e envia email com link `FRONT_URL/primeiro-acesso?token=<64 hex>`

### Primeiro acesso / Esqueci minha senha (mesma mecânica, entradas visuais distintas)
1. Página de login expõe dois links — **"Primeiro acesso"** e **"Esqueci minha senha"** — ambos
   abrem um input de email e chamam `POST /auth/forgot-password {email}`
2. Backend envia token **sempre que o email existir** (`password NULL` = primeiro acesso,
   senha definida = recuperação); a distinção entre os casos é apenas visual no frontend
3. Resposta é sempre 200 com mensagem genérica ("se o email estiver cadastrado, você
   receberá as instruções") — não revela existência do email (anti-enumeração)
4. Tokens anteriores permanecem válidos até expirarem ou serem usados (sem invalidação
   cruzada — o cenário "token expirou, pedir outro" é coberto pela expiração natural)
5. `POST /auth/reset-password {token, newPassword}` valida hash + expiração + `used_at`,
   grava a senha e marca o token como usado

### Login com `password NULL`
Retorna erro específico "usuário pendente de primeiro acesso" — não entra, orienta ao fluxo.

### Trocar a própria senha (logado, qualquer role)
`POST /auth/change-password {currentPassword, newPassword}` — valida a senha atual.

### Limpar senha (org/admin)
`POST /admin/users/{id}/clear-password` → `password = NULL` + email com novo token.
O usuário volta ao estado de primeiro acesso.

## Endpoints

| Rota | Auth | Quem |
|---|---|---|
| `POST /auth/forgot-password` | pública | sempre 200 genérico |
| `POST /auth/reset-password` | pública | valida token |
| `POST /auth/change-password` | JWT | qualquer role (a própria senha) — rota física: `POST /account/change-password` |
| `GET /auth/me` | JWT | qualquer role — `{email, name, role}` para o frontend montar a nav — rota física: `GET /account/me` |
| `GET /admin/users` | JWT | organization + admin |
| `POST /admin/users` | JWT | org: role `admin` ou `partner`; admin: só `partner` |
| `PUT /admin/users/{id}` | JWT | org: qualquer user; admin: só partner — edita `name`/`role`; **email não é editável** |
| `DELETE /admin/users/{id}` | JWT | org: qualquer user; admin: só partner — **nunca a si mesmo** |
| `POST /admin/users/{id}/clear-password` | JWT | org: qualquer user; admin: só partner — nunca a si mesmo |

O `RoleMiddleware("admin")` atual passa a aceitar `organization` e `admin` nas rotas
`/admin` existentes. `partner` só acessa `/auth/*` (públicas) e `/account/*` (JWT,
próprias do usuário). As rotas autenticadas de conta ficam sob o prefixo `/account`
(subrouter com `ValidateJwtTokenMiddleware`) porque o prefixo `/auth` é público —
daí `GET /account/me` e `POST /account/change-password`.

## Matriz de permissões

| Ação | organization | admin | partner |
|---|---|---|---|
| Listar usuários | ✓ | ✓ | ✗ |
| Adicionar `admin` | ✓ | ✗ | ✗ |
| Adicionar `partner` | ✓ | ✓ | ✗ |
| Editar usuário | qualquer | só partner | ✗ |
| Excluir usuário | qualquer (não a si) | só partner (não a si) | ✗ |
| Limpar senha | qualquer (não a si) | só partner (não a si) | ✗ |
| Trocar a própria senha | ✓ | ✓ | ✓ |
| Demais páginas `/admin` (loans, clients…) | ✓ | ✓ | ✗ |

## Email

Interface `EmailSender` com `SendPasswordResetEmail(to, link)`; implementação SMTP via
`SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASSWORD`, `SMTP_FROM` no `.env`. Sem essas
variáveis, a implementação de dev loga o link no console — o fluxo inteiro é testável
localmente sem credenciais. O link usa `FRONT_URL` (ex. `http://localhost:5173`).

A senha continua **MD5** (padrão atual de login) por consistência; migrar para bcrypt é
melhoria futura registrada abaixo.

## Frontend (`cash-by-card-front`, React + MUI + Vite)

- **`/primeiro-acesso?token=`** (pública): nova senha + confirmação → sucesso → redireciona ao login
- **`/alterar-senha`** (logado): senha atual + nova — todos os roles
- **`/usuarios`** (org + admin): form de adicionar (roles ofertadas conforme o próprio role)
  + tabela (nome, email, role, status da senha "pendente/ok") com ações editar, limpar senha,
  excluir
- **`/login`**: links "Primeiro acesso" e "Esqueci minha senha" com input de email
  (mesmo endpoint, mensagens idênticas)
- **Nav**: item "Usuários" visível só para org + admin; `partner` logado vê apenas
  "Alterar senha" (as demais páginas chamam APIs `/admin` que o backend rejeita)
- APIs novas em `src/apis/user/`, seguindo o padrão `apiFetch`

## Segurança

- Token: 32 bytes aleatórios (crypto/rand), 64 hex no email; no banco apenas SHA-256
- Expiração 30 min; single-use (`used_at`); validação por hash — token vazado do banco não
  serve para nada
- `forgot-password` não revela existência de email; sem rate limit nesta fase (ver futuro)
- Auto-exclusão e auto-limpeza de senha bloqueadas no backend
- MD5 mantido apenas por consistência com o login existente (ver melhorias futuras)

## Testes

- **Backend (TDD)**: use cases com mock de repositório (validações, matriz de permissões,
  token expirado/usado/inválido, login com password NULL); handlers com httptest
- **Frontend**: vitest nos componentes novos, seguindo o padrão de
  `src/sections/loan/loan-filters-panel.test.jsx`
- **Verificação ao vivo**: fluxo completo via curl (criar usuário → token no console →
  reset → login → troca de senha → limpar senha → novo token)

## Fora de escopo (futuro)

- Migração de senha MD5 → bcrypt
- Rate limit em `forgot-password` (anti-spam de emails)
- Purge automático de tokens expirados (coluna `used_at`/`expires_at` já permite limpeza manual)
- Escopo de dados por partner nas páginas existentes (hoje partner não acessa nada além de `/auth/*`)
