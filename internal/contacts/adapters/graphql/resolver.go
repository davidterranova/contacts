package graphql

import (
	"context"
	"io"

	"github.com/davidterranova/contacts/internal/contacts/adapters/graphql/model"
	"github.com/davidterranova/contacts/internal/contacts/domain"
	"github.com/davidterranova/contacts/internal/contacts/usecase"
	"github.com/davidterranova/contacts/pkg/auth"
	"github.com/davidterranova/contacts/pkg/user"
	"github.com/rs/zerolog/log"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type App interface {
	ListContacts(ctx context.Context, cmdIssuedBy user.User) ([]*domain.Contact, error)
	ExportContact(ctx context.Context, cmd usecase.CmdExportContact, cmdIssuedBy user.User) (io.Writer, error)
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
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:list failed to get user from context")
		return nil, err
	}

	contacts, err := r.app.ListContacts(ctx, user)
	if err != nil {
		return nil, err
	}

	return toGQLContacts(contacts), nil
}

func (r *Resolver) CreateContact(ctx context.Context, input model.NewContact) (*model.Contact, error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:list failed to get user from context")
		return nil, err
	}

	contact, err := r.app.CreateContact(
		ctx,
		usecase.CmdCreateContact{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			Phone:     input.Phone,
		},
		user,
	)
	if err != nil {
		return nil, err
	}

	return toGQLContact(contact), nil
}

func toGQLContact(contact *domain.Contact) *model.Contact {
	return &model.Contact{
		ID:               contact.AggregateId().String(),
		CreatedAt:        contact.CreatedAt().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        contact.UpdatedAt().Format("2006-01-02T15:04:05Z"),
		FirstName:        contact.FirstName,
		LastName:         contact.LastName,
		Email:            contact.Email,
		Phone:            contact.Phone,
		AggregateVersion: contact.AggregateVersion(),
	}
}

func toGQLContacts(contacts []*domain.Contact) []*model.Contact {
	var gqlContacts = make([]*model.Contact, 0, len(contacts))
	for _, contact := range contacts {
		gqlContacts = append(gqlContacts, toGQLContact(contact))
	}

	return gqlContacts
}
