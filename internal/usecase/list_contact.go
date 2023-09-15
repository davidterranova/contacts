package usecase

import (
	"context"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/internal/ports"
)

type QueryListContact struct {
	CreatedBy domain.User
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
	return handleRepositoryError(h.repo.List(
		ctx,
		ports.NewFilter(ports.WithCreatedBy(query.CreatedBy.Id)),
	))
}
