package day_variant

import "errors"

var (
	ErrNotFound     = errors.New("day variant not found")
	ErrInvalidInput = errors.New("invalid day variant data")
)
