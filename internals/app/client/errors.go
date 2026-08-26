package client

import "errors"

var (
	ErrNotFound     = errors.New("client not found")
	ErrInvalidInput = errors.New("invalid client data")
)
