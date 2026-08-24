package repositories

import (
	"context"
	"time"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	"gorm.io/gorm"
)

type passwordResetTokenRepository struct {
	db *gorm.DB
}

func NewPasswordResetTokenRepository(db *gorm.DB) PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

func (r *passwordResetTokenRepository) CreateToken(ctx context.Context, entity *entities.PasswordResetToken) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *passwordResetTokenRepository) FindValidTokenByHash(ctx context.Context, tokenHash string) (*entities.PasswordResetToken, error) {
	var entity *entities.PasswordResetToken
	err := r.db.WithContext(ctx).
		Where("token_hash = ? AND used_at IS NULL AND expires_at > ?", tokenHash, time.Now()).
		Limit(1).
		Find(&entity).Error
	if err != nil {
		return nil, err
	}
	if entity == nil || entity.ID == 0 {
		return nil, nil
	}
	return entity, nil
}

func (r *passwordResetTokenRepository) MarkTokenUsed(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&entities.PasswordResetToken{}).
		Where("id = ?", id).
		Update("used_at", time.Now()).Error
}

func (r *passwordResetTokenRepository) DeleteTokensByUser(ctx context.Context, userID int) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&entities.PasswordResetToken{}).Error
}
