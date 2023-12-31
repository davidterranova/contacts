package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/davidterranova/contacts/pkg/user"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/go-playground/validator"
	uuid "github.com/google/uuid"
)

type CmdUpdateContact struct {
	Updater   user.User `validate:"required"`
	ContactId string    `validate:"required,uuid"`
	FirstName string    `validate:"omitempty,min=2,max=255"`
	LastName  string    `validate:"omitempty,min=2,max=255"`
	Email     string    `validate:"omitempty,email"`
	Phone     string    `validate:"omitempty,e164"` // https://en.wikipedia.org/wiki/E.164
}

type UpdateContact struct {
	repo      ContactRepository
	validator *validator.Validate
}

func NewUpdateContact(repo ContactRepository) UpdateContact {
	return UpdateContact{
		repo:      repo,
		validator: validator.New(),
	}
}

func (h UpdateContact) Update(ctx context.Context, cmd CmdUpdateContact) (*domain.Contact, error) {
	err := h.validator.Struct(cmd)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCommand, err)
	}

	contactUUID, err := uuid.Parse(cmd.ContactId)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCommand, err)
	}

	contact, err := h.repo.Update(ctx, contactUUID, func(c domain.Contact) (domain.Contact, error) {
		return updateContactFn(c, cmd)
	})

	return handleRepositoryError(contact, err)
}

func updateContactFn(c domain.Contact, cmd CmdUpdateContact) (domain.Contact, error) {
	updated := false

	if c.CreatedBy != cmd.Updater.Id() {
		return c, fmt.Errorf("%w: %s", ErrForbidden, "contact can only be updated by its creator")
	}

	if cmd.FirstName != "" {
		updated = true
		c.FirstName = cmd.FirstName
	}

	if cmd.LastName != "" {
		updated = true
		c.LastName = cmd.LastName
	}

	if cmd.Email != "" {
		updated = true
		c.Email = cmd.Email
	}

	if cmd.Phone != "" {
		updated = true
		c.Phone = cmd.Phone
	}

	if updated {
		c.UpdatedAt = time.Now().UTC()
	}

	return c, nil
}
