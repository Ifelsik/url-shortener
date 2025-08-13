package url

import (
	"errors"
)

var (
	ErrEmptyRequest = errors.New("empty request")

	ErrNotFound = errors.New("resource not found")
)
