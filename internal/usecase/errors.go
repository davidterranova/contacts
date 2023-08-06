package usecase

import (
	"errors"
	"fmt"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/pkg/eventsourcing"
)

var (
	ErrInternal       = errors.New("internal error")
	ErrInvalidCommand = errors.New("invalid command")
	ErrNotFound       = errors.New("not found")
)

func handleErrs(contact *domain.Contact, err error) (*domain.Contact, error) {
	switch {
	case errors.Is(err, eventsourcing.ErrAggregateNotFound):
		return nil, fmt.Errorf("%w: %s", ErrNotFound, err)
	case err != nil:
		return nil, fmt.Errorf("%w: %s", ErrInternal, err)
	default:
		return contact, nil
	}
}
