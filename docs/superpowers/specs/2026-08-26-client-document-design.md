# Design: Documento do cliente (1 por cliente — PDF/JPEG/PNG)

- **Data**: 2026-08-26 | **Branch**: `feature/user-access-password` | **Status**: aprovado

## Decisões

| Tema | Decisão |
|---|---|
| Storage | **Filesystem**: pasta `DOCS_DIR` (env; default `./docs`; volume no servidor), arquivo único `{cpf}.{ext}`, upload **substitui** o existente |
| Formatos | PDF, JPEG, PNG — validação por **magic bytes** (`%PDF-`, `\xFF\xD8`, `\x89PNG`), extensão é ignorada |
| Quantidade | **1 documento por cliente** (upsert) |
| Tamanho máx | 5MB |
| Acesso | download/visualização/remoção **só org/admin** (JWT + RoleMiddleware existentes) |
| Parceiro | sem upload/visualização de documento (gestão de admin) |
| Banco | `client` ganha `document_name`, `document_type`, `document_size` (migration manual + dump) |

## Rotas

| Rota | Comportamento |
|---|---|
| `POST /admin/clients/files/{cpf}` (refeito) | multipart com **exatamente 1 arquivo**; valida magic bytes + 5MB; salva `{cpf}.{ext}` em DOCS_DIR apagando o anterior; atualiza colunas; erros PT ("formato não suportado — use PDF, JPEG ou PNG", "envie apenas um arquivo", "arquivo maior que 5MB", "cliente não encontrado") |
| `GET /admin/clients/{id}/document` (novo) | streama com Content-Type correto + `Content-Disposition: inline`; 404 se não existe |
| `DELETE /admin/clients/{id}/document` (novo) | apaga arquivo + limpa colunas |

Excluir cliente remove o documento junto. O `saveFile` antigo (pasta `assets/`, múltiplos arquivos, sem validação) sai do código.

## Frontend

Formulário do cliente (modo admin): upload de arquivo único com `accept="application/pdf,image/jpeg,image/png"`, feedback de tamanho; quando existe documento, mostra nome/tamanho + botões **Ver** (abre `GET .../document` em nova aba com Bearer via fetch→blob→objectURL) e **Remover**; subir outro substitui. Modo parceiro inalterado (sem upload).

## Testes

Backend TDD: magic bytes (pdf/jpeg/png ok; gif/_txt/heic rejeitados), limite 5MB, multi-arquivo rejeitado, substituição apaga o anterior, metadados gravados, download content-type, delete limpa. Frontend vitest: estados do upload (sem doc/com doc/erro). E2E com PDF e JPEG reais.

## Fora de escopo

Conversão HEIC server-side; thumbnails; versionamento de documentos; paginação.
