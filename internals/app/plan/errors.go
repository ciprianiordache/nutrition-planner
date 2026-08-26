package plan

import "errors"

var (
	ErrNotFound     = errors.New("plan not found")
	ErrInvalidInput = errors.New("invalid plan data")
)
