package usecase

import (
	"context"
	"fmt"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/go-playground/validator"
	uuid "github.com/google/uuid"
)

type CmdCreateContact struct {
	CreatedBy string `validate:"required,uuid"`

	FirstName string `validate:"min=2,max=255"`
	LastName  string `validate:"min=2,max=255"`
	Email     string `validate:"required,email"`
	Phone     string `validate:"e164"` // https://en.wikipedia.org/wiki/E.164
}

type CreateContact struct {
	repo      ContactRepository
	validator *validator.Validate
}

func NewCreateContact(repo ContactRepository) CreateContact {
	return CreateContact{
		repo:      repo,
		validator: validator.New(),
	}
}

func (h CreateContact) Create(ctx context.Context, cmd CmdCreateContact) (*domain.Contact, error) {
	err := h.validator.Struct(cmd)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCommand, err)
	}

	createdBy, err := uuid.Parse(cmd.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCommand, err)
	}

	contact := domain.New(createdBy)
	contact.FirstName = cmd.FirstName
	contact.LastName = cmd.LastName
	contact.Email = cmd.Email
	contact.Phone = cmd.Phone

	return handleRepositoryError(h.repo.Save(ctx, contact))
}
