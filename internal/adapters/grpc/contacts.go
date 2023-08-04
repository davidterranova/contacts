package grpc

import (
	"context"
	"fmt"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/internal/usecase"
	"github.com/google/uuid"
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
		usecase.NewCmdCreateContact(
			req.FirstName,
			req.LastName,
			req.Email,
			req.Phone,
		),
	)
	if err != nil {
		return nil, err
	}

	return &CreateContactResponse{
		Contact: toPBContact(contact),
	}, nil
}

func (h *Handler) UpdateContact(ctx context.Context, req *UpdateContactRequest) (*UpdateContactResponse, error) {
	contactId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid resource uuid: %s", err)
	}

	contact, err := h.app.UpdateContact(
		ctx,
		usecase.NewCmdUpdateContact(
			contactId,
			req.FirstName,
			req.LastName,
			req.Email,
			req.Phone,
		),
	)
	if err != nil {
		return nil, err
	}

	return &UpdateContactResponse{
		Contact: toPBContact(contact),
	}, nil
}

func (h *Handler) DeleteContact(ctx context.Context, req *DeleteContactRequest) (*DeleteContactResponse, error) {
	contactId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid resource uuid: %s", err)
	}

	err = h.app.DeleteContact(
		ctx,
		usecase.NewCmdDeleteContact(contactId),
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
