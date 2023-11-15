package http

import (
	"net/http"
	"strconv"

	"github.com/davidterranova/contacts/internal/admin/usecase"
	"github.com/davidterranova/contacts/pkg/eventsourcing"
	"github.com/davidterranova/contacts/pkg/xhttp"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	aggregateIdParam   = "aggregate_id"
	aggregateTypeParam = "aggregate_type"
)

type EventHandler struct {
	app App
}

func NewEventHandler(app App) *EventHandler {
	return &EventHandler{
		app: app,
	}
}

func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	var query usecase.QueryListEvent
	ctx := r.Context()

	aggregateIdParam := queryParam(r, aggregateIdParam)
	if aggregateIdParam != nil {
		aggregateId, err := uuid.Parse(*aggregateIdParam)
		if err != nil {
			log.Ctx(ctx).Warn().Err(err).Msg("admin_events:list failed to parse aggregate_id")
			xhttp.WriteError(ctx, w, http.StatusBadRequest, "failed to parse aggregate_id", err)
			return
		}
		query.AggregateId = &aggregateId
	}

	aggregateTypeParam := queryParam(r, aggregateTypeParam)
	if aggregateTypeParam != nil {
		aggregateType := eventsourcing.AggregateType(*aggregateTypeParam)
		query.AggregateType = &aggregateType
	}

	eventPublisjed := queryParam(r, "published")
	if eventPublisjed != nil {
		published, err := strconv.ParseBool(*eventPublisjed)
		if err != nil {
			log.Ctx(ctx).Warn().Err(err).Msg("admin_events:list failed to parse published")
			xhttp.WriteError(ctx, w, http.StatusBadRequest, "failed to parse published", err)
			return
		}
		query.Published = &published
	}

	events, err := h.app.ListEvents(ctx, query)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("admin_events:list failed to list events")
		xhttp.WriteError(ctx, w, http.StatusInternalServerError, "failed to list events", err)
		return
	}

	toReturnEvents := fromDomainList(events)
	xhttp.WriteObject(ctx, w, http.StatusOK, toReturnEvents)
}

func queryParam(r *http.Request, key string) *string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}

	return &value
}
