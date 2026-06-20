package store

import (
	"booking-service/internal/domain"
	"sync"
)

// EventTypeStore is an in-memory store for EventType records.
type EventTypeStore struct {
	mu    sync.RWMutex
	items map[string]domain.EventType
	order []string // preserve insertion order for listing
}

// NewEventTypeStore creates a new EventTypeStore.
func NewEventTypeStore() *EventTypeStore {
	return &EventTypeStore{
		items: make(map[string]domain.EventType),
	}
}

// List returns all EventTypes in insertion order.
func (s *EventTypeStore) List() []domain.EventType {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]domain.EventType, 0, len(s.order))
	for _, id := range s.order {
		result = append(result, s.items[id])
	}
	return result
}

// Get returns the EventType with the given ID, or ErrNotFound.
func (s *EventTypeStore) Get(id string) (domain.EventType, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	et, ok := s.items[id]
	if !ok {
		return domain.EventType{}, domain.ErrNotFound
	}
	return et, nil
}

// Create adds a new EventType. Returns ErrAlreadyExists if id is taken.
func (s *EventTypeStore) Create(et domain.EventType) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[et.ID]; exists {
		return domain.ErrAlreadyExists
	}
	s.items[et.ID] = et
	s.order = append(s.order, et.ID)
	return nil
}
