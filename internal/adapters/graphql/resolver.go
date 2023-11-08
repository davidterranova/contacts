package graphql

import (
	"context"

	"github.com/davidterranova/contacts/internal/adapters/graphql/model"
	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/internal/usecase"
	"github.com/davidterranova/contacts/pkg/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type App interface {
	ListContacts(ctx context.Context, query usecase.QueryListContact) ([]*domain.Contact, error)
	CreateContact(ctx context.Context, cmd usecase.CmdCreateContact, cmdIssuedBy user.User) (*domain.Contact, error)
	UpdateContact(ctx context.Context, cmd usecase.CmdUpdateContact, cmdIssuedBy user.User) (*domain.Contact, error)
	DeleteContact(ctx context.Context, cmd usecase.CmdDeleteContact, cmdIssuedBy user.User) error
}

type Resolver struct {
	app App
}

func NewResolver(app App) *Resolver {
	return &Resolver{
		app: app,
	}
}

func (r *Resolver) ListContacts(ctx context.Context) ([]*model.Contact, error) {
	contacts, err := r.app.ListContacts(ctx, usecase.QueryListContact{})
	if err != nil {
		return nil, err
	}

	return toGQLContacts(contacts), nil
}

func (r *Resolver) CreateContact(ctx context.Context, input model.NewContact) (*model.Contact, error) {
	contact, err := r.app.CreateContact(
		ctx,
		usecase.CmdCreateContact{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			Phone:     input.Phone,
		},
		user.NewUnauthenticated(), // TODO: to fix
	)
	if err != nil {
		return nil, err
	}

	return toGQLContact(contact), nil
}

func toGQLContact(contact *domain.Contact) *model.Contact {
	return &model.Contact{
		ID:        contact.Id.String(),
		CreatedAt: contact.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: contact.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		FirstName: contact.FirstName,
		LastName:  contact.LastName,
		Email:     contact.Email,
		Phone:     contact.Phone,
	}
}

func toGQLContacts(contacts []*domain.Contact) []*model.Contact {
	var gqlContacts = make([]*model.Contact, 0, len(contacts))
	for _, contact := range contacts {
		gqlContacts = append(gqlContacts, toGQLContact(contact))
	}

	return gqlContacts
}
