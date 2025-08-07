package identifier

import "github.com/google/uuid"

// Interface Identifier describes an methods
// for generating unique identifiers
type Identifier interface {
	String() string
}	

// struct uuidProvider implements Identifier via UUID
type uuidProvider struct{}

func NewUUIDProvider() *uuidProvider {
	return &uuidProvider{}
}

func (u *uuidProvider) String() string {
	return uuid.NewString()
}
