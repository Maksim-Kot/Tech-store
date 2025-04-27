package repository

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrNotEnough    = errors.New("not enough")
	ErrEditConflict = errors.New("edit conflict")
)
