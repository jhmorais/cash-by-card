package input

type UpdateClient struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name"`
	PixType   int    `json:"pix_type"`
	PixKey    string `json:"pix_key"`
	PartnerID int    `json:"partner_id"`
}
