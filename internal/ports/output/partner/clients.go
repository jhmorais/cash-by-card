package output

import "github.com/jhmorais/cash-by-card/internal/domain/entities"

type PartnerClients struct {
	Clients []*entities.Client `json:"clients"`
}

type PartnerCreateClient struct {
	ClientID int `json:"clientId"`
}

type PartnerUpdateClient struct {
	ClientID int `json:"clientId"`
}
