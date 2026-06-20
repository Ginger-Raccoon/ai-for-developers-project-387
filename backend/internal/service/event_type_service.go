package service

import (
	"booking-service/internal/domain"
	"booking-service/internal/store"
	"fmt"
)

// EventTypeService handles business logic for EventTypes.
type EventTypeService struct {
	store *store.EventTypeStore
}

// NewEventTypeService creates a new EventTypeService.
func NewEventTypeService(s *store.EventTypeStore) *EventTypeService {
	return &EventTypeService{store: s}
}

// List returns all event types.
func (s *EventTypeService) List() []domain.EventType {
	return s.store.List()
}

// Get returns the event type with the given ID, or ErrNotFound.
func (s *EventTypeService) Get(id string) (domain.EventType, error) {
	return s.store.Get(id)
}

// Create validates and creates a new EventType.
// Returns:
//   - ErrValidation if title is empty or duration <= 0
//   - ErrAlreadyExists if id is taken
func (s *EventTypeService) Create(req domain.CreateEventTypeRequest) (domain.EventType, error) {
	if req.Title == "" {
		return domain.EventType{}, fmt.Errorf("%w: title must not be empty", domain.ErrValidation)
	}
	if req.Duration <= 0 {
		return domain.EventType{}, fmt.Errorf("%w: duration must be greater than 0", domain.ErrValidation)
	}

	et := domain.EventType{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
	}

	if err := s.store.Create(et); err != nil {
		return domain.EventType{}, err
	}
	return et, nil
}
