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

type createPartnerUseCase struct {
	partnerRepository repositories.PartnerRepository
}

func NewCreatePartnerUseCase(partnerRepository repositories.PartnerRepository) contracts.CreatePartnerUseCase {

	return &createPartnerUseCase{
		partnerRepository: partnerRepository,
	}
}

func (c *createPartnerUseCase) Execute(ctx context.Context, createPartner *input.CreatePartner) (*output.CreatePartner, error) {

	if len(createPartner.Name) > 250 {
		createPartner.Name = createPartner.Name[:250]
	}

	if createPartner.Phone == "" {
		return nil, fmt.Errorf("telefone não pode ser vazio")
	}

	if createPartner.CPF == "" {
		return nil, fmt.Errorf("cpf não pode ser vazio")
	}

	partnerEntity := &entities.Partner{
		Name:      createPartner.Name,
		PixKey:    createPartner.PixKey,
		CPF:       createPartner.CPF,
		Phone:     createPartner.Phone,
		Address:   createPartner.Address,
		Email:     createPartner.Email,
		PixType:   createPartner.PixType,
		CreatedAt: time.Now(),
	}

	partner, err := c.partnerRepository.FindPartnerByCPF(ctx, createPartner.CPF)
	if err != nil {
		return nil, fmt.Errorf("failed to get partner: %v", err)
	}

	if len(partner) > 0 {
		return nil, fmt.Errorf("falha, já existe um parceiro com esse cpf")
	}

	partner, err = c.partnerRepository.FindPartnerByEmail(ctx, createPartner.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get partner: %v", err)
	}

	if len(partner) > 0 {
		return nil, fmt.Errorf("falha, já existe um parceiro com esse email")
	}

	err = c.partnerRepository.CreatePartner(ctx, partnerEntity)
	if err != nil {
		return nil, fmt.Errorf("não foi possivel salvar o parceiro: %v", err)
	}

	createPartnerOutput := &output.CreatePartner{
		PartnerID: partnerEntity.ID,
	}

	return createPartnerOutput, nil
}
