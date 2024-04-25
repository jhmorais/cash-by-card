package repositories

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (d *userRepository) CreateUser(ctx context.Context, entity *entities.User) error {
	return d.db.
		Session(&gorm.Session{FullSaveAssociations: false}).
		Create(entity).
		Error
}

func (d *userRepository) UpdateUser(ctx context.Context, entity *entities.User) error {
	return d.db.
		Session(&gorm.Session{FullSaveAssociations: false}).
		Omit("created_at").
		Save(entity).
		Error
}

func (d *userRepository) DeleteUser(ctx context.Context, entity *entities.User) error {
	return d.db.
		Session(&gorm.Session{FullSaveAssociations: false}).
		Delete(entity).
		Error
}

func (d *userRepository) FindUserByID(ctx context.Context, id int) (*entities.User, error) {
	var entity *entities.User

	err := d.db.
		Preload(clause.Associations).
		Where("id = ?", id).
		Limit(1).
		Find(&entity).Error

	return entity, err
}

func (d *userRepository) FindUserByEmail(ctx context.Context, email string) ([]*entities.User, error) {
	var entity []*entities.User

	err := d.db.
		Preload(clause.Associations).
		Where("email = ?", email).
		Limit(100).
		Find(&entity).Error

	return entity, err
}

func (d *userRepository) FindUserByEmailandPassword(ctx context.Context, email string, password string) ([]*entities.User, error) {
	var entity []*entities.User

	err := d.db.
		Preload(clause.Associations).
		Where("email = ? AND password = ?", email, password).
		Limit(100).
		Find(&entity).Error

	return entity, err
}

func (d *userRepository) ListUser(ctx context.Context) ([]*entities.User, error) {
	//TODO impl pagination
	var entities []*entities.User

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
