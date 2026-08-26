package services

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/client"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/client"
	"github.com/jhmorais/cash-by-card/internal/documents"
)

var (
	docPdfHead = []byte("%PDF-1.4\n%\xe2\xe3\xcf\xd3")
	docGifHead = []byte("GIF89a some gif payload")
)

type uploadFile struct {
	filename string
	content  []byte
}

// mockClientRepository embute a interface e sobrescreve só o que importa.
type mockClientRepository struct {
	repositories.ClientRepository
	findByCPFFunc func(ctx context.Context, cpf string) ([]*entities.Client, error)
	updateDocFunc func(ctx context.Context, id int, name, docType string, size int) error
}

func (m *mockClientRepository) FindClientByCPF(ctx context.Context, cpf string) ([]*entities.Client, error) {
	return m.findByCPFFunc(ctx, cpf)
}

func (m *mockClientRepository) UpdateClientDocument(ctx context.Context, id int, name, docType string, size int) error {
	return m.updateDocFunc(ctx, id, name, docType, size)
}

type mockFindClientByIDUseCase struct {
	contracts.FindClientByIDUseCase
	executeFunc func(ctx context.Context, clientID int) (*output.FindClient, error)
}

func (m *mockFindClientByIDUseCase) Execute(ctx context.Context, clientID int) (*output.FindClient, error) {
	return m.executeFunc(ctx, clientID)
}

type mockDeleteClientUseCase struct {
	contracts.DeleteClientUseCase
	executeFunc func(ctx context.Context, clientID int) (*output.DeleteClient, error)
}

func (m *mockDeleteClientUseCase) Execute(ctx context.Context, clientID int) (*output.DeleteClient, error) {
	return m.executeFunc(ctx, clientID)
}

func documentsTestHandler(repo repositories.ClientRepository, store *documents.Store, findByID *mockFindClientByIDUseCase, deleteUC *mockDeleteClientUseCase) *Handler {
	h := &Handler{
		ClientRepository: repo,
		Documents:        store,
	}
	if findByID != nil {
		h.FindClientByIDUseCase = findByID
	}
	if deleteUC != nil {
		h.DeleteClientUseCase = deleteUC
	}
	return h
}

func multipartUploadRequest(t *testing.T, target string, vars map[string]string, files map[string]uploadFile) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	for field, f := range files {
		fw, err := mw.CreateFormFile(field, f.filename)
		if err != nil {
			t.Fatalf("CreateFormFile falhou: %v", err)
		}
		if _, err := fw.Write(f.content); err != nil {
			t.Fatalf("write falhou: %v", err)
		}
	}
	mw.Close()

	req := httptest.NewRequest(http.MethodPost, target, &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	return mux.SetURLVars(req, vars)
}

const testCPF = "123.456.789-09"

func clientFound(id int) []*entities.Client {
	return []*entities.Client{{ID: id, Name: "Maria", CPF: testCPF}}
}

