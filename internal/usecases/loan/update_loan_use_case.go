package loan

import (
	"context"
	"fmt"
	"time"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/loan"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/loan"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/loan"
)

type updateLoanUseCase struct {
	loanRepository repositories.LoanRepository
}

func NewUpdateLoanUseCase(loanRepository repositories.LoanRepository) contracts.UpdateLoanUseCase {

	return &updateLoanUseCase{
		loanRepository: loanRepository,
	}
}

func (c *updateLoanUseCase) Execute(ctx context.Context, updateLoan *input.UpdateLoan) (*output.CreateLoan, error) {

	if updateLoan.ClientID < 0 {
		return nil, fmt.Errorf("cannot update a loan without ClientId")
	}

	if !(updateLoan.PaymentStatus == "pending" || updateLoan.PaymentStatus == "paid") {
		return nil, fmt.Errorf("cannot update a loan with invalid payment status")
	}

	if updateLoan.AskValue == 0 {
		return nil, fmt.Errorf("cannot update a loan without AskValue")
	}
	if updateLoan.Amount == 0 {
		return nil, fmt.Errorf("cannot update a loan without Amount")
	}

	if updateLoan.OperationPercent < 0 {
		return nil, fmt.Errorf("cannot update a loan without valid OperationPercent")
	}

	if updateLoan.NumberCards <= 0 {
		return nil, fmt.Errorf("cannot update loan without valid number of cards")
	}

	if len(updateLoan.Cards) == 0 {
		return nil, fmt.Errorf("cannot update loan without card")
	}

	if updateLoan.PartnerID < 0 {
		return nil, fmt.Errorf("cannot update a loan without PartnerId")
	}

	if updateLoan.GrossProfit == 0 {
		return nil, fmt.Errorf("cannot update a loan without valid GrossProfit")
	}
	if updateLoan.Profit == 0 {
		return nil, fmt.Errorf("cannot update a loan without valid Profit")
	}
	if !(updateLoan.Type == 1 || updateLoan.Type == 2) {
		return nil, fmt.Errorf("cannot update a loan with invalid type")
	}

	if updateLoan.ClientAmount == 0 {
		return nil, fmt.Errorf("cannot update a loan without valid client amount")
	}

	cardsEntity := []entities.Card{}
	for _, card := range updateLoan.Cards {
		cardsEntity = append(cardsEntity, entities.Card{
			PaymentType:       card.PaymentType,
			Value:             card.Value,
			Brand:             card.Brand,
			Installments:      card.Installments,
			InstallmentsValue: card.InstallmentsValue,
			LoanID:            card.LoanID,
			CardMachineID:     card.CardMachineID,
		})
	}
	loanEntity := &entities.Loan{
		ID:               updateLoan.ID,
		PaymentStatus:    updateLoan.PaymentStatus,
		AskValue:         updateLoan.AskValue,
		Amount:           updateLoan.Amount,
		Cards:            cardsEntity,
		NumberCards:      updateLoan.NumberCards,
		ClientID:         updateLoan.ClientID,
		PartnerID:        &updateLoan.PartnerID,
		PartnerAmount:    updateLoan.PartnerAmount,
		PartnerPercent:   updateLoan.PartnerPercent,
		Profit:           updateLoan.Profit,
		GrossProfit:      updateLoan.GrossProfit,
		OperationPercent: updateLoan.OperationPercent,
		ClientAmount:     updateLoan.ClientAmount,
		Type:             updateLoan.Type,
		UpdatedAt:        time.Now(),
	}

	errUpdate := c.loanRepository.UpdateLoan(ctx, loanEntity)
	if errUpdate != nil {
		return nil, fmt.Errorf("cannot update loan at database: %v", errUpdate)
	}

	updateLoanOutput := &output.CreateLoan{
		LoanID: loanEntity.ID,
	}

	return updateLoanOutput, nil
}
