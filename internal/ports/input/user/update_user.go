package input

type UpdateUser struct {
	ID       int    `json:"id,omitempty"`
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"-"`
	// Role     string `json:"role"`
}
