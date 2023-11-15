package http

import "github.com/davidterranova/contacts/internal/admin/domain"

type Event struct {
	EventId          string `json:"eventId"`
	EventIssuesAt    string `json:"eventIssuesAt"`
	EventIssuedBy    string `json:"eventIssuedBy"`
	EventType        string `json:"eventType"`
	AggregateType    string `json:"aggregateType"`
	AggregateId      string `json:"aggregateId"`
	AggregateVersion int    `json:"aggregateVersion"`
	EventData        string `json:"eventData"`
	Published        bool   `json:"published"`
}

func fromDomainList(events []*domain.Event) []*Event {
	var list = make([]*Event, 0, len(events))
	for _, e := range events {
		list = append(list, fromDomain(e))
	}

	return list
}

func fromDomain(e *domain.Event) *Event {
	return &Event{
		EventId:          e.EventId.String(),
		EventIssuesAt:    e.EventIssuesAt.Format("2006-01-02T15:04:05Z"),
		EventIssuedBy:    e.EventIssuedBy.String(),
		EventType:        e.EventType,
		AggregateType:    string(e.AggregateType),
		AggregateId:      e.AggregateId.String(),
		AggregateVersion: e.AggregateVersion,
		EventData:        string(e.EventData),
		Published:        e.Published,
	}
}
