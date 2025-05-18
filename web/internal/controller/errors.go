package controller

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrNotEnough          = errors.New("not enough")
	ErrEditConflict       = errors.New("edit conflict")
)
