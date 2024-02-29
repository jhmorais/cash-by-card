package loan

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/loan"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/loan"
)

type updateLoanPaymentStatusUseCase struct {
	loanRepository repositories.LoanRepository
}

func NewUpdateLoanPaymentStatusUseCase(loanRepository repositories.LoanRepository) contracts.UpdateLoanPaymentStatusUseCase {

	return &updateLoanPaymentStatusUseCase{
		loanRepository: loanRepository,
	}
}

func (c *updateLoanPaymentStatusUseCase) Execute(ctx context.Context, updateLoanPaymentStatus *input.UpdateLoanPaymentStatus) error {

	if !(updateLoanPaymentStatus.PaymentStatus == "pending" || updateLoanPaymentStatus.PaymentStatus == "paid") {
		return fmt.Errorf("cannot update a loan with invalid payment status")
	}

	err := c.loanRepository.UpdateLoanPaymentStatus(ctx, updateLoanPaymentStatus.ID, updateLoanPaymentStatus.PaymentStatus)
	if err != nil {
		return fmt.Errorf("cannot update PaymentStatus: %v", err)
	}

	return nil

}
