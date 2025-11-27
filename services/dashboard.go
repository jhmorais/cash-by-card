package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jhmorais/cash-by-card/utils"
)

func (h *Handler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	monthParam := r.URL.Query().Get("month")
	yearParam := r.URL.Query().Get("year")

	if monthParam == "" || yearParam == "" {
		utils.WriteErrModel(w, http.StatusBadRequest,
			utils.NewErrorResponse("month and year query params are required"))
		return
	}

	month, err := strconv.Atoi(monthParam)
	if err != nil || month < 1 || month > 12 {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("invalid month parameter"))
		return
	}

	year, err := strconv.Atoi(yearParam)
	if err != nil || year <= 0 {
		utils.WriteErrModel(w, http.StatusBadRequest,
			utils.NewErrorResponse("invalid year"))
		return
	}

	response, err := h.DashboardUseCase.Execute(ctx, month, year)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse(fmt.Sprintf("failed to get dashboard data: %s", err.Error())))
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse("failed to marshal dashboard response"))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}
