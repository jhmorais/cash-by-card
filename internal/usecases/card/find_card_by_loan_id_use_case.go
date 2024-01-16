package card

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/card"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/card"
)

type findCardByLoanIDUseCase struct {
	cardRepository repositories.CardRepository
}

func NewFindCardByLoanIDUseCase(cardRepository repositories.CardRepository) contracts.FindCardByLoanIDUseCase {

	return &findCardByLoanIDUseCase{
		cardRepository: cardRepository,
	}
}

func (c *findCardByLoanIDUseCase) Execute(ctx context.Context, LoanID int) (*output.ListCard, error) {

	cardsEntity, err := c.cardRepository.FindCardByLoanID(ctx, LoanID)
	if err != nil {
		return nil, fmt.Errorf("erro to find card with LoanId: '%v' at database: '%v'", LoanID, err)
	}

	if len(cardsEntity) == 0 {
		return nil, fmt.Errorf("cards not found")
	}

	output := &output.ListCard{
		Cards: cardsEntity,
	}

	return output, nil
}
