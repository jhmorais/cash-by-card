package repositories

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type cardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) CardRepository {
	return &cardRepository{db: db}
}

func (d *cardRepository) CreateCard(ctx context.Context, entity *entities.Card) error {
	return d.db.
		Session(&gorm.Session{FullSaveAssociations: false}).
		Create(entity).
		Error
}

func (d *cardRepository) UpdateCard(ctx context.Context, entity *entities.Card) error {
	return d.db.
		Session(&gorm.Session{FullSaveAssociations: false}).
		Save(entity).
		Error
}

func (d *cardRepository) DeleteCard(ctx context.Context, entity *entities.Card) error {
	return d.db.
		Session(&gorm.Session{FullSaveAssociations: false}).
		Delete(entity).
		Error
}

func (d *cardRepository) FindCardByID(ctx context.Context, id int) (*entities.Card, error) {
	var entity *entities.Card

	err := d.db.
		Preload(clause.Associations).
		Where("id = ?", id).
		Limit(1).
		Find(&entity).Error

	return entity, err
}

func (d *cardRepository) FindCardByLoanID(ctx context.Context, LoanID int) ([]*entities.Card, error) {
	var entity []*entities.Card

	err := d.db.
		Preload(clause.Associations).
		Where("loan_id = ?", LoanID).
		Find(&entity).Error

	return entity, err
}

func (d *cardRepository) ListCard(ctx context.Context) ([]*entities.Card, error) {
	//TODO impl pagination
	var entities []*entities.Card

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
