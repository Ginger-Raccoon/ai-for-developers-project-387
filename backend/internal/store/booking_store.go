package store

import (
	"booking-service/internal/domain"
	"sort"
	"sync"
	"time"
)

// BookingStore is an in-memory store for Booking records.
type BookingStore struct {
	mu    sync.RWMutex
	items map[string]domain.Booking
}

// NewBookingStore creates a new BookingStore.
func NewBookingStore() *BookingStore {
	return &BookingStore{
		items: make(map[string]domain.Booking),
	}
}

// ListFuture returns all bookings with start >= now, sorted ascending by start.
func (s *BookingStore) ListFuture() []domain.Booking {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now().UTC()
	var result []domain.Booking
	for _, b := range s.items {
		if b.Start.After(now) {
			result = append(result, b)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Start.Before(result[j].Start)
	})
	return result
}

// All returns all bookings (for slot overlap checks).
func (s *BookingStore) All() []domain.Booking {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]domain.Booking, 0, len(s.items))
	for _, b := range s.items {
		result = append(result, b)
	}
	return result
}

// Create adds a new Booking.
func (s *BookingStore) Create(b domain.Booking) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[b.ID] = b
}

// HasOverlap checks whether any existing booking overlaps [start, end).
// Returns the conflicting booking if found.
func (s *BookingStore) HasOverlap(start, end time.Time) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, b := range s.items {
		// overlap: b.Start < end && b.End > start
		if b.Start.Before(end) && b.End.After(start) {
			return true
		}
	}
	return false
}
