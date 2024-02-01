package repositories

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/domain/entities"
)

type LoanRepository interface {
	CreateLoan(ctx context.Context, entity *entities.Loan) error
	DeleteLoan(ctx context.Context, entity *entities.Loan) error
	UpdateLoan(ctx context.Context, entity *entities.Loan) error
	FindLoanByID(ctx context.Context, id int) (*entities.Loan, error)
	FindLoanByClientID(ctx context.Context, ClientID int) ([]*entities.Loan, error)
	FindLoanByPartnerID(ctx context.Context, PartnerID int) ([]*entities.Loan, error)
	ListLoan(ctx context.Context) ([]*entities.Loan, error)
}
