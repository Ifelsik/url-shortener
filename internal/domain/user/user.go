package user

import "time"

type User struct {
	ID        uint64
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}
