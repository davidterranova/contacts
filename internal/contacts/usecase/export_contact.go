package usecase

import (
	"context"
	"fmt"
	"io"

	"github.com/davidterranova/contacts/internal/contacts/domain"
	"github.com/davidterranova/contacts/internal/contacts/exports"
	"github.com/davidterranova/contacts/pkg/user"
	"github.com/go-playground/validator"
	uuid "github.com/google/uuid"
)

type CmdExportContact struct {
	ContactId string `validate:"required,uuid"`
}

type ExportContactHandler struct {
	validator *validator.Validate
	lister    ContactReadModel
}

func NewExportContact(lister ContactReadModel) ExportContactHandler {
	return ExportContactHandler{
		validator: validator.New(),
		lister:    lister,
	}
}

func (h ExportContactHandler) Export(ctx context.Context, cmd CmdExportContact, cmdIssuedBy user.User) (io.Writer, error) {
	err := h.validator.Struct(cmd)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCommand, err)
	}

	uuid, err := uuid.Parse(cmd.ContactId)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCommand, err)
	}

	contact, err := h.lister.Get(ctx, QueryContact{
		ContactId: &uuid,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInternal, err)
	}

	err = checkExportPolicy(cmdIssuedBy, contact)
	if err != nil {
		return nil, err
	}

	return exports.VCardRenderer{}.Render(contact)
}

func checkExportPolicy(requestor user.User, aggregate *domain.Contact) error {
	if requestor.Id() == aggregate.CreatedBy.Id() {
		return nil
	}

	return ErrForbidden
}
