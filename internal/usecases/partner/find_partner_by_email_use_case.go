package partner

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/partner"
)

type findPartnerByEmailUseCase struct {
	partnerRepository repositories.PartnerRepository
}

func NewFindPartnerByEmailUseCase(partnerRepository repositories.PartnerRepository) contracts.FindPartnerByEmailUseCase {

	return &findPartnerByEmailUseCase{
		partnerRepository: partnerRepository,
	}
}

func (c *findPartnerByEmailUseCase) Execute(ctx context.Context, email string) (*output.ListPartner, error) {

	partnerEntity, err := c.partnerRepository.FindPartnerByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("erro to find partner with email: '%s' at database: '%v'", email, err)
	}

	if len(partnerEntity) == 0 {
		return nil, fmt.Errorf("partners not found")
	}

	output := &output.ListPartner{
		Partners: partnerEntity,
	}

	return output, nil
}
