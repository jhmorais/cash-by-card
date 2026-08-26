package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/jhmorais/cash-by-card/internal/documents"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/client"
	"github.com/jhmorais/cash-by-card/utils"
)

// maxDocumentSize é o tamanho máximo do documento do cliente (5MB).
const maxDocumentSize = 5 * 1024 * 1024

// contentTypeByExt devolve o Content-Type do documento por extensão.
func contentTypeByExt(ext string) string {
	switch ext {
	case "pdf":
		return "application/pdf"
	case "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	default:
		return "application/octet-stream"
	}
}

func (h *Handler) ListClients(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	response, err := h.ListClientUseCase.Execute(ctx)
	if err != nil {
		utils.WriteErrModel(w, http.StatusNotFound,
			utils.NewErrorResponse(fmt.Sprintf("failed to get clients, error: '%s'", err.Error())))
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal client response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func (h *Handler) GetClientByID(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id, err := utils.RetrieveParam(r, "id")
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading id"))
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error cast id to int"))
		return
	}

	response, err := h.FindClientByIDUseCase.Execute(ctx, idInt)
	if err != nil {
		utils.WriteErrModel(w, http.StatusNotFound,
			utils.NewErrorResponse(fmt.Sprintf("failed to find client, error: '%s'", err.Error())))
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal client response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func (h *Handler) GetClient(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	name, err := utils.RetrieveParam(r, "name")
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading name"))
		return
	}

	response, err := h.FindClientByNameUseCase.Execute(ctx, name)
	if err != nil {
		utils.WriteErrModel(w, http.StatusNotFound,
			utils.NewErrorResponse(fmt.Sprintf("failed to find client, error: '%s'", err.Error())))
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal client response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func (h *Handler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id, err := utils.RetrieveParam(r, "id")
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading id"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading request body"))
		return
	}

	client := input.UpdateClient{}
	err = json.Unmarshal(body, &client)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("failed to parse request body"))
		return
	}

	client.ID, err = strconv.Atoi(id)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("failed to parse param id to int"))
		return
	}

	response, err := h.UpdateClientUseCase.Execute(ctx, &client)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest,
			utils.NewErrorResponse(err.Error()))
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal client response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func (h *Handler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id, err := utils.RetrieveParam(r, "id")
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading id"))
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error cast id to int"))
		return
	}

	// apaga o documento do cliente junto (arquivo ausente é tolerado)
	if response, err := h.FindClientByIDUseCase.Execute(ctx, idInt); err == nil && response != nil && response.Client != nil {
		_ = h.Documents.Delete(response.Client.CPF)
	}

	response, err := h.DeleteClientUseCase.Execute(ctx, idInt)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest,
			utils.NewErrorResponse((err.Error())))
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal client response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func (h *Handler) CreateClient(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading request body"))
		return
	}

	client := input.CreateClient{}
	err = json.Unmarshal(body, &client)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("Campos não preenchidos"))
		return
	}

	response, err := h.CreateClientUseCase.Execute(ctx, &client)
	if err != nil {
		var errorStatus *utils.RequestError
		if errors.As(err, &errorStatus) {
			utils.WriteErrModel(w, errorStatus.StatusCode,
				utils.NewErrorResponse(errorStatus.Error()))
			return
		}
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse(err.Error()))
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal client response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

