package usecase

import (
	"errors"
	"fmt"

	"github.com/davidterranova/contacts/internal/contacts/domain"
	"github.com/davidterranova/cqrs/eventsourcing"
)

var (
	ErrInternal       = errors.New("internal error")
	ErrInvalidCommand = errors.New("invalid command")
	ErrNotFound       = errors.New("not found")
	ErrForbidden      = errors.New("forbidden")
)

func handleErrs(contact *domain.Contact, err error) (*domain.Contact, error) {
	switch {
	case errors.Is(err, eventsourcing.ErrAggregateNotFound), errors.Is(err, ErrNotFound):
		return nil, fmt.Errorf("%w: %s", ErrNotFound, err)
	case errors.Is(err, ErrForbidden):
		return nil, err
	case err != nil:
		return nil, fmt.Errorf("%w: %s", ErrInternal, err)
	default:
		return contact, nil
	}
}
