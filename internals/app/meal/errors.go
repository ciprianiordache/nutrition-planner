package meal

import "errors"

var (
	ErrNotFound     = errors.New("meal not found")
	ErrInvalidInput = errors.New("invalid meal data")
)
