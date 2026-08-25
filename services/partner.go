package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/client"
	"github.com/jhmorais/cash-by-card/utils"
)

func (h *Handler) PartnerListClients(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	requesterEmail := utils.EmailFromContext(r.Context())
	response, err := h.PartnerClientsUseCase.ListClients(ctx, requesterEmail)
	if err != nil {
		writeUserError(w, err)
		return
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal partner clients response"))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func (h *Handler) PartnerCreateClient(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading request body"))
		return
	}
	client := input.CreateClient{}
	if err := json.Unmarshal(body, &client); err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("failed to parse request body"))
		return
	}
	requesterEmail := utils.EmailFromContext(r.Context())
	response, err := h.PartnerClientsUseCase.CreateClient(ctx, requesterEmail, &client)
	if err != nil {
		writeUserError(w, err)
		return
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal partner client response"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(jsonResponse))
}

func (h *Handler) PartnerUpdateClient(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
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
	client := input.UpdateClient{ID: id}
	if err := json.Unmarshal(body, &client); err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("failed to parse request body"))
		return
	}
	client.ID = id
	requesterEmail := utils.EmailFromContext(r.Context())
	response, err := h.PartnerClientsUseCase.UpdateClient(ctx, requesterEmail, &client)
	if err != nil {
		writeUserError(w, err)
		return
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal partner client response"))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func (h *Handler) PartnerCardMachines(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	response, err := h.ListCardMachineUseCase.Execute(ctx)
	if err != nil {
		utils.WriteErrModel(w, http.StatusNotFound,
			utils.NewErrorResponse(fmt.Sprintf("failed to get cardMachines, error: '%s'", err.Error())))
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal cardMachines response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func (h *Handler) PartnerReport(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	requesterEmail := utils.EmailFromContext(r.Context())
	response, err := h.PartnerReportUseCase.Execute(ctx, requesterEmail)
	if err != nil {
		writeUserError(w, err)
		return
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal partner report response"))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}