// CreateClientDocuments recebe o documento único do cliente (PDF, JPEG ou
// PNG, até 5MB), valida por magic bytes, salva em DOCS_DIR/{cpf}.{ext}
// substituindo o anterior e grava os metadados no banco.
func (h *Handler) CreateClientDocuments(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	cpf, err := utils.RetrieveParam(r, "cpf")
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading cpf"))
		return
	}

	clients, err := h.ClientRepository.FindClientByCPF(ctx, cpf)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("failed to find client"))
		return
	}
	if len(clients) == 0 || clients[0].ID == 0 {
		utils.WriteErrModel(w, http.StatusNotFound, utils.NewErrorResponse("cliente não encontrado"))
		return
	}
	client := clients[0]

	if err := r.ParseMultipartForm(maxDocumentSize); err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error parsing multipart form"))
		return
	}

	var files []*multipart.FileHeader
	for _, headers := range r.MultipartForm.File {
		files = append(files, headers...)
	}
	if len(files) != 1 {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("envie apenas um arquivo"))
		return
	}

	fileHeader := files[0]
	if fileHeader.Size > maxDocumentSize {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("arquivo maior que 5MB"))
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError, utils.NewErrorResponse("error opening file"))
		return
	}
	defer src.Close()

	// lê o começo do arquivo para validar os magic bytes
	head := make([]byte, 512)
	n, err := io.ReadFull(src, head)
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		utils.WriteErrModel(w, http.StatusInternalServerError, utils.NewErrorResponse("error reading file"))
		return
	}
	head = head[:n]

	ext, err := h.Documents.Save(cpf, head, src, fileHeader.Size)
	if err != nil {
		if errors.Is(err, documents.ErrUnsupportedFormat) {
			utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse(err.Error()))
			return
		}
		utils.WriteErrModel(w, http.StatusInternalServerError, utils.NewErrorResponse("error saving file"))
		return
	}

	if err := h.ClientRepository.UpdateClientDocument(ctx, client.ID, fileHeader.Filename, ext, int(fileHeader.Size)); err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("failed to update client document"))
		return
	}

	metadata := map[string]interface{}{
		"documentName": fileHeader.Filename,
		"documentType": ext,
		"documentSize": fileHeader.Size,
	}
	jsonResponse, err := json.Marshal(metadata)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal client documents response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

// GetClientDocument streama o documento do cliente com o Content-Type
// correto e Content-Disposition inline (visualização no navegador).
func (h *Handler) GetClientDocument(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id, err := utils.RetrieveParam(r, "id")
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading id"))
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error cast id to int"))
		return
	}

	response, err := h.FindClientByIDUseCase.Execute(ctx, idInt)
	if err != nil || response == nil || response.Client == nil {
		utils.WriteErrModel(w, http.StatusNotFound, utils.NewErrorResponse("cliente não encontrado"))
		return
	}
	client := response.Client

	if client.DocumentType == "" {
		utils.WriteErrModel(w, http.StatusNotFound, utils.NewErrorResponse("documento não encontrado"))
		return
	}

	file, err := os.Open(h.Documents.Path(client.CPF, client.DocumentType))
	if err != nil {
		if os.IsNotExist(err) {
			utils.WriteErrModel(w, http.StatusNotFound, utils.NewErrorResponse("documento não encontrado"))
			return
		}
		utils.WriteErrModel(w, http.StatusInternalServerError, utils.NewErrorResponse("error opening document"))
		return
	}
	defer file.Close()

	filename := client.DocumentName
	if filename == "" {
		filename = client.CPF + "." + client.DocumentType
	}

	w.Header().Set("Content-Type", contentTypeByExt(client.DocumentType))
	w.Header().Set("Content-Disposition", mime.FormatMediaType("inline", map[string]string{"filename": filename}))
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, file); err != nil {
		// headers já foram enviados; nada mais a fazer
		return
	}
}

// DeleteClientDocument apaga o arquivo do documento e limpa as colunas.
func (h *Handler) DeleteClientDocument(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id, err := utils.RetrieveParam(r, "id")
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading id"))
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error cast id to int"))
		return
	}

	response, err := h.FindClientByIDUseCase.Execute(ctx, idInt)
	if err != nil || response == nil || response.Client == nil {
		utils.WriteErrModel(w, http.StatusNotFound, utils.NewErrorResponse("cliente não encontrado"))
		return
	}

	if err := h.Documents.Delete(response.Client.CPF); err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError, utils.NewErrorResponse("error deleting document"))
		return
	}

	if err := h.ClientRepository.UpdateClientDocument(ctx, idInt, "", "", 0); err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("failed to update client document"))
		return
	}

	jsonResponse, err := json.Marshal(map[string]string{"message": "documento removido"})
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal client documents response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}
