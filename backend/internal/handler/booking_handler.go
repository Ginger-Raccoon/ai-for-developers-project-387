package handler

import (
	"booking-service/internal/domain"
	"booking-service/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

// BookingHandler handles HTTP requests for bookings.
type BookingHandler struct {
	svc *service.BookingService
}

// NewBookingHandler creates a new BookingHandler.
func NewBookingHandler(svc *service.BookingService) *BookingHandler {
	return &BookingHandler{svc: svc}
}

// List handles GET /api/v1/bookings
func (h *BookingHandler) List(w http.ResponseWriter, r *http.Request) {
	bookings := h.svc.ListFuture()
	if bookings == nil {
		bookings = []domain.Booking{}
	}
	writeJSON(w, http.StatusOK, bookings)
}

// Create handles POST /api/v1/bookings
func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	// We need to parse the JSON manually to handle the time field.
	var raw struct {
		EventTypeID string `json:"eventTypeId"`
		GuestName   string `json:"guestName"`
		GuestEmail  string `json:"guestEmail"`
		Start       string `json:"start"`
	}

	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	// Enforce validation order: steps 1 & 2 (400) must fire before time parsing
	// so that a missing guestName is reported before a bad start format.
	if strings.TrimSpace(raw.GuestName) == "" {
		writeError(w, http.StatusBadRequest, "guestName must not be empty")
		return
	}
	if !strings.Contains(raw.GuestEmail, "@") || !strings.Contains(raw.GuestEmail, ".") {
		writeError(w, http.StatusBadRequest, "guestEmail must be a valid email address")
		return
	}

	start, err := time.Parse(time.RFC3339, raw.Start)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid start time format, expected RFC3339")
		return
	}

	req := domain.CreateBookingRequest{
		EventTypeID: raw.EventTypeID,
		GuestName:   raw.GuestName,
		GuestEmail:  raw.GuestEmail,
		Start:       start,
	}

	booking, err := h.svc.Create(req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrValidation):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, domain.ErrNotFound):
			writeError(w, http.StatusNotFound, "event type not found")
		case errors.Is(err, domain.ErrConflict):
			writeError(w, http.StatusConflict, err.Error())
		case errors.Is(err, domain.ErrUnprocessable):
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusCreated, booking)
}
