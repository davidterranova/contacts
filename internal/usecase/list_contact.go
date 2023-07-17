package usecase

import (
	"context"

	"github.com/davidterranova/contacts/internal/domain"
)

type QueryListContact struct{}

type ListContactHandler struct {
	repo ContactRepository
}

func NewListContact(repo ContactRepository) ListContactHandler {
	return ListContactHandler{
		repo: repo,
	}
}

func (h ListContactHandler) List(ctx context.Context, query QueryListContact) ([]*domain.Contact, error) {
	return handleRepositoryError(h.repo.List(ctx))
}
