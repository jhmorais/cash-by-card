package output

import (
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
)

type ListCardMachine struct {
	CardMachines []*entities.CardMachine
}
