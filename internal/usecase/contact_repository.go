//go:generate mockgen -destination=mock_contact_repository.go -package=usecase . ContactRepository
package usecase

import (
	"context"

	"github.com/davidterranova/contacts/internal/domain"
	_ "github.com/golang/mock/mockgen/model"
	uuid "github.com/google/uuid"
)

type ContactRepository interface {
	List(ctx context.Context, filter domain.Filter) ([]*domain.Contact, error)
	Create(ctx context.Context, contact *domain.Contact) (*domain.Contact, error)

	// Update with "Repository pattern" making a clean separation of concerns
	// between the use case and the persistence layer
	// while pushing atomicity to the persistence layer
	Update(ctx context.Context, id uuid.UUID, updateFn func(c domain.Contact) (domain.Contact, error)) (*domain.Contact, error)
	Delete(ctx context.Context, id uuid.UUID, deleterFn func(c domain.Contact) error) error
}
