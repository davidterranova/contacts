//go:generate mockgen -destination=mock_contact_repository.go -package=usecase . ContactRepository
package usecase

import (
	"context"

	"github.com/davidterranova/contacts/internal/domain"
	_ "github.com/golang/mock/mockgen/model"
	uuid "github.com/google/uuid"
)

type ContactRepository interface {
	List(ctx context.Context) ([]*domain.Contact, error)
	Save(ctx context.Context, contact *domain.Contact) (*domain.Contact, error)
	Update(ctx context.Context, id uuid.UUID, updateFn func(c domain.Contact) (domain.Contact, error)) (*domain.Contact, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
