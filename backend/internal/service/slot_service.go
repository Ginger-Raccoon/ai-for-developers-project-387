package service

import (
	"booking-service/internal/domain"
	"booking-service/internal/store"
	"time"
)

// SlotService handles business logic for generating free time slots.
type SlotService struct {
	eventTypeStore *store.EventTypeStore
	bookingStore   *store.BookingStore
}

// NewSlotService creates a new SlotService.
func NewSlotService(ets *store.EventTypeStore, bs *store.BookingStore) *SlotService {
	return &SlotService{
		eventTypeStore: ets,
		bookingStore:   bs,
	}
}

// ListSlots returns free slots for the given eventTypeId.
// If date is zero, returns slots for the next 14 days.
// If date is provided, returns slots only for that day.
func (s *SlotService) ListSlots(eventTypeID string, date time.Time) ([]domain.Slot, error) {
	et, err := s.eventTypeStore.Get(eventTypeID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	var windowStart, windowEnd time.Time

	if date.IsZero() {
		// No date filter: entire 14-day window
		windowStart = today
		windowEnd = today.AddDate(0, 0, 14)
	} else {
		// Specific date: just that day
		d := date.UTC()
		windowStart = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
		windowEnd = windowStart.AddDate(0, 0, 1)
	}

	duration := time.Duration(et.Duration) * time.Minute
	bookings := s.bookingStore.All()

	var slots []domain.Slot
	current := windowStart
	for !current.Add(duration).After(windowEnd) {
		slotStart := current
		slotEnd := current.Add(duration)

		if isFree(slotStart, slotEnd, bookings) {
			slots = append(slots, domain.Slot{
				Start: slotStart,
				End:   slotEnd,
			})
		}
		current = current.Add(duration)
	}

	if slots == nil {
		slots = []domain.Slot{}
	}
	return slots, nil
}

// isFree returns true if no booking overlaps [slotStart, slotEnd).
func isFree(slotStart, slotEnd time.Time, bookings []domain.Booking) bool {
	for _, b := range bookings {
		// overlap: b.Start < slotEnd && b.End > slotStart
		if b.Start.Before(slotEnd) && b.End.After(slotStart) {
			return false
		}
	}
	return true
}
