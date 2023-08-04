package usecase

import (
	"errors"
)

var (
	ErrInternal       = errors.New("internal error")
	ErrInvalidCommand = errors.New("invalid command")
	ErrNotFound       = errors.New("not found")
)
