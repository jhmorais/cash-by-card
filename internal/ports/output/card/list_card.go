package output

import (
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
)

type ListCard struct {
	Cards []*entities.Card
}
