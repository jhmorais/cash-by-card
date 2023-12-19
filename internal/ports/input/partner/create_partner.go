package input

type CreatePartner struct {
	Name    string `json:"name"`
	CPF     string `json:"cpf"`
	PixKey  string `json:"pixKey"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	PixType int    `json:"pixType"`
	Email   string `json:"email"`
}
