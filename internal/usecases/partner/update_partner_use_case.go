package partner

import (
	"context"
	"fmt"
	"time"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/partner"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/partner"
)

type updatePartnerUseCase struct {
	partnerRepository repositories.PartnerRepository
}

func NewUpdatePartnerUseCase(partnerRepository repositories.PartnerRepository) contracts.UpdatePartnerUseCase {

	return &updatePartnerUseCase{
		partnerRepository: partnerRepository,
	}
}

func (c *updatePartnerUseCase) Execute(ctx context.Context, updatePartner *input.UpdatePartner) (*output.CreatePartner, error) {
	if updatePartner.Name == "" {
		return nil, fmt.Errorf("failed name partner is empty")
	}

	if updatePartner.PixKey == "" {
		return nil, fmt.Errorf("failed pix key partner is empty")
	}

	partner, err := c.partnerRepository.FindPartnerByCPF(ctx, updatePartner.CPF)
	if err != nil {
		return nil, fmt.Errorf("failed to get partner")
	}

	if len(partner) > 0 && partner[0].ID != updatePartner.ID {
		return nil, fmt.Errorf("failed, already exists partner with the same cpf")
	}

	if len(updatePartner.PixKey) > 250 {
		updatePartner.PixKey = updatePartner.PixKey[:250]
	}

	partnerEntity := &entities.Partner{
		ID:        updatePartner.ID,
		Name:      updatePartner.Name,
		PixKey:    updatePartner.PixKey,
		CPF:       updatePartner.CPF,
		Telefone:  updatePartner.Phone,
		Endereco:  updatePartner.Address,
		CreatedAt: time.Now(),
	}

	errUpdate := c.partnerRepository.UpdatePartner(ctx, partnerEntity)
	if errUpdate != nil {
		return nil, fmt.Errorf("cannot update partner at database: %v", errUpdate)
	}

	createPartnerOutput := &output.CreatePartner{
		PartnerID: partnerEntity.ID,
	}

	return createPartnerOutput, nil
}
