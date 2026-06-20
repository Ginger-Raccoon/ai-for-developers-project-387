package service

import (
	"booking-service/internal/domain"
	"booking-service/internal/store"
	"crypto/rand"
	"fmt"
	"strings"
	"time"
)

// BookingService handles business logic for bookings.
type BookingService struct {
	eventTypeStore *store.EventTypeStore
	bookingStore   *store.BookingStore
}

// NewBookingService creates a new BookingService.
func NewBookingService(ets *store.EventTypeStore, bs *store.BookingStore) *BookingService {
	return &BookingService{
		eventTypeStore: ets,
		bookingStore:   bs,
	}
}

// ListFuture returns all future bookings sorted ascending by start.
func (s *BookingService) ListFuture() []domain.Booking {
	return s.bookingStore.ListFuture()
}

// Create validates and creates a new Booking, strictly in the required order.
func (s *BookingService) Create(req domain.CreateBookingRequest) (domain.Booking, error) {
	// Step 1: guestName not empty → 400
	if strings.TrimSpace(req.GuestName) == "" {
		return domain.Booking{}, fmt.Errorf("%w: guestName must not be empty", domain.ErrValidation)
	}

	// Step 2: guestEmail contains '@' and '.' → 400
	if !strings.Contains(req.GuestEmail, "@") || !strings.Contains(req.GuestEmail, ".") {
		return domain.Booking{}, fmt.Errorf("%w: guestEmail must be a valid email address", domain.ErrValidation)
	}

	// Step 3: eventTypeId exists → 404
	et, err := s.eventTypeStore.Get(req.EventTypeID)
	if err != nil {
		return domain.Booking{}, err // ErrNotFound
	}

	// Step 4: validate 14-day window
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	windowStart := today
	windowEnd := today.AddDate(0, 0, 14)

	duration := time.Duration(et.Duration) * time.Minute
	start := req.Start.UTC()
	end := start.Add(duration)

	if start.Before(windowStart) || end.After(windowEnd) {
		return domain.Booking{}, fmt.Errorf("%w: slot is outside the 14-day booking window", domain.ErrUnprocessable)
	}

	// Step 5: alignment check — (start - windowStart) must be divisible by duration
	diff := start.Sub(windowStart)
	if int(diff.Minutes())%et.Duration != 0 {
		return domain.Booking{}, fmt.Errorf("%w: slot does not align with valid slot grid", domain.ErrUnprocessable)
	}

	// Step 6: no overlap with existing bookings → 409
	if s.bookingStore.HasOverlap(start, end) {
		return domain.Booking{}, fmt.Errorf("%w: time slot is already taken", domain.ErrConflict)
	}

	// Step 7: create booking
	id, err := newUUID()
	if err != nil {
		return domain.Booking{}, fmt.Errorf("failed to generate booking ID: %w", err)
	}

	b := domain.Booking{
		ID:             id,
		EventTypeID:    et.ID,
		EventTypeTitle: et.Title,
		GuestName:      req.GuestName,
		GuestEmail:     req.GuestEmail,
		Start:          start,
		End:            end,
		CreatedAt:      time.Now().UTC(),
	}
	s.bookingStore.Create(b)
	return b, nil
}

// newUUID generates a random UUID v4 using crypto/rand.
func newUUID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	// Set version 4
	b[6] = (b[6] & 0x0f) | 0x40
	// Set variant bits
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%12x",
		b[0:4],
		b[4:6],
		b[6:8],
		b[8:10],
		b[10:],
	), nil
}
