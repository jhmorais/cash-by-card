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

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading request body"))
		return
	}
	forgot := input.ForgotPassword{}
	if err := json.Unmarshal(body, &forgot); err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("failed to parse request body"))
		return
	}
	if err := h.ForgotPasswordUseCase.Execute(ctx, &forgot); err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError,
			utils.NewErrorResponse(fmt.Sprintf("failed to process forgot password, error: '%s'", err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"message":"Se o email estiver cadastrado, você receberá as instruções por email","type":"SUCCESS"}`)
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading request body"))
		return
	}
	reset := input.ResetPassword{}
	if err := json.Unmarshal(body, &reset); err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("failed to parse request body"))
		return
	}
	if err := h.ResetPasswordUseCase.Execute(ctx, &reset); err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest,
			utils.NewErrorResponse(fmt.Sprintf("failed to reset password, error: '%s'", err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"message":"senha definida com sucesso","type":"SUCCESS"}`)
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("error reading request body"))
		return
	}
	change := input.ChangePassword{}
	if err := json.Unmarshal(body, &change); err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest, utils.NewErrorResponse("failed to parse request body"))
		return
	}
	email := utils.EmailFromContext(r.Context())
	if err := h.ChangePasswordUseCase.Execute(ctx, email, &change); err != nil {
		utils.WriteErrModel(w, http.StatusBadRequest,
			utils.NewErrorResponse(fmt.Sprintf("failed to change password, error: '%s'", err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"message":"senha alterada com sucesso","type":"SUCCESS"}`)
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	email := utils.EmailFromContext(r.Context())
	user, err := h.GetUserUseCase.Execute(context.Background(), email)
	if err != nil {
		utils.WriteErrModel(w, http.StatusNotFound,
			utils.NewErrorResponse(fmt.Sprintf("failed to get user, error: '%s'", err.Error())))
		return
	}
	jsonResponse, err := json.Marshal(user)
	if err != nil {
		utils.WriteErrModel(w, http.StatusInternalServerError, utils.NewErrorResponse("Failed to marshal user response"))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}
