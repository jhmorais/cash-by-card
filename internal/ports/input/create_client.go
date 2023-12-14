package input

type CreateClient struct {
	Name      string `json:"name"`
	PixType   int    `json:"pixType"`
	PixKey    string `json:"pixKey"`
	PartnerID int    `json:"partnerId"`
}
