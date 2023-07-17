package ports

import (
	"context"
	"errors"
	"fmt"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("contact not found")

// InMemoryContactRepository is not thread safe
type InMemoryContactRepository struct {
	contacts map[uuid.UUID]*domain.Contact
}

func NewInMemoryContactRepository() *InMemoryContactRepository {
	return &InMemoryContactRepository{
		contacts: map[uuid.UUID]*domain.Contact{},
	}
}

func (r *InMemoryContactRepository) Get(_ context.Context, id uuid.UUID) (*domain.Contact, error) {
	contact, ok := r.contacts[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, id)
	}

	return contact, nil
}

func (r *InMemoryContactRepository) List(ctx context.Context) ([]*domain.Contact, error) {
	contacts := make([]*domain.Contact, 0, len(r.contacts))
	for _, contact := range r.contacts {
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

func (r *InMemoryContactRepository) Save(_ context.Context, contact *domain.Contact) (*domain.Contact, error) {
	r.contacts[contact.Id] = contact
	return contact, nil
}

func (r *InMemoryContactRepository) Update(ctx context.Context, id uuid.UUID, updateFn func(c domain.Contact) (domain.Contact, error)) (*domain.Contact, error) {
	originalConact, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	updatedContact, err := updateFn(*originalConact)
	if err != nil {
		return nil, err
	}

	return r.Save(ctx, &updatedContact)
}

func (r *InMemoryContactRepository) Delete(_ context.Context, id uuid.UUID) error {
	if _, ok := r.contacts[id]; !ok {
		return fmt.Errorf("%w: %s", ErrNotFound, id)
	}
	delete(r.contacts, id)

	return nil
}
