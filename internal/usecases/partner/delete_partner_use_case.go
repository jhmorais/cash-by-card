package partner

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/partner"
)

type deletePartnerUseCase struct {
	partnerRepository repositories.PartnerRepository
}

func NewDeletePartnerUseCase(partnerRepository repositories.PartnerRepository) contracts.DeletePartnerUseCase {

	return &deletePartnerUseCase{
		partnerRepository: partnerRepository,
	}
}

func (c *deletePartnerUseCase) Execute(ctx context.Context, partnerID int) (*output.DeletePartner, error) {

	partnerEntity, err := c.partnerRepository.FindPartnerByID(ctx, partnerID)
	if err != nil {
		return nil, fmt.Errorf("failed to find partner '%d' at database: '%v'", partnerID, err)
	}

	if partnerEntity == nil || partnerEntity.ID == 0 {
		return nil, fmt.Errorf("partnerID not found")
	}

	err = c.partnerRepository.DeletePartner(ctx, partnerEntity)
	if err != nil {
		return nil, fmt.Errorf("falha para deletar o parceiro de id '%d'", partnerEntity.ID)
	}

	output := &output.DeletePartner{
		PartnerID:   partnerEntity.ID,
		PartnerName: partnerEntity.Name,
	}

	return output, nil
}
