package usecase

import (
	"context"
	"fmt"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/internal/ports"
	"github.com/google/uuid"
)

type QueryListContact struct {
	CreatedBy string
}

type ListContactHandler struct {
	repo ContactRepository
}

func NewListContact(repo ContactRepository) ListContactHandler {
	return ListContactHandler{
		repo: repo,
	}
}

func (h ListContactHandler) List(ctx context.Context, query QueryListContact) ([]*domain.Contact, error) {
	createdBy, err := uuid.Parse(query.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCommand, err)
	}

	return handleRepositoryError(h.repo.List(
		ctx,
		ports.NewFilter(ports.WithCreatedBy(createdBy)),
	))
}
