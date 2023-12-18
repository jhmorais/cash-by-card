package partner

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/partner"
)

type findPartnerByIDUseCase struct {
	partnerRepository repositories.PartnerRepository
}

func NewFindPartnerByIDUseCase(partnerRepository repositories.PartnerRepository) contracts.FindPartnerByIDUseCase {

	return &findPartnerByIDUseCase{
		partnerRepository: partnerRepository,
	}
}

func (c *findPartnerByIDUseCase) Execute(ctx context.Context, partnerID int) (*output.FindPartner, error) {

	partnerEntity, err := c.partnerRepository.FindPartnerByID(ctx, partnerID)
	if err != nil {
		return nil, fmt.Errorf("erro to find partner '%d' at database: '%v'", partnerID, err)
	}

	if partnerEntity == nil || partnerEntity.ID == 0 {
		return nil, fmt.Errorf("partnerID not found")
	}

	output := &output.FindPartner{
		Partner: partnerEntity,
	}

	return output, nil
}
