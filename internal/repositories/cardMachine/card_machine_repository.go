package repositories

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type cardMachineRepository struct {
	db *gorm.DB
}

func NewCardMachineRepository(db *gorm.DB) CardMachineRepository {
	return &cardMachineRepository{db: db}
}

func (d *cardMachineRepository) CreateCardMachine(ctx context.Context, entity *entities.CardMachine) error {
	return d.db.
		Session(&gorm.Session{FullSaveAssociations: false}).
		Create(entity).
		Error
}

func (d *cardMachineRepository) UpdateCardMachine(ctx context.Context, entity *entities.CardMachine) error {
	return d.db.
		Session(&gorm.Session{FullSaveAssociations: false}).
		Omit("created_at").
		Save(entity).
		Error
}

func (d *cardMachineRepository) DeleteCardMachine(ctx context.Context, entity *entities.CardMachine) error {
	return d.db.
		Session(&gorm.Session{FullSaveAssociations: false}).
		Delete(entity).
		Error
}

func (d *cardMachineRepository) FindCardMachineByID(ctx context.Context, id int) (*entities.CardMachine, error) {
	var entity *entities.CardMachine

	err := d.db.
		Preload(clause.Associations).
		Where("id = ?", id).
		Limit(1).
		Find(&entity).Error

	return entity, err
}

func (d *cardMachineRepository) ListCardMachine(ctx context.Context) ([]*entities.CardMachine, error) {
	//TODO impl pagination
	var entities []*entities.CardMachine

	err := d.db.
		Preload(clause.Associations).
		Limit(100).
		Order("created_at desc").
		Find(&entities).Error

	if err != nil {
		return nil, err
	}

	return entities, nil
}
