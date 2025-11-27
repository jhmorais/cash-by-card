package repositories

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	dashboard "github.com/jhmorais/cash-by-card/internal/ports/output/dashboard"
)

type LoanRepository interface {
	CreateLoan(ctx context.Context, entity *entities.Loan) error
	DeleteLoan(ctx context.Context, entity *entities.Loan) error
	UpdateLoan(ctx context.Context, entity *entities.Loan) error
	FindLoanByID(ctx context.Context, id int) (*entities.Loan, error)
	FindLoanByClientID(ctx context.Context, ClientID int) ([]*entities.Loan, error)
	FindLoanByPartnerID(ctx context.Context, PartnerID int) ([]*entities.Loan, error)
	ListLoan(ctx context.Context) ([]*entities.Loan, error)
	UpdateLoanPaymentStatus(ctx context.Context, LoanID int, paymentStatus string) error
	GetTotals(ctx context.Context, month int, year int) (*dashboard.Dashboard, error)
	GetBestPartners(ctx context.Context, month int, year int) ([]dashboard.BestPartner, error)
	GetMonthlyLoans(ctx context.Context, year int) (*dashboard.MonthlyLoans, error)
}
