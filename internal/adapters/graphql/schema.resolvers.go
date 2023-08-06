package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.35

import (
	"context"

	"github.com/davidterranova/contacts/internal/adapters/graphql/model"
	"github.com/davidterranova/contacts/internal/usecase"
)

// CreateContact is the resolver for the createContact field.
func (r *mutationResolver) CreateContact(ctx context.Context, req model.NewContact) (*model.Contact, error) {
	contact, err := r.app.CreateContact(
		ctx,
		usecase.CmdCreateContact{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Phone:     req.Phone,
		},
	)
	if err != nil {
		return nil, err
	}

	return toGQLContact(contact), nil
}

// UpdateContact is the resolver for the updateContact field.
func (r *mutationResolver) UpdateContact(ctx context.Context, id string, req model.NewContact) (*model.Contact, error) {
	contact, err := r.app.UpdateContact(
		ctx,
		usecase.CmdUpdateContact{
			ContactId: id,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Phone:     req.Phone,
		},
	)
	if err != nil {
		return nil, err
	}

	return toGQLContact(contact), nil
}

// DeleteContact is the resolver for the deleteContact field.
func (r *mutationResolver) DeleteContact(ctx context.Context, id string) (*model.Contact, error) {
	err := r.app.DeleteContact(ctx, usecase.CmdDeleteContact{ContactId: id})
	if err != nil {
		return nil, err
	}

	return &model.Contact{ID: id}, nil
}

// ListContacts is the resolver for the listContacts field.
func (r *queryResolver) ListContacts(ctx context.Context) ([]*model.Contact, error) {
	contacts, err := r.app.ListContacts(ctx, usecase.QueryListContact{})
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
