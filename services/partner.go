package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jhmorais/cash-by-card/utils"
)

func (h *Handler) ListPartners(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	response, err := h.ListPartnerUseCase.Execute(ctx)
	if err != nil {
		utils.WriteErrModel(w, http.StatusNotFound,
			utils.NewErrorResponse(fmt.Sprintf("failed to get partners, error: '%s'", err.Error())))
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("Failed to marshal partner response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}
