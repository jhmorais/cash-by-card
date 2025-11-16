package repositories

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type loanRepository struct {
	db *gorm.DB
}

func NewLoanRepository(db *gorm.DB) LoanRepository {
	return &loanRepository{db: db}
}

func (d *loanRepository) CreateLoan(ctx context.Context, entity *entities.Loan) error {
	return d.db.
		Session(&gorm.Session{FullSaveAssociations: false}).
		Create(entity).
		Error
}

func (d *loanRepository) UpdateLoan(ctx context.Context, entity *entities.Loan) error {
	return d.db.
		Session(&gorm.Session{FullSaveAssociations: false}).
		Omit("created_at").
		Save(entity).
		Error
}

func (d *loanRepository) DeleteLoan(ctx context.Context, entity *entities.Loan) error {
	return d.db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Delete(entity).
		Error
}

func (d *loanRepository) FindLoanByID(ctx context.Context, id int) (*entities.Loan, error) {
	var entity *entities.Loan

	err := d.db.
		Preload(clause.Associations).
		Where("id = ?", id).
		Limit(1).
		Find(&entity).Error

	return entity, err
}

func (d *loanRepository) FindLoanByClientID(ctx context.Context, ClientID int) ([]*entities.Loan, error) {
	var entity []*entities.Loan

	err := d.db.
		Preload(clause.Associations).
		Where("client_id = ?", ClientID).
		Limit(100).
		Find(&entity).Error

	return entity, err
}

func (d *loanRepository) FindLoanByPartnerID(ctx context.Context, PartnerID int) ([]*entities.Loan, error) {
	var entity []*entities.Loan

	err := d.db.
		Preload(clause.Associations).
		Where("partner_id = ?", PartnerID).
		Limit(100).
		Find(&entity).Error

	return entity, err
}

func (d *loanRepository) UpdateLoanPaymentStatus(ctx context.Context, LoanID int, paymentStatus string) error {
	err := d.db.
		Model(&entities.Loan{}).
		Where("id = ?", LoanID).
		Update("payment_status", paymentStatus).
		Error
	return err
}

func (d *loanRepository) ListLoan(ctx context.Context) ([]*entities.Loan, error) {
	//TODO impl pagination
	var entities []*entities.Loan

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
