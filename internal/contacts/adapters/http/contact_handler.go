//go:generate mockgen -destination=mock_app.go -package=http . App
package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/davidterranova/contacts/internal/contacts/domain"
	"github.com/davidterranova/contacts/internal/contacts/usecase"
	"github.com/davidterranova/contacts/pkg/auth"
	"github.com/davidterranova/contacts/pkg/user"
	"github.com/davidterranova/contacts/pkg/xhttp"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type App interface {
	ListContacts(ctx context.Context, cmdIssuedBy user.User) ([]*domain.Contact, error)
	ExportContact(ctx context.Context, cmd usecase.CmdExportContact, cmdIssuedBy user.User) (io.Writer, error)
	CreateContact(ctx context.Context, cmd usecase.CmdCreateContact, cmdIssuedBy user.User) (*domain.Contact, error)
	UpdateContact(ctx context.Context, cmd usecase.CmdUpdateContact, cmdIssuedBy user.User) (*domain.Contact, error)
	DeleteContact(ctx context.Context, cmd usecase.CmdDeleteContact, cmdIssuedBy user.User) error
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
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:list failed to get user from context")
		xhttp.WriteError(ctx, w, http.StatusInternalServerError, "failed to get user from context", err)
		return
	}

	contacts, err := h.app.ListContacts(ctx, user)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:list failed to list contacts")
		xhttp.WriteError(ctx, w, http.StatusInternalServerError, "failed to list contacts", err)
		return
	}

	toReturnContacts := fromDomainList(contacts)
	xhttp.WriteObject(ctx, w, http.StatusOK, toReturnContacts)
}

func (h *ContactHandler) Export(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:export failed to get user from context")
		xhttp.WriteError(ctx, w, http.StatusInternalServerError, "failed to get user from context", err)
		return
	}

	contactId := mux.Vars(r)[pathContactId]
	export, err := h.app.ExportContact(ctx, usecase.CmdExportContact{ContactId: contactId}, user)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidCommand):
			xhttp.WriteError(ctx, w, http.StatusBadRequest, "command validation failed", err)
		case errors.Is(err, usecase.ErrNotFound):
			xhttp.WriteError(ctx, w, http.StatusNotFound, "contact not found", err)
		case errors.Is(err, usecase.ErrForbidden):
			xhttp.WriteError(ctx, w, http.StatusForbidden, "access denied to contact", err)
		default:
			log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:export failed to export contact")
			xhttp.WriteError(ctx, w, http.StatusInternalServerError, "failed to export contact", err)
		}
		return
	}

	w.Header().Set("Content-Type", "text/vcard")
	w.Header().Set("Content-Disposition", "attachment;filename=contacts.vcf")
	w.WriteHeader(http.StatusOK)

	_, err = fmt.Fprint(w, export)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:export failed to write response")
	}
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
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:create failed to get user from context")
		xhttp.WriteError(ctx, w, http.StatusInternalServerError, "failed to get user from context", err)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:create failed to decode request")
		xhttp.WriteError(ctx, w, http.StatusBadRequest, "failed to decode request", err)
		return
	}

	contact, err := h.app.CreateContact(
		ctx,
		usecase.CmdCreateContact{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Phone:     req.Phone,
		},
		user,
	)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidCommand):
			xhttp.WriteError(ctx, w, http.StatusBadRequest, "contact validation failed", err)
		default:
			log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:create failed to create contact")
			xhttp.WriteError(ctx, w, http.StatusInternalServerError, "failed to create contact", err)
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
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:update failed to get user from context")
		xhttp.WriteError(ctx, w, http.StatusInternalServerError, "failed to get user from context", err)
		return
	}
	contactId := mux.Vars(r)[pathContactId]

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:update failed to decode request")
		xhttp.WriteError(ctx, w, http.StatusBadRequest, "failed to decode request", err)
		return
	}

	contact, err := h.app.UpdateContact(
		ctx,
		usecase.CmdUpdateContact{
			ContactId: contactId,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Phone:     req.Phone,
		},
		user,
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
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:delete failed to get user from context")
		xhttp.WriteError(ctx, w, http.StatusInternalServerError, "failed to get user from context", err)
		return
	}

	contactId := mux.Vars(r)[pathContactId]

	err = h.app.DeleteContact(ctx, usecase.CmdDeleteContact{ContactId: contactId}, user)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidCommand):
			xhttp.WriteError(ctx, w, http.StatusBadRequest, "contact validation failed", err)
		case errors.Is(err, usecase.ErrNotFound):
			xhttp.WriteError(ctx, w, http.StatusNotFound, "contact not found", err)
		default:
			log.Ctx(ctx).Warn().Err(err).Msg("user_contacts:delete failed to update contact")
			xhttp.WriteError(ctx, w, http.StatusInternalServerError, "failed to update contact", err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
