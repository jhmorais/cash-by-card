package repositories

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type partnerRepository struct {
	db *gorm.DB
}

func NewPartnerRepository(db *gorm.DB) PartnerRepository {
	return &partnerRepository{db: db}
}

func (d *partnerRepository) CreatePartner(ctx context.Context, entity *entities.Partner) error {
	return d.db.
		Session(&gorm.Session{FullSaveAssociations: false}).
		Create(entity).
		Error
}

func (d *partnerRepository) UpdatePartner(ctx context.Context, entity *entities.Partner) error {
	return d.db.
		Session(&gorm.Session{FullSaveAssociations: false}).
		Omit("created_at").
		Save(entity).
		Error
}

func (d *partnerRepository) DeletePartner(ctx context.Context, entity *entities.Partner) error {
	return d.db.
		Session(&gorm.Session{FullSaveAssociations: false}).
		Delete(entity).
		Error
}

func (d *partnerRepository) FindPartnerByID(ctx context.Context, id int) (*entities.Partner, error) {
	var entity *entities.Partner

	err := d.db.
		Preload(clause.Associations).
		Where("id = ?", id).
		Limit(1).
		Find(&entity).Error

	return entity, err
}

func (d *partnerRepository) FindPartnerByName(ctx context.Context, name string) ([]*entities.Partner, error) {
	var entity []*entities.Partner

	err := d.db.
		Preload(clause.Associations).
		Where("name = ?", name).
		Limit(100).
		Find(&entity).Error

	return entity, err
}

func (d *partnerRepository) FindPartnerByEmail(ctx context.Context, email string) ([]*entities.Partner, error) {
	var entity []*entities.Partner

	err := d.db.
		Preload(clause.Associations).
		Where("email = ?", email).
		Limit(100).
		Find(&entity).Error

	return entity, err
}

func (d *partnerRepository) FindPartnerByCPF(ctx context.Context, cpf string) ([]*entities.Partner, error) {
	var entity []*entities.Partner

	err := d.db.
		Preload(clause.Associations).
		Where("cpf = ?", cpf).
		Limit(100).
		Find(&entity).Error

	return entity, err
}

func (d *partnerRepository) FindPartnerByPartnerID(ctx context.Context, partnerID int, name string) ([]*entities.Partner, error) {
	var entity []*entities.Partner

	err := d.db.
		Preload(clause.Associations).
		Where("partner_id = ?", partnerID).
		Where("name = ?", name).
		Limit(1).
		Find(&entity).Error

	return entity, err
}

func (d *partnerRepository) ListPartner(ctx context.Context) ([]*entities.Partner, error) {
	//TODO impl pagination
	var entities []*entities.Partner

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
