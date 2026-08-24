---
name: run-dev
description: Sobe o ambiente local do cash-by-card (MySQL no Docker + API Go nativa) e faz smoke test com login JWT — use quando pedir para rodar, testar ou demonstrar a API localmente
---

# Rodar o cash-by-card localmente (WSL2)

API REST Go + MySQL 5.7 no Docker. Sem frontend. Porta **3000**.

## Pré-condições (já instaladas nesta máquina em 2026-08-21)

- Go 1.24 em `~/.local/go` (symlink `~/.local/bin/go`; `go.mod` exige toolchain 1.24, o apt só tem 1.22)
- Docker (`docker.io` + `docker-compose-v2`). Sessões abertas antes do `usermod` precisam de `sg docker -c "..."` para acessar o socket sem sudo

## Run

```bash
# 1. MySQL (dump inicial em docker/mysql/dump.sql cria schema+6 tabelas no 1º boot)
make compose-up   # ou: sg docker -c "docker compose up -d"

# aguardar healthy
for i in $(seq 1 30); do [ "$(sg docker -c "docker inspect -f '{{.State.Health.Status}}' mysqlcontainer")" = healthy ] && break; sleep 3; done

# 2. API (rodar SEMPRE da raiz do repo — o .env é lido do CWD)
go run cmd/restserver/main.go &> /tmp/cashbycard.log &

# aguardar ready (log: "Starting SERVER, LISTEN PORT: 3000")
for i in $(seq 1 20); do curl -s -o /dev/null --max-time 2 -X POST http://localhost:3000/auth/login && break; sleep 2; done
```

## Smoke test

```bash
JWT=$(curl -s -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@local.dev","password":"admin123"}' | sed 's/^"//; s/"$//')

curl -s http://localhost:3000/admin/loans -H "Authorization: Bearer $JWT"
# → {"loans":[...],"total":N,"page":1,"limit":10,"totalPages":M}
```

Todas as rotas `/admin` e `/public` exigem Bearer token (sem ele: "token contains an invalid
number of segments"). `/auth/login` é a única aberta.

## Stop

```bash
lsof -ti:3000 -sTCP:LISTEN | xargs -r kill   # libera a porta da API
make compose-stop                             # para o MySQL sem perder dados
```

Não use `pkill -f "restserver"` — o padrão casa com a própria linha de comando do shell
que o invocou e mata a sessão (verificado na prática).

## Environment (`.env` na raiz, gitignored — recriar se sumir)

| Variável | Valor local | Nota |
|---|---|---|
| `MYSQL_CONNECTION_STRING` | `tcp(localhost:3306)/database?charset=utf8&parseTime=True&loc=Local` | |
| `MYSQL_USER` / `MYSQL_PASSWORD` | `root` / `password` | |
| `SERVER_PORT` | `3000` | sem isso a app escuta em porta aleatória |
| `JWT_SECRET_KEY` | qualquer valor | sem isso login/auth falham |

Mudou o `.env`? Reinicie a server — as vars são lidas só no boot.

## Gotchas do app (verificados)

- **`go run` não recarrega**: depois de editar código Go, mate o server (port kill acima) e
  suba de novo — senão você testa contra um binário velho e filtros novos "não funcionam".
- Não existe rota de criação de usuário. O admin `admin@local.dev` / `admin123`
  (senha = MD5 hex, role `admin`) foi inserido direto no banco local.
- `paymentType` de card só aceita `onlineTax` ou `presentialTax`.
- `POST /admin/loans` cria o loan ANTES de validar os cards — card inválido
  deixa loan órfão. Limpar com `DELETE /admin/loans/{id}` se necessário.
- Query params de filtro/paginação são camelCase (`paymentStatus`, `clientName`,
  `amountMin`, `page`, `limit` ≤ 100) e só existem em `GET /admin/loans`.
- Reset total do banco: `sg docker -c "docker compose down -v"` + `make compose-up`
  (o dump roda de novo num datadir vazio).
