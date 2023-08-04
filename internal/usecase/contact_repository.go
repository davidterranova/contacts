//go:generate mockgen -destination=mock_contact_repository.go -package=usecase . ContactRepository
package usecase

import (
	"context"

	"github.com/davidterranova/contacts/internal/domain"
	_ "github.com/golang/mock/mockgen/model"
)

type ContactLister interface {
	List(ctx context.Context) ([]*domain.Contact, error)
}
