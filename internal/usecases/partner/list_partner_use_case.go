package partner

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/partner"
)

type listPartnerUseCase struct {
	partnerRepository repositories.PartnerRepository
}

func NewListPartnerUseCase(partnerRepository repositories.PartnerRepository) contracts.ListPartnerUseCase {

	return &listPartnerUseCase{
		partnerRepository: partnerRepository,
	}
}

func (l *listPartnerUseCase) Execute(ctx context.Context) (*output.ListPartner, error) {
	var err error
	output := &output.ListPartner{Partners: []*entities.Partner{}}

	output.Partners, err = l.partnerRepository.ListPartner(ctx)
	if err != nil {
		return nil, fmt.Errorf("error when list partners on database: %v", err)
	}

	return output, nil
}
