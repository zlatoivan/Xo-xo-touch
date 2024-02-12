package user

type User struct {
	ID       string
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRepo interface {
	Register(username, password string) (*User, error)
	Authorize(username, password string) (*User, error)
}
