package input

type CreateClient struct {
	Name      string `json:"name"`
	PixType   int    `json:"pix_type"`
	PixKey    string `json:"pix_key"`
	PartnerID int    `json:"partner_id"`
}
