package food

import "errors"

var (
	ErrNotFound     = errors.New("food not found")
	ErrInvalidInput = errors.New("invalid food data")
)
