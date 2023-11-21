package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.35

import (
	"context"

	"github.com/davidterranova/contacts/internal/contacts/adapters/graphql/model"
	"github.com/davidterranova/contacts/internal/contacts/usecase"
	"github.com/davidterranova/contacts/pkg/auth"
	"github.com/rs/zerolog/log"
)

// CreateContact is the resolver for the createContact field.
func (r *mutationResolver) CreateContact(ctx context.Context, input model.NewContact) (*model.Contact, error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:create failed to get user from context")
		return nil, auth.ErrUnauthorized
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

// UpdateContact is the resolver for the updateContact field.
func (r *mutationResolver) UpdateContact(ctx context.Context, id string, input model.NewContact) (*model.Contact, error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:update failed to get user from context")
		return nil, auth.ErrUnauthorized
	}

	contact, err := r.app.UpdateContact(
		ctx,
		usecase.CmdUpdateContact{
			ContactId: id,
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

// DeleteContact is the resolver for the deleteContact field.
func (r *mutationResolver) DeleteContact(ctx context.Context, id string) (*model.Contact, error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:delete failed to get user from context")
		return nil, auth.ErrUnauthorized
	}

	err = r.app.DeleteContact(ctx, usecase.CmdDeleteContact{ContactId: id}, user)
	if err != nil {
		return nil, err
	}

	return &model.Contact{ID: id}, nil
}

// ListContacts is the resolver for the listContacts field.
func (r *queryResolver) ListContacts(ctx context.Context) ([]*model.Contact, error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:list failed to get user from context")
		return nil, auth.ErrUnauthorized
	}

	contacts, err := r.app.ListContacts(ctx, usecase.QueryListContact{
		User: user,
	})
	if err != nil {
		return nil, err
	}

	return toGQLContacts(contacts), nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
