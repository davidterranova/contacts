package domain

import "github.com/google/uuid"

type Filter interface {
	CreatedBy() *uuid.UUID
}
