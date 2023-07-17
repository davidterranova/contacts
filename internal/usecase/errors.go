package usecase

import (
	"errors"
	"fmt"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/internal/ports"
)

var (
	ErrInternal       = errors.New("internal error")
	ErrInvalidCommand = errors.New("invalid command")
	ErrNotFound       = errors.New("not found")
)

type contactResponse interface {
	*domain.Contact | []*domain.Contact
}

func handleRepositoryError[T contactResponse](c T, err error) (T, error) {
	if err == nil {
		return c, nil
	}

	switch {
	case errors.Is(err, ports.ErrNotFound):
		return nil, fmt.Errorf("%w: %s", ErrNotFound, err)
	default:
		return nil, fmt.Errorf("%w: %s", ErrInternal, err)
	}
}
