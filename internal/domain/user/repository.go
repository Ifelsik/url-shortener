package user

type UserRepository interface {
	Add(user *User) error
	GetByToken(token string) (*User, error)
}
