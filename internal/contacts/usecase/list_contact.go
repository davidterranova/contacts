package usecase

import (
	"context"

	"github.com/davidterranova/contacts/internal/contacts/domain"
	"github.com/davidterranova/contacts/pkg/user"
)

type QueryListContact struct {
	User *user.User
}

type ListContactHandler struct {
	lister ContactLister
}

func NewListContact(lister ContactLister) ListContactHandler {
	return ListContactHandler{
		lister: lister,
	}
}

func (h ListContactHandler) List(ctx context.Context, query QueryListContact) ([]*domain.Contact, error) {
	return h.lister.List(ctx, query)
}
