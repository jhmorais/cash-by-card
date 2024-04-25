package input

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"-"`
	Role     string `json:"role"`
}
