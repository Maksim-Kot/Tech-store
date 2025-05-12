package repository

import "errors"

var (
	ErrNotFound   = errors.New("order not found")
	ErrNotCreated = errors.New("order not created")
)
