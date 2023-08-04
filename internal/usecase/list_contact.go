package usecase

import (
	"context"

	"github.com/davidterranova/contacts/internal/domain"
)

type QueryListContact struct{}

type ListContactHandler struct {
	lister ContactLister
}

func NewListContact(lister ContactLister) ListContactHandler {
	return ListContactHandler{
		lister: lister,
	}
}

func (h ListContactHandler) List(ctx context.Context, query QueryListContact) ([]*domain.Contact, error) {
	return h.lister.List(ctx)
}
