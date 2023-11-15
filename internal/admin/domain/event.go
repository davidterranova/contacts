package domain

import (
	"encoding/json"
	"time"

	"github.com/davidterranova/contacts/pkg/user"

	"github.com/davidterranova/contacts/pkg/eventsourcing"
	"github.com/google/uuid"
)

type Event struct {
	EventId          uuid.UUID
	EventIssuesAt    time.Time
	EventIssuedBy    user.User
	EventType        string
	AggregateType    eventsourcing.AggregateType
	AggregateId      uuid.UUID
	AggregateVersion int
	EventData        json.RawMessage

	Published bool
}
