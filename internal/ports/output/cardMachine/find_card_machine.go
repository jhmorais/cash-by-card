package output

import (
	"time"
)

type FindCardMachine struct {
	CardMachine *FindCardMachineOutput // adicionar os mesmos campos do entity mas as taxas serao do tipo string
}

type FindCardMachineOutput struct {
	ID            int                    `gorm:"id" json:"id"`
	Brand         string                 `gorm:"size:250" json:"brand"`
	Name          string                 `gorm:"size:250" json:"name"`
	PresentialTax map[string]interface{} `json:"presentialTax"` // quando for devolver para o front fazer um parse para string, JHSON.parse no front para acessar ao json
	OnlineTax     map[string]interface{} `json:"onlineTax"`     // quando for devolver para o front fazer um parse para string, JHSON.parse no front para acessar ao json
	Installments  int                    `gorm:"index" json:"installments"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
