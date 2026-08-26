# Client Document Implementation Plan

> **For agentic workers:** execute via superpowers:subagent-driven-development. Spec: `docs/superpowers/specs/2026-08-26-client-document-design.md`. Branch `feature/user-access-password`, TDD, commits por task, NUNCA commitar `.env`/`docs/`.

### Task B1: Backend — documento do cliente (TDD)

1. **Migration** `db/migrations/2026-08-26-client-document.sql` (+ dump.sql + banco local): `ALTER TABLE client ADD COLUMN document_name VARCHAR(255) NULL, ADD COLUMN document_type VARCHAR(50) NULL, ADD COLUMN document_size INT NULL;`
2. **Entity** `entities.Client` += `DocumentName/DocumentType/DocumentSize string,int` com json tags `documentName/documentType/documentSize` (e output PartnerClients herda).
3. **`internal/documents/`** (novo pacote, sem gorm — puro fs, testável):
   - `Store{Dir string}` (de `DOCS_DIR` env, default `./docs`), `NewStore() *Store`
   - `DetectType(head []byte) string` — magic bytes: `%PDF-`→"pdf", `\xFF\xD8`→"jpeg", `\x89PNG`→"png", senão ""
   - `Save(cpf string, head []byte, r io.Reader, size int64) (ext string, err error)` — valida tipo (err "formato não suportado — use PDF, JPEG ou PNG"), escreve `Dir/{cpf}.{ext}` com `os.CreateTemp`+rename (substitui; remove outro `{cpf}.*` anterior)
   - `Path(cpf, ext)`, `Delete(cpf)` (remove todos `{cpf}.*`, tolera ausente)
   - **TDD**: unit tests com bytes reais (PDF header, JPEG FFD8, PNG 8950, GIF rejeitado), substituição troca extensão e apaga o velho
4. **Handler refeito** `services/clients.go` `CreateClientDocuments`: cliente existe (via FindClientByCPF); `ParseMultipartForm(5<<20)`; exatamente 1 arquivo ("envie apenas um arquivo"); size ≤5MB; `documents.Save`; atualiza colunas via repo `UpdateClientDocument(ctx, id, name, type, size)` (novo método map-update); 200 com metadados.
5. **Rotas novas**: `GET /admin/clients/{id}/document` (procura cliente por ID → cpf → stream `Path`, Content-Type por ext, inline; 404 limpo) e `DELETE /admin/clients/{id}/document` (documents.Delete + limpa colunas) em `rest_server.go`.
6. **Delete de cliente** remove documento: no handler/use case de DeleteClient existente, antes de deletar, `documents.Delete(cpf)`.
7. Remover `saveFile`/assets antigos. `go build`/`go test ./... -count=1` verdes. Commit: `feat: documento unico por cliente em filesystem com validacao por magic bytes`.

### Task B2: E2E backend (curl)

Gerar PDF/JPEG mínimos (`printf '%s'` com headers), POST com `-F`, verificar: aceita e substitui (arquivo no DOCS_DIR com uma extensão só), rejeita gif/txt, rejeita 2 arquivos, GET content-type correto, DELETE limpa, delete do cliente remove o arquivo.

### Task F1: Frontend — upload no formulário do cliente

`src/sections/client/form-new-client.jsx` + `src/apis/client/`: upload único (accept `application/pdf,image/jpeg,image/png`), usa `uploadClientFiles` ajustado p/ 1 arquivo; exibe documento atual (`client.documentName/documentSize` da listagem) com botões **Ver** (fetch Bearer → blob → objectURL → `window.open`) e **Remover** (`DELETE /admin/clients/{id}/document`); erros do backend em Alert. vitest nos estados. `npm test` + `npm run build`. Commit: `feat: upload e visualizacao do documento unico do cliente`.

### Task F2: Review final + E2E manual

Subagent review do diff contra a spec; suites finais; entrega para teste do usuário (login org → cliente → subir PDF → ver → substituir por JPEG → remover).
