package handler

import (
	"booking-service/internal/domain"
	"booking-service/internal/service"
	"encoding/json"
	"errors"
	"net/http"
)

// EventTypeHandler handles HTTP requests for event types.
type EventTypeHandler struct {
	svc *service.EventTypeService
}

// NewEventTypeHandler creates a new EventTypeHandler.
func NewEventTypeHandler(svc *service.EventTypeService) *EventTypeHandler {
	return &EventTypeHandler{svc: svc}
}

// List handles GET /api/v1/event-types
func (h *EventTypeHandler) List(w http.ResponseWriter, r *http.Request) {
	items := h.svc.List()
	// Always return an array, never null
	if items == nil {
		items = []domain.EventType{}
	}
	writeJSON(w, http.StatusOK, items)
}

// Create handles POST /api/v1/event-types
func (h *EventTypeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateEventTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	et, err := h.svc.Create(req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrValidation):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, domain.ErrAlreadyExists):
			writeError(w, http.StatusConflict, "event type with this id already exists")
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusCreated, et)
}

// Read handles GET /api/v1/event-types/{id}
func (h *EventTypeHandler) Read(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	et, err := h.svc.Get(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "event type not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, et)
}
