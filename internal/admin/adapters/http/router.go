package http

import (
	"context"

	"github.com/davidterranova/contacts/internal/admin/domain"
	"github.com/davidterranova/contacts/internal/admin/usecase"
	"github.com/davidterranova/contacts/pkg/xhttp"
	"github.com/gorilla/mux"
)

type App interface {
	ListEvents(ctx context.Context, query usecase.QueryListEvent) ([]*domain.Event, error)
}

func New(root *mux.Router, app App, authFn xhttp.AuthFn) *mux.Router {
	mountV1Events(root, authFn, app)

	return root
}

func mountV1Events(router *mux.Router, authFn xhttp.AuthFn, app App) {
	eventsHandler := NewEventHandler(app)
	v1 := router.PathPrefix("/v1/events").Subrouter()
	v1.Use(
		mux.CORSMethodMiddleware(v1),
	)

	if authFn != nil {
		v1.Use(xhttp.AuthMiddleware(authFn))
	}

	v1.HandleFunc("", eventsHandler.List).Methods("GET")
}
