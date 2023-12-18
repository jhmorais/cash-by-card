package partner

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/partner"
)

type findPartnerUseCase struct {
	partnerRepository repositories.PartnerRepository
}

func NewFindPartnerUseCase(partnerRepository repositories.PartnerRepository) contracts.FindPartnerUseCase {

	return &findPartnerUseCase{
		partnerRepository: partnerRepository,
	}
}

func (c *findPartnerUseCase) Execute(ctx context.Context, partnerID int, name string) (*output.FindPartner, error) {
	if partnerID == 0 || name == "" {
		return nil, fmt.Errorf("failed brand or name are empty")
	}

	partnerEntity, err := c.partnerRepository.FindPartnerByPartnerID(ctx, partnerID, name)
	if err != nil {
		return nil, fmt.Errorf("erro to find partner with brand: '%d' at database: '%v'", partnerID, err)
	}

	if len(partnerEntity) == 0 || partnerEntity[0].ID == 0 {
		return nil, fmt.Errorf("partner not found")
	}

	output := &output.FindPartner{
		Partner: partnerEntity[0],
	}

	return output, nil
}
