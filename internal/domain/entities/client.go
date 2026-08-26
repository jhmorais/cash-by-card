package entities

import "time"

type Client struct {
	ID        int    `gorm:"id" json:"id"`
	Name      string `gorm:"size:250" json:"name"`
	PixType   int    `json:"pixType"`
	PixKey    string `json:"pixKey"`
	Phone     string `gorm:"size:15" json:"phone"`
	CPF       string `gorm:"size:15" json:"cpf"`
	PartnerID *int   `gorm:"index" json:"partnerId"`
	Documents string `json:"documents"`

	// Documento único do cliente (PDF/JPEG/PNG) guardado no filesystem
	// em DOCS_DIR/{cpf}.{ext}; aqui ficam só os metadados.
	DocumentName string `gorm:"column:document_name" json:"documentName"`
	DocumentType string `gorm:"column:document_type" json:"documentType"`
	DocumentSize int    `gorm:"column:document_size" json:"documentSize"`

	CreatedAt time.Time
	UpdatedAt time.Time
	Partner   Partner `gorm:"foreignKey:PartnerID" json:"partner"`
}
