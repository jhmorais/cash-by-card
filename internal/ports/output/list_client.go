package output

import (
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
)

type ListClient struct {
	Clients []*entities.Client
}
