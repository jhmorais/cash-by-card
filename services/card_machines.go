package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/cardMachine"
	"github.com/jhmorais/cash-by-card/utils"
)

func (h *Handler) ListCardMachines(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	response, err := h.ListCardMachineUseCase.Execute(ctx)
	if err != nil {
		utils.WriteErrModel(w, http.StatusNotFound,
			utils.NewErrorResponse(fmt.Sprintf("failed to get cardMarchines, error: '%s'", err.Error())))
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal cardMarchines response"))
		return
	}
	// for _, cardMachine := range response.CardMachines{
	// 	var brands []string
	// 	errBrandUnmarshal := json.Unmarshal([]byte(cardMachine.Brand), brands)

	// }

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func (h *Handler) GetCardMachineByID(w http.ResponseWriter, r *http.Request) {
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

	response, err := h.FindCardMachineByIDUseCase.Execute(ctx, idInt)
	if err != nil {
		utils.WriteErrModel(w, http.StatusNotFound,
			utils.NewErrorResponse(fmt.Sprintf("failed to find cardMarchine, error: '%s'", err.Error())))
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal cardMarchine response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func (h *Handler) UpdateCardMachine(w http.ResponseWriter, r *http.Request) {
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

	cardMachine := input.UpdateCardMachine{}
	err = json.Unmarshal(body, &cardMachine)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("failed to parse request body"))
		return
	}

	cardMachine.ID, err = strconv.Atoi(id)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("failed to parse param id to int"))
		return
	}

	response, err := h.UpdateCardMachineUseCase.Execute(ctx, &cardMachine)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest,
			utils.NewErrorResponse(fmt.Sprintf("failed to update cardMachine, error:'%s'", err.Error())))
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal cardMachine response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func (h *Handler) DeleteCardMachine(w http.ResponseWriter, r *http.Request) {
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

	response, err := h.DeleteCardMachineUseCase.Execute(ctx, idInt)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest,
			utils.NewErrorResponse(fmt.Sprintf("failed to delete cardMachine, error: '%s'", err.Error())))
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal cardMachine response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func (h *Handler) CreateCardMachine(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading request body"))
		return
	}

	cardMachine := input.CreateCardMachine{}
	err = json.Unmarshal(body, &cardMachine)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("failed to parse request body"))
		return
	}

	response, err := h.CreateCardMachineUseCase.Execute(ctx, &cardMachine)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse(fmt.Sprintf("failed to create card, error: '%s'", err.Error())))
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal card response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}
