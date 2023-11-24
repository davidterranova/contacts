package http

import (
	"net/http"

	"github.com/davidterranova/contacts/pkg/xhttp"
	"github.com/gorilla/mux"
)

const pathContactId = "contactId"

// New returns a new contacts API router
func New(root *mux.Router, app App, authFn xhttp.AuthFn) *mux.Router {
	mountV1Contacts(root, authFn, app)

	return root
}

func mountV1Contacts(root *mux.Router, authFn xhttp.AuthFn, app App) {
	contactsHandler := NewContactHandler(app)
	v1 := root.PathPrefix("/v1/contacts").Subrouter()
	v1.Use(
		mux.CORSMethodMiddleware(v1),
	)

	if authFn != nil {
		v1.Use(xhttp.AuthMiddleware(authFn))
	}

	v1.HandleFunc("", contactsHandler.List).Methods(http.MethodGet)
	v1.HandleFunc("", contactsHandler.Create).Methods(http.MethodPost)
	v1.HandleFunc("/{"+pathContactId+"}", contactsHandler.Update).Methods(http.MethodPut)
	v1.HandleFunc("/{"+pathContactId+"}", contactsHandler.Delete).Methods(http.MethodDelete)
	v1.HandleFunc("/{"+pathContactId+"}:export", contactsHandler.Export).Methods(http.MethodGet)
}
