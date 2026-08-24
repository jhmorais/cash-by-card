package input

type ResetPassword struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}
