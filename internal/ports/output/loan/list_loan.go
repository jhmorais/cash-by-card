package output

import (
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
)

type ListLoan struct {
	Loans      []*entities.Loan `json:"loans"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"totalPages"`
}
