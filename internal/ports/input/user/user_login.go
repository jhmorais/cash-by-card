package input

type UserLogin struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	// Role     string `json:"role"`
}
