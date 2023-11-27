package grpc

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/davidterranova/contacts/internal/contacts/domain"
	"github.com/davidterranova/contacts/internal/contacts/usecase"
	"github.com/davidterranova/contacts/pkg/auth"
	"github.com/davidterranova/contacts/pkg/user"
	"github.com/rs/zerolog/log"
)

const layout = "2006-01-02T15:04:05Z"

type App interface {
	ListContacts(ctx context.Context, cmdIssuedBy user.User) ([]*domain.Contact, error)
	ExportContact(ctx context.Context, cmd usecase.CmdExportContact, cmdIssuedBy user.User) (io.Writer, error)
	CreateContact(ctx context.Context, cmd usecase.CmdCreateContact, cmdIssuedBy user.User) (*domain.Contact, error)
	UpdateContact(ctx context.Context, cmd usecase.CmdUpdateContact, cmdIssuedBy user.User) (*domain.Contact, error)
	DeleteContact(ctx context.Context, cmd usecase.CmdDeleteContact, cmdIssuedBy user.User) error
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
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:list failed to get user from context")
		return nil, err
	}

	contacts, err := h.app.ListContacts(ctx, user)
	if err != nil {
		return nil, err
	}

	return &ListContactsResponse{
		Contacts: toPBContactList(contacts...),
	}, nil
}

func (h *Handler) ExportContact(ctx context.Context, req *ExportContactRequest) (*ExportContactResponse, error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:export failed to get user from context")
		return nil, err
	}

	writer, err := h.app.ExportContact(
		ctx,
		usecase.CmdExportContact{ContactId: req.Id},
		user,
	)
	if err != nil {
		return nil, err
	}

	response := &ExportContactResponse{}
	buf := new(bytes.Buffer)

	_, err = fmt.Fprint(writer, buf)
	if err != nil {
		return nil, err
	}

	response.Vcard = buf.Bytes()
	return response, nil
}

func (h *Handler) CreateContact(ctx context.Context, req *CreateContactRequest) (*CreateContactResponse, error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:create failed to get user from context")
		return nil, err
	}

	contact, err := h.app.CreateContact(
		ctx,
		usecase.CmdCreateContact{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Phone:     req.Phone,
		},
		user,
	)
	if err != nil {
		return nil, err
	}

	return &CreateContactResponse{
		Contact: toPBContact(contact),
	}, nil
}

func (h *Handler) UpdateContact(ctx context.Context, req *UpdateContactRequest) (*UpdateContactResponse, error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:update failed to get user from context")
		return nil, err
	}

	contact, err := h.app.UpdateContact(
		ctx,
		usecase.CmdUpdateContact{
			ContactId: req.Id,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Phone:     req.Phone,
		},
		user,
	)
	if err != nil {
		return nil, err
	}

	return &UpdateContactResponse{
		Contact: toPBContact(contact),
	}, nil
}

func (h *Handler) DeleteContact(ctx context.Context, req *DeleteContactRequest) (*DeleteContactResponse, error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:delete failed to get user from context")
		return nil, err
	}

	err = h.app.DeleteContact(
		ctx,
		usecase.CmdDeleteContact{ContactId: req.Id},
		user,
	)

	return nil, err
}

func (h *Handler) mustEmbedUnimplementedContactsServer() {}

func toPBContact(contact *domain.Contact) *Contact {
	return &Contact{
		Id:               contact.AggregateId().String(),
		CreatedAt:        contact.CreatedAt().Format(layout),
		UpdatedAt:        contact.UpdatedAt().Format(layout),
		FirstName:        contact.FirstName,
		LastName:         contact.LastName,
		Email:            contact.Email,
		Phone:            contact.Phone,
		AggregateVersion: int32(contact.AggregateVersion()),
	}
}

func toPBContactList(contacts ...*domain.Contact) []*Contact {
	var pbContacts = make([]*Contact, 0, len(contacts))
	for _, contact := range contacts {
		pbContacts = append(pbContacts, toPBContact(contact))
	}

	return pbContacts
}
