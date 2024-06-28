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

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	response, err := h.ListClientUseCase.Execute(ctx)
	if err != nil {
		utils.WriteErrModel(w, http.StatusNotFound,
			utils.NewErrorResponse(fmt.Sprintf("failed to get users, error: '%s'", err.Error())))
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal user response"))
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
	ctx := context.Background()
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading request body"))
		return
	}

	user := input.CreateUser{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("failed to parse request body"))
		return
	}

	response, err := h.CreateUserUseCase.Execute(ctx, &user)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse(fmt.Sprintf("failed to create user, error: '%s'", err.Error())))
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal user response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}
