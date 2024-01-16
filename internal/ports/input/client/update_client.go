package input

type UpdateClient struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name"`
	PixType   int    `json:"pixType"`
	PixKey    string `json:"pixKey"`
	PartnerID int    `json:"partnerId"`
	Phone     string `json:"phone"`
	CPF       string `json:"cpf"`
	Documents string `json:"documents"`
}
