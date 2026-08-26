package supplement

import "errors"

var (
	ErrNotFound     = errors.New("supplement not found")
	ErrInvalidInput = errors.New("invalid supplement data")
)
