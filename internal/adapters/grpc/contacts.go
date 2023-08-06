package grpc

import (
	"context"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/internal/usecase"
)

const layout = "2006-01-02T15:04:05Z"

type App interface {
	ListContacts(ctx context.Context, query usecase.QueryListContact) ([]*domain.Contact, error)
	CreateContact(ctx context.Context, cmd usecase.CmdCreateContact) (*domain.Contact, error)
	UpdateContact(ctx context.Context, cmd usecase.CmdUpdateContact) (*domain.Contact, error)
	DeleteContact(ctx context.Context, cmd usecase.CmdDeleteContact) error
}

type Handler struct {
	app App
}

func NewHandler(app App) *Handler {
	return &Handler{
		app: app,
	}
}

func (h *Handler) ListContacts(ctx context.Context, req *ListContactsRequest) (*ListContactsResponse, error) {
	contacts, err := h.app.ListContacts(ctx, usecase.QueryListContact{})
	if err != nil {
		return nil, err
	}

	return &ListContactsResponse{
		Contacts: toPBContactList(contacts...),
	}, nil
}

func (h *Handler) CreateContact(ctx context.Context, req *CreateContactRequest) (*CreateContactResponse, error) {
	contact, err := h.app.CreateContact(
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

	return &CreateContactResponse{
		Contact: toPBContact(contact),
	}, nil
}

func (h *Handler) UpdateContact(ctx context.Context, req *UpdateContactRequest) (*UpdateContactResponse, error) {
	contact, err := h.app.UpdateContact(
		ctx,
		usecase.CmdUpdateContact{
			ContactId: req.Id,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Phone:     req.Phone,
		},
	)
	if err != nil {
		return nil, err
	}

	return &UpdateContactResponse{
		Contact: toPBContact(contact),
	}, nil
}

func (h *Handler) DeleteContact(ctx context.Context, req *DeleteContactRequest) (*DeleteContactResponse, error) {
	err := h.app.DeleteContact(
		ctx,
		usecase.CmdDeleteContact{ContactId: req.Id},
	)

	return nil, err
}

func (h *Handler) mustEmbedUnimplementedContactsServer() {}

func toPBContact(contact *domain.Contact) *Contact {
	return &Contact{
		Id:        contact.Id.String(),
		CreatedAt: contact.CreatedAt.Format(layout),
		UpdatedAt: contact.UpdatedAt.Format(layout),
		FirstName: contact.FirstName,
		LastName:  contact.LastName,
		Email:     contact.Email,
		Phone:     contact.Phone,
	}
}

func toPBContactList(contacts ...*domain.Contact) []*Contact {
	var pbContacts = make([]*Contact, 0, len(contacts))
	for _, contact := range contacts {
		pbContacts = append(pbContacts, toPBContact(contact))
	}

	return pbContacts
}
