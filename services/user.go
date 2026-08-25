package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	"github.com/jhmorais/cash-by-card/utils"
)

// writeUserError mapeia erros de permissao dos use cases para 403; demais para 400.
func writeUserError(w http.ResponseWriter, err error) {
	msg := err.Error()
	if strings.Contains(msg, "permiss") || strings.Contains(msg, "só pode") || strings.Contains(msg, "próprio") || strings.Contains(msg, "24h") {
		utils.WriteErrModel(w, http.StatusForbidden, utils.NewErrorResponse(msg))
		return
	}
	utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse(msg))
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	requesterEmail := utils.EmailFromContext(r.Context())
	response, err := h.ListUserUseCase.Execute(context.Background(), requesterEmail)
	if err != nil {
		writeUserError(w, err)
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

// func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
// 	ctx := context.Background()
// 	id, err := utils.RetrieveParam(r, "id")
// 	if err != nil {
// 		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading id"))
// 		return
// 	}

// 	idInt, err := strconv.Atoi(id)
// 	if err != nil {
// 		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error cast id to int"))
// 		return
// 	}

// 	response, err := h.FindClientByIDUseCase.Execute(ctx, idInt)
// 	if err != nil {
// 		utils.WriteErrModel(w, http.StatusNotFound,
// 			utils.NewErrorResponse(fmt.Sprintf("failed to find user, error: '%s'", err.Error())))
// 		return
// 	}

// 	jsonResponse, err := json.Marshal(response)
// 	if err != nil {
// 		utils.WriteErrModel(w, http.StatusInternalServerError,
// 			utils.NewErrorResponse("Failed to marshal user response"))
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprint(w, string(jsonResponse))
// }

// func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
// 	ctx := context.Background()
// 	name, err := utils.RetrieveParam(r, "name")
// 	if err != nil {
// 		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading name"))
// 		return
// 	}

// 	response, err := h.FindClientByNameUseCase.Execute(ctx, name)
// 	if err != nil {
// 		utils.WriteErrModel(w, http.StatusNotFound,
// 			utils.NewErrorResponse(fmt.Sprintf("failed to find user, error: '%s'", err.Error())))
// 		return
// 	}

// 	jsonResponse, err := json.Marshal(response)
// 	if err != nil {
// 		utils.WriteErrModel(w, http.StatusInternalServerError,
// 			utils.NewErrorResponse("Failed to marshal user response"))
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprint(w, string(jsonResponse))
// }

// func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
// 	ctx := context.Background()
// 	id, err := utils.RetrieveParam(r, "id")
// 	if err != nil {
// 		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading id"))
// 		return
// 	}

// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading request body"))
// 		return
// 	}

// 	user := input.UpdateClient{}
// 	err = json.Unmarshal(body, &user)
// 	if err != nil {
// 		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("failed to parse request body"))
// 		return
// 	}

// 	user.ID, err = strconv.Atoi(id)
// 	if err != nil {
// 		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("failed to parse param id to int"))
// 		return
// 	}

// 	response, err := h.UpdateClientUseCase.Execute(ctx, &user)
// 	if err != nil {
// 		utils.WriteErrModel(w, http.StatusBadRequest,
// 			utils.NewErrorResponse(fmt.Sprintf("failed to update user, error:'%s'", err.Error())))
// 		return
// 	}

// 	jsonResponse, err := json.Marshal(response)
// 	if err != nil {
// 		utils.WriteErrModel(w, http.StatusInternalServerError,
// 			utils.NewErrorResponse("Failed to marshal user response"))
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprint(w, string(jsonResponse))
// }

// func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
// 	ctx := context.Background()
// 	id, err := utils.RetrieveParam(r, "id")
// 	if err != nil {
// 		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading id"))
// 		return
// 	}

// 	idInt, err := strconv.Atoi(id)
// 	if err != nil {
// 		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error cast id to int"))
// 		return
// 	}

// 	response, err := h.DeleteClientUseCase.Execute(ctx, idInt)
// 	if err != nil {
// 		utils.WriteErrModel(w, http.StatusBadRequest,
// 			utils.NewErrorResponse(fmt.Sprintf("failed to delete user, error: '%s'", err.Error())))
// 		return
// 	}

// 	jsonResponse, err := json.Marshal(response)
// 	if err != nil {
// 		utils.WriteErrModel(w, http.StatusInternalServerError,
// 			utils.NewErrorResponse("Failed to marshal user response"))
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprint(w, string(jsonResponse))
// }

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
		writeUserError(w, err)
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
		writeUserError(w, err)
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
		writeUserError(w, err)
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
		writeUserError(w, err)
		return
	}
	jsonResponse, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}
