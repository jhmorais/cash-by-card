package sample

import (
	"context"
	"time"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	"github.com/jhmorais/cash-by-card/internal/repositories"
	"gorm.io/gorm"
)

func DBSeed(db *gorm.DB) error {
	deviceRepository := repositories.NewClientRepository(db)

	err := createClient(deviceRepository, "Iphone 14", 1)
	if err != nil {
		return err
	}

	err = createClient(deviceRepository, "Iphone 13", 2)
	if err != nil {
		return err
	}

	err = createClient(deviceRepository, "Galaxy S21", 3)
	if err != nil {
		return err
	}

	err = createClient(deviceRepository, "Mi", 4)
	if err != nil {
		return err
	}

	return nil
}

func createClient(clientRepository repositories.ClientRepository, name string, partnerID int) error {
	ctx := context.Background()
	client, err := clientRepository.FindClientByPartnerID(ctx, partnerID, name)
	if err != nil {
		return err
	}

	_, err = clientRepository.ListClient(ctx)
	if err != nil {
		return err
	}

	if client[0].ID == 0 {
		client[0] = &entities.Client{
			Name:      name,
			CreatedAt: time.Now(),
		}
		err := clientRepository.CreateClient(context.Background(), client[0])
		if err != nil {
			return err
		}
	}

	return nil
}
