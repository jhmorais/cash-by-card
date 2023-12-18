package partner

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/partner"
)

type findPartnerByNameBrandUseCase struct {
	partnerRepository repositories.PartnerRepository
}

func NewFindPartnerByNameUseCase(partnerRepository repositories.PartnerRepository) contracts.FindPartnerByNameUseCase {

	return &findPartnerByNameBrandUseCase{
		partnerRepository: partnerRepository,
	}
}

func (c *findPartnerByNameBrandUseCase) Execute(ctx context.Context, name string) (*output.ListPartner, error) {

	partnerEntity, err := c.partnerRepository.FindPartnerByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("erro to find partner with brand: '%s' at database: '%v'", name, err)
	}

	if len(partnerEntity) == 0 {
		return nil, fmt.Errorf("partners not found")
	}

	output := &output.ListPartner{
		Partners: partnerEntity,
	}

	return output, nil
}
