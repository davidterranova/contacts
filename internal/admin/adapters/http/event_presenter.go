package http

import "github.com/davidterranova/contacts/internal/admin/domain"

type Event struct {
	EventId          string `json:"event_id"`
	EventIssuesAt    string `json:"event_issues_at"`
	EventIssuedBy    string `json:"event_issued_by"`
	EventType        string `json:"event_type"`
	AggregateType    string `json:"aggregate_type"`
	AggregateId      string `json:"aggregate_id"`
	AggregateVersion int    `json:"aggregate_version"`
	EventData        string `json:"event_data"`
	Published        bool   `json:"event_published"`
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
