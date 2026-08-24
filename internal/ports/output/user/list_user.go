package output

type UserItem struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	Role               string `json:"role"`
	PendingFirstAccess bool   `json:"pendingFirstAccess"`
}

type ListUser struct {
	Users []UserItem `json:"users"`
}