func TestCreateClientDocuments_Success(t *testing.T) {
	store := &documents.Store{Dir: t.TempDir()}
	content := append(append([]byte{}, docPdfHead...), []byte("conteudo do pdf")...)

	var gotID int
	var gotName, gotType string
	var gotSize int
	repo := &mockClientRepository{
		findByCPFFunc: func(ctx context.Context, cpf string) ([]*entities.Client, error) {
			return clientFound(7), nil
		},
		updateDocFunc: func(ctx context.Context, id int, name, docType string, size int) error {
			gotID, gotName, gotType, gotSize = id, name, docType, size
			return nil
		},
	}
	handler := documentsTestHandler(repo, store, nil, nil)

	req := multipartUploadRequest(t, "/admin/clients/files/"+testCPF,
		map[string]string{"cpf": testCPF},
		map[string]uploadFile{"file": {filename: "contrato.pdf", content: content}})
	rec := httptest.NewRecorder()
	handler.CreateClientDocuments(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	// arquivo salvo com extensão detectada e conteúdo íntegro
	saved, err := os.ReadFile(store.Path(testCPF, "pdf"))
	if err != nil {
		t.Fatalf("arquivo não foi salvo: %v", err)
	}
	if !bytes.Equal(saved, content) {
		t.Fatalf("conteúdo salvo difere (%d bytes)", len(saved))
	}

	if gotID != 7 || gotName != "contrato.pdf" || gotType != "pdf" || gotSize != len(content) {
		t.Fatalf("UpdateClientDocument chamado com (%d, %q, %q, %d)", gotID, gotName, gotType, gotSize)
	}

	var body struct {
		DocumentName string `json:"documentName"`
		DocumentType string `json:"documentType"`
		DocumentSize int    `json:"documentSize"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("resposta não é json válido: %v", err)
	}
	if body.DocumentName != "contrato.pdf" || body.DocumentType != "pdf" || body.DocumentSize != len(content) {
		t.Fatalf("metadados da resposta errados: %+v", body)
	}
}

func TestCreateClientDocuments_MultipleFilesRejected(t *testing.T) {
	store := &documents.Store{Dir: t.TempDir()}
	repo := &mockClientRepository{
		findByCPFFunc: func(ctx context.Context, cpf string) ([]*entities.Client, error) {
			return clientFound(7), nil
		},
	}
	handler := documentsTestHandler(repo, store, nil, nil)

	req := multipartUploadRequest(t, "/admin/clients/files/"+testCPF,
		map[string]string{"cpf": testCPF},
		map[string]uploadFile{
			"file1": {filename: "a.pdf", content: docPdfHead},
			"file2": {filename: "b.pdf", content: docPdfHead},
		})
	rec := httptest.NewRecorder()
	handler.CreateClientDocuments(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "envie apenas um arquivo") {
		t.Fatalf("mensagem inesperada: %s", rec.Body.String())
	}
	entries, err := os.ReadDir(store.Dir)
	if err != nil {
		t.Fatalf("ler docs dir: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("nenhum arquivo deveria ter sido salvo, achei %d", len(entries))
	}
}

func TestCreateClientDocuments_InvalidFormatRejected(t *testing.T) {
	store := &documents.Store{Dir: t.TempDir()}
	repo := &mockClientRepository{
		findByCPFFunc: func(ctx context.Context, cpf string) ([]*entities.Client, error) {
			return clientFound(7), nil
		},
	}
	handler := documentsTestHandler(repo, store, nil, nil)

	req := multipartUploadRequest(t, "/admin/clients/files/"+testCPF,
		map[string]string{"cpf": testCPF},
		map[string]uploadFile{"file": {filename: "foto.gif", content: docGifHead}})
	rec := httptest.NewRecorder()
	handler.CreateClientDocuments(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "formato não suportado") {
		t.Fatalf("mensagem inesperada: %s", rec.Body.String())
	}
}

func TestCreateClientDocuments_TooLargeRejected(t *testing.T) {
	store := &documents.Store{Dir: t.TempDir()}
	repo := &mockClientRepository{
		findByCPFFunc: func(ctx context.Context, cpf string) ([]*entities.Client, error) {
			return clientFound(7), nil
		},
	}
	handler := documentsTestHandler(repo, store, nil, nil)

	content := make([]byte, 5*1024*1024+1)
	copy(content, docPdfHead)

	req := multipartUploadRequest(t, "/admin/clients/files/"+testCPF,
		map[string]string{"cpf": testCPF},
		map[string]uploadFile{"file": {filename: "grande.pdf", content: content}})
	rec := httptest.NewRecorder()
	handler.CreateClientDocuments(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "arquivo maior que 5MB") {
		t.Fatalf("mensagem inesperada: %s", rec.Body.String())
	}
}

func TestCreateClientDocuments_ClientNotFound(t *testing.T) {
	store := &documents.Store{Dir: t.TempDir()}
	repo := &mockClientRepository{
		findByCPFFunc: func(ctx context.Context, cpf string) ([]*entities.Client, error) {
			return nil, nil
		},
	}
	handler := documentsTestHandler(repo, store, nil, nil)

	req := multipartUploadRequest(t, "/admin/clients/files/"+testCPF,
		map[string]string{"cpf": testCPF},
		map[string]uploadFile{"file": {filename: "a.pdf", content: docPdfHead}})
	rec := httptest.NewRecorder()
	handler.CreateClientDocuments(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "cliente não encontrado") {
		t.Fatalf("mensagem inesperada: %s", rec.Body.String())
	}
}

func seedDocument(t *testing.T, store *documents.Store, ext string, content []byte) {
	t.Helper()
	if err := os.WriteFile(store.Path(testCPF, ext), content, 0o644); err != nil {
		t.Fatalf("seed falhou: %v", err)
	}
}

func TestGetClientDocument_StreamsFile(t *testing.T) {
	store := &documents.Store{Dir: t.TempDir()}
	content := append(append([]byte{}, docPdfHead...), []byte("pdf body")...)
	seedDocument(t, store, "pdf", content)

	findByID := &mockFindClientByIDUseCase{
		executeFunc: func(ctx context.Context, clientID int) (*output.FindClient, error) {
			return &output.FindClient{Client: &entities.Client{
				ID: 7, CPF: testCPF, DocumentName: "contrato.pdf", DocumentType: "pdf",
			}}, nil
		},
	}
	handler := documentsTestHandler(nil, store, findByID, nil)

	req := mux.SetURLVars(httptest.NewRequest(http.MethodGet, "/admin/clients/7/document", nil),
		map[string]string{"id": "7"})
	rec := httptest.NewRecorder()
	handler.GetClientDocument(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/pdf" {
		t.Fatalf("Content-Type = %q, wants application/pdf", ct)
	}
	if cd := rec.Header().Get("Content-Disposition"); !strings.Contains(cd, "inline") || !strings.Contains(cd, "contrato.pdf") {
		t.Fatalf("Content-Disposition = %q", cd)
	}
	if !bytes.Equal(rec.Body.Bytes(), content) {
		t.Fatal("conteúdo streamado difere")
	}
}

func TestGetClientDocument_NoDocument(t *testing.T) {
	store := &documents.Store{Dir: t.TempDir()}
	findByID := &mockFindClientByIDUseCase{
		executeFunc: func(ctx context.Context, clientID int) (*output.FindClient, error) {
			return &output.FindClient{Client: &entities.Client{ID: 7, CPF: testCPF}}, nil
		},
	}
	handler := documentsTestHandler(nil, store, findByID, nil)

	req := mux.SetURLVars(httptest.NewRequest(http.MethodGet, "/admin/clients/7/document", nil),
		map[string]string{"id": "7"})
	rec := httptest.NewRecorder()
	handler.GetClientDocument(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "failed") || strings.Contains(rec.Body.String(), "error") {
		t.Fatalf("404 deveria ser limpo: %s", rec.Body.String())
	}
}

func TestGetClientDocument_ClientNotFound(t *testing.T) {
	store := &documents.Store{Dir: t.TempDir()}
	findByID := &mockFindClientByIDUseCase{
		executeFunc: func(ctx context.Context, clientID int) (*output.FindClient, error) {
			return nil, context.DeadlineExceeded // qualquer erro vira 404 limpo
		},
	}
	handler := documentsTestHandler(nil, store, findByID, nil)

	req := mux.SetURLVars(httptest.NewRequest(http.MethodGet, "/admin/clients/99/document", nil),
		map[string]string{"id": "99"})
	rec := httptest.NewRecorder()
	handler.GetClientDocument(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestDeleteClientDocument_RemovesFileAndClearsColumns(t *testing.T) {
	store := &documents.Store{Dir: t.TempDir()}
	seedDocument(t, store, "png", []byte{0x89, 'P', 'N', 'G'})

	var gotName, gotType string
	var gotID, gotSize int
	repo := &mockClientRepository{
		updateDocFunc: func(ctx context.Context, id int, name, docType string, size int) error {
			gotID, gotName, gotType, gotSize = id, name, docType, size
			return nil
		},
	}
	findByID := &mockFindClientByIDUseCase{
		executeFunc: func(ctx context.Context, clientID int) (*output.FindClient, error) {
			return &output.FindClient{Client: &entities.Client{ID: 7, CPF: testCPF, DocumentType: "png"}}, nil
		},
	}
	handler := documentsTestHandler(repo, store, findByID, nil)

	req := mux.SetURLVars(httptest.NewRequest(http.MethodDelete, "/admin/clients/7/document", nil),
		map[string]string{"id": "7"})
	rec := httptest.NewRecorder()
	handler.DeleteClientDocument(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if _, err := os.Stat(store.Path(testCPF, "png")); !os.IsNotExist(err) {
		t.Fatal("arquivo deveria ter sido apagado")
	}
	if gotID != 7 || gotName != "" || gotType != "" || gotSize != 0 {
		t.Fatalf("colunas não foram limpas: (%d, %q, %q, %d)", gotID, gotName, gotType, gotSize)
	}
}

func TestDeleteClient_RemovesDocumentFile(t *testing.T) {
	store := &documents.Store{Dir: t.TempDir()}
	seedDocument(t, store, "pdf", docPdfHead)

	findByID := &mockFindClientByIDUseCase{
		executeFunc: func(ctx context.Context, clientID int) (*output.FindClient, error) {
			return &output.FindClient{Client: &entities.Client{ID: 7, CPF: testCPF, DocumentType: "pdf"}}, nil
		},
	}
	deleteUC := &mockDeleteClientUseCase{
		executeFunc: func(ctx context.Context, clientID int) (*output.DeleteClient, error) {
			return &output.DeleteClient{ClientID: clientID, ClientName: "Maria"}, nil
		},
	}
	handler := documentsTestHandler(nil, store, findByID, deleteUC)

	req := mux.SetURLVars(httptest.NewRequest(http.MethodDelete, "/admin/clients/7", nil),
		map[string]string{"id": "7"})
	rec := httptest.NewRecorder()
	handler.DeleteClient(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if _, err := os.Stat(store.Path(testCPF, "pdf")); !os.IsNotExist(err) {
		t.Fatal("documento deveria ter sido apagado junto com o cliente")
	}
}

var _ = io.Discard
