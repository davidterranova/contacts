package domain

import (
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	Id uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time

	CreatedBy uuid.UUID

	FirstName string
	LastName  string
	Email     string
	Phone     string
}

func New(createdBy uuid.UUID) *Contact {
	now := time.Now().UTC()

	return &Contact{
		Id:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: createdBy,
	}
}
