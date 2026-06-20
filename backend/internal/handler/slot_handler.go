package handler

import (
	"booking-service/internal/service"
	"errors"
	"net/http"
	"time"

	"booking-service/internal/domain"
)

// SlotHandler handles HTTP requests for slots.
type SlotHandler struct {
	svc *service.SlotService
}

// NewSlotHandler creates a new SlotHandler.
func NewSlotHandler(svc *service.SlotService) *SlotHandler {
	return &SlotHandler{svc: svc}
}

// List handles GET /api/v1/event-types/{eventTypeId}/slots
func (h *SlotHandler) List(w http.ResponseWriter, r *http.Request) {
	eventTypeID := r.PathValue("eventTypeId")

	var date time.Time
	if dateStr := r.URL.Query().Get("date"); dateStr != "" {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid date format, expected YYYY-MM-DD")
			return
		}
		date = parsed
	}

	slots, err := h.svc.ListSlots(eventTypeID, date)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "event type not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, slots)
}
