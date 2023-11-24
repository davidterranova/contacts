package usecase

import (
	"context"

	"github.com/davidterranova/contacts/internal/contacts/domain"
	"github.com/davidterranova/contacts/pkg/user"
)

type ListContactHandler struct {
	lister ContactReadModel
}

func NewListContact(lister ContactReadModel) ListContactHandler {
	return ListContactHandler{
		lister: lister,
	}
}

func (h ListContactHandler) List(ctx context.Context, cmdIssuedBy user.User) ([]*domain.Contact, error) {
	return h.lister.List(ctx, QueryContact{
		Requestor: cmdIssuedBy,
	})
}
