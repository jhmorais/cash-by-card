package output

import (
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
)

type ListLoan struct {
	Loans []*entities.Loan
}
