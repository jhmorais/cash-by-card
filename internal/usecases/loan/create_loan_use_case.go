package loan

import (
	"context"
	"fmt"
	"time"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	inputCard "github.com/jhmorais/cash-by-card/internal/ports/input/card"
	inputLoan "github.com/jhmorais/cash-by-card/internal/ports/input/loan"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/loan"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/loan"
)

type createLoanUseCase struct {
	loanRepository repositories.LoanRepository
	cardUsecase    contracts.CreateCardUseCase
}

func NewCreateLoanUseCase(loanRepository repositories.LoanRepository, cardUsecase contracts.CreateCardUseCase) contracts.CreateLoanUseCase {

	return &createLoanUseCase{
		loanRepository: loanRepository,
		cardUsecase:    cardUsecase,
	}
}

func (c *createLoanUseCase) Execute(ctx context.Context, createLoan *inputLoan.CreateLoan) (*output.CreateLoan, error) {

	if createLoan.ClientID < 0 {
		return nil, fmt.Errorf("cannot create a loan without ClientId")
	}

	if createLoan.PaymentStatus == "" {
		createLoan.PaymentStatus = "pending"
	}

	if !(createLoan.PaymentStatus == "pending" || createLoan.PaymentStatus == "paid") {
		return nil, fmt.Errorf("cannot create a loan with invalid payment status")
	}

	if createLoan.AskValue == 0 {
		return nil, fmt.Errorf("cannot create a loan without AskValue")
	}
	if createLoan.Amount == 0 {
		return nil, fmt.Errorf("cannot create a loan without Amount")
	}

	if createLoan.OperationPercent < 0 {
		return nil, fmt.Errorf("cannot create a loan without valid OperationPercent")
	}

	if createLoan.NumberCards <= 0 {
		return nil, fmt.Errorf("cannot create loan without valid number of cards")
	}

	if len(createLoan.Cards) == 0 {
		return nil, fmt.Errorf("cannot create loan without card")
	}

	if createLoan.PartnerID < 0 {
		return nil, fmt.Errorf("cannot create a loan without PartnerId")
	}

	if createLoan.GrossProfit == 0 {
		return nil, fmt.Errorf("cannot create a loan without valid GrossProfit")
	}
	if createLoan.Profit == 0 {
		return nil, fmt.Errorf("cannot create a loan without valid Profit")
	}

	loanEntity := &entities.Loan{
		ClientID:         createLoan.ClientID,
		AskValue:         createLoan.AskValue,
		Amount:           createLoan.Amount,
		OperationPercent: createLoan.OperationPercent,
		NumberCards:      createLoan.NumberCards,
		PartnerID:        createLoan.PartnerID,
		GrossProfit:      createLoan.GrossProfit,
		PartnerPercent:   createLoan.PartnerPercent,
		PartnerAmount:    createLoan.PartnerAmount,
		Profit:           createLoan.Profit,
		PaymentStatus:    createLoan.PaymentStatus,
		CreatedAt:        time.Now(),
	}

	err := c.loanRepository.CreateLoan(ctx, loanEntity)
	if err != nil {
		return nil, fmt.Errorf("cannot save loan at database: %v", err)
	}

	cardsInput := []inputCard.CreateCard{}
	for _, card := range createLoan.Cards {
		cardsInput = append(cardsInput, inputCard.CreateCard{
			PaymentType:   card.PaymentType,
			Value:         card.Value,
			Brand:         card.Brand,
			Installments:  card.Installments,
			LoanID:        loanEntity.ID,
			CardMachineID: card.CardMachineID,
		})
	}

	_, err = c.cardUsecase.Execute(ctx, &cardsInput)
	if err != nil {
		return nil, fmt.Errorf("cannot save cards at database: %v", err)
	}

	createLoanOutput := &output.CreateLoan{
		LoanID: loanEntity.ID,
	}

	return createLoanOutput, nil
}
