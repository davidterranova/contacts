//go:generate mockgen -destination=mock_app.go -package=http . App
package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/internal/usecase"
	"github.com/davidterranova/contacts/pkg/xhttp"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type App interface {
	ListContacts(ctx context.Context, query usecase.QueryListContact) ([]*domain.Contact, error)
	CreateContact(ctx context.Context, cmd usecase.CmdCreateContact) (*domain.Contact, error)
	UpdateContact(ctx context.Context, cmd usecase.CmdUpdateContact) (*domain.Contact, error)
	DeleteContact(ctx context.Context, cmd usecase.CmdDeleteContact) error
}

type ContactHandler struct {
	app App
}

func NewContactHandler(app App) *ContactHandler {
	return &ContactHandler{
		app: app,
	}
}

func (h *ContactHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contacts, err := h.app.ListContacts(ctx, usecase.QueryListContact{})
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:list failed to list contacts")
		xhttp.WriteError(ctx, w, http.StatusInternalServerError, "failed to list contacts", err)
		return
	}

	toReturnContacts := fromDomainList(contacts)
	xhttp.WriteObject(ctx, w, http.StatusOK, toReturnContacts)
}

type createContactRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required"`
}

func (h *ContactHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createContactRequest
	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:create failed to decode request")
		xhttp.WriteError(ctx, w, http.StatusBadRequest, "failed to decode request", err)
		return
	}

	contact, err := h.app.CreateContact(
		ctx,
		usecase.NewCmdCreateContact(
			req.FirstName,
			req.LastName,
			req.Email,
			req.Phone,
		),
	)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidCommand):
			xhttp.WriteError(ctx, w, http.StatusBadRequest, "contact validation failed", err)
		default:
			log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:update failed to update contact")
			xhttp.WriteError(ctx, w, http.StatusInternalServerError, "failed to update contact", err)
		}
		return
	}

	toReturnContact := fromDomain(contact)
	xhttp.WriteObject(ctx, w, http.StatusCreated, toReturnContact)
}

type updateContactRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

func (h *ContactHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req updateContactRequest
	ctx := r.Context()
	contactId, err := uuid.Parse(mux.Vars(r)[pathContactId])
	if err != nil {
		xhttp.WriteError(ctx, w, http.StatusBadRequest, "invalid resource uuid", err)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:update failed to decode request")
		xhttp.WriteError(ctx, w, http.StatusBadRequest, "failed to decode request", err)
		return
	}

	contact, err := h.app.UpdateContact(
		ctx,
		usecase.NewCmdUpdateContact(
			contactId,
			req.FirstName,
			req.LastName,
			req.Email,
			req.Phone,
		),
	)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidCommand):
			xhttp.WriteError(ctx, w, http.StatusBadRequest, "contact validation failed", err)
		case errors.Is(err, usecase.ErrNotFound):
			xhttp.WriteError(ctx, w, http.StatusNotFound, "contact not found", err)
		default:
			log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:update failed to update contact")
			xhttp.WriteError(ctx, w, http.StatusInternalServerError, "failed to update contact", err)
		}
		return
	}

	toReturnContact := fromDomain(contact)
	xhttp.WriteObject(ctx, w, http.StatusOK, toReturnContact)
}

func (h *ContactHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	contactId, err := uuid.Parse(mux.Vars(r)[pathContactId])
	if err != nil {
		xhttp.WriteError(ctx, w, http.StatusBadRequest, "invalid resource uuid", err)
		return
	}

	err = h.app.DeleteContact(ctx, usecase.NewCmdDeleteContact(contactId))
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidCommand):
			xhttp.WriteError(ctx, w, http.StatusBadRequest, "contact validation failed", err)
		case errors.Is(err, usecase.ErrNotFound):
			xhttp.WriteError(ctx, w, http.StatusNotFound, "contact not found", err)
		default:
			log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:update failed to update contact")
			xhttp.WriteError(ctx, w, http.StatusInternalServerError, "failed to update contact", err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
