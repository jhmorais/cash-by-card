package output

type FindUser struct {
	ID    int    `gorm:"id" json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}
