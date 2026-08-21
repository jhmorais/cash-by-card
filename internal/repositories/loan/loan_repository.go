package repositories

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/loan"
	dashboard "github.com/jhmorais/cash-by-card/internal/ports/output/dashboard"
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

func (d *loanRepository) ListLoan(ctx context.Context, filter *input.ListLoanFilter, pagination *input.Pagination) ([]*entities.Loan, int64, error) {
	var loans []*entities.Loan
	var total int64

	if pagination == nil {
		pagination = &input.Pagination{Page: 1, Limit: 10}
	}

	countQuery := applyLoanFilters(d.db.WithContext(ctx).Model(&entities.Loan{}), filter)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	findQuery := applyLoanFilters(d.db.WithContext(ctx).Model(&entities.Loan{}), filter)
	offset := (pagination.Page - 1) * pagination.Limit
	err := findQuery.
		Preload(clause.Associations).
		Order("loan.created_at desc, loan.id desc").
		Offset(offset).
		Limit(pagination.Limit).
		Find(&loans).Error
	if err != nil {
		return nil, 0, err
	}

	return loans, total, nil
}

// applyLoanFilters appends conditional WHERE/JOIN clauses based on non-nil filter fields.
func applyLoanFilters(db *gorm.DB, filter *input.ListLoanFilter) *gorm.DB {
	if filter == nil {
		return db
	}

	if filter.PaymentStatus != nil {
		db = db.Where("loan.payment_status = ?", *filter.PaymentStatus)
	}
	clientFilter := (filter.ClientName != nil && *filter.ClientName != "") ||
		(filter.ClientCPF != nil && *filter.ClientCPF != "")
	if clientFilter {
		db = db.Joins("JOIN client ON client.id = loan.client_id")
	}
	if filter.ClientName != nil && *filter.ClientName != "" {
		db = db.Where("client.name LIKE ?", "%"+*filter.ClientName+"%")
	}
	if filter.ClientCPF != nil && *filter.ClientCPF != "" {
		db = db.Where("client.cpf = ?", *filter.ClientCPF)
	}
	partnerFilter := (filter.PartnerName != nil && *filter.PartnerName != "") ||
		(filter.PartnerCPF != nil && *filter.PartnerCPF != "")
	if partnerFilter {
		db = db.Joins("JOIN partner ON partner.id = loan.partner_id")
	}
	if filter.PartnerName != nil && *filter.PartnerName != "" {
		db = db.Where("partner.name LIKE ?", "%"+*filter.PartnerName+"%")
	}
	if filter.PartnerCPF != nil && *filter.PartnerCPF != "" {
		db = db.Where("partner.cpf = ?", *filter.PartnerCPF)
	}

	return db
}

func (d *loanRepository) GetTotals(ctx context.Context, month int, year int) (*dashboard.Dashboard, error) {
	var result dashboard.Dashboard

	err := d.db.
		WithContext(ctx).
		Model(&entities.Loan{}).
		Select(`
            COUNT(*) AS total_loans,
            COALESCE(SUM(client_amount), 0) AS total_value,
            COALESCE(SUM(gross_profit), 0) AS gross_profit,
            COALESCE(SUM(profit), 0) AS profit,
			COALESCE(SUM(partner_amount), 0) AS partner_profit,
			COALESCE(SUM(amount), 0) AS machine_amount
        `).
		Where("EXTRACT(YEAR FROM created_at) = ?", year).
		Where("MONTH(created_at) = ?", month).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *loanRepository) GetBestPartners(ctx context.Context, month int, year int) ([]dashboard.BestPartner, error) {
	var result []dashboard.BestPartner

	err := d.db.WithContext(ctx).
		Table("loan").
		Select(`
            partner.name AS partner,
            COUNT(loan.id) AS qtt,
			COALESCE(SUM(loan.partner_amount), 0) AS profit
        `).
		Joins("INNER JOIN partner ON partner.id = loan.partner_id").
		Where("EXTRACT(YEAR FROM loan.created_at) = ?", year).
		Where("EXTRACT(MONTH FROM loan.created_at) = ?", month).
		// Where("MONTH(loan.created_at) = ?", month).
		Group("partner.name").
		Order("qtt DESC").
		Limit(5).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *loanRepository) GetMonthlyLoans(ctx context.Context, year int) (*dashboard.MonthlyLoans, error) {
	loans := &dashboard.MonthlyLoans{
		Labels: []string{},
		Total:  []float64{},
		Gross:  []float64{},
		Net:    []float64{},
	}

	type row struct {
		Month         int
		Total         float64
		Gross         float64
		Net           float64
		PartnerProfit float64
		MachineAmount float64
	}

	var rows []row

	err := d.db.WithContext(ctx).
		Table("loan").
		Select(`
            MONTH(created_at) AS month,
            COALESCE(SUM(client_amount), 0) AS total,
            COALESCE(SUM(gross_profit), 0) AS gross,
            COALESCE(SUM(profit), 0) AS net,
			COALESCE(SUM(partner_amount), 0) AS partner_profit,
			COALESCE(SUM(amount), 0) AS machine_amount
        `).
		Where("EXTRACT(YEAR FROM created_at) = ?", year).
		Group("month").
		Order("month").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	// Prepara arrays no formato que o front espera
	monthNames := []string{"Jan", "Fev", "Mar", "Abr", "Mai", "Jun", "Jul", "Ago", "Set", "Out", "Nov", "Dez"}

	for _, r := range rows {
		loans.Labels = append(loans.Labels, monthNames[r.Month-1])
		loans.Total = append(loans.Total, r.Total)
		loans.Gross = append(loans.Gross, r.Gross)
		loans.Net = append(loans.Net, r.Net)
		loans.PartnerProfit = append(loans.PartnerProfit, r.PartnerProfit)
		loans.MachineAmount = append(loans.MachineAmount, r.MachineAmount)
	}

	return loans, nil
}
