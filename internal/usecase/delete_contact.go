package usecase

import (
	"context"
	"fmt"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/go-playground/validator"
	uuid "github.com/google/uuid"
)

type CmdDeleteContact struct {
	ContactId string `validate:"required,uuid"`
}

type DeleteContactHandler struct {
	repo      ContactRepository
	validator *validator.Validate
}

func NewDeleteContact(repo ContactRepository) DeleteContactHandler {
	return DeleteContactHandler{
		repo:      repo,
		validator: validator.New(),
	}
}

func (h DeleteContactHandler) Delete(ctx context.Context, cmd CmdDeleteContact) error {
	err := h.validator.Struct(cmd)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidCommand, err)
	}

	contactUUID, err := uuid.Parse(cmd.ContactId)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidCommand, err)
	}

	_, err = handleRepositoryError[*domain.Contact](nil, h.repo.Delete(ctx, contactUUID))
	return err
}
