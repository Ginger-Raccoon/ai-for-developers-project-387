package domain

import (
	"errors"
	"time"
)

// Sentinel errors
var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrValidation    = errors.New("validation error")
	ErrConflict      = errors.New("conflict")
	ErrUnprocessable = errors.New("unprocessable")
)

// EventType represents a type of meeting/call
type EventType struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"` // minutes
}

// CreateEventTypeRequest is the input for creating an EventType
type CreateEventTypeRequest struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
}

// Slot represents a free time slot
type Slot struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Booking represents a confirmed booking
type Booking struct {
	ID             string    `json:"id"`
	EventTypeID    string    `json:"eventTypeId"`
	EventTypeTitle string    `json:"eventTypeTitle"`
	GuestName      string    `json:"guestName"`
	GuestEmail     string    `json:"guestEmail"`
	Start          time.Time `json:"start"`
	End            time.Time `json:"end"`
	CreatedAt      time.Time `json:"createdAt"`
}

// CreateBookingRequest is the input for creating a Booking
type CreateBookingRequest struct {
	EventTypeID string    `json:"eventTypeId"`
	GuestName   string    `json:"guestName"`
	GuestEmail  string    `json:"guestEmail"`
	Start       time.Time `json:"start"`
}

// ApiError is the standard error response body
type ApiError struct {
	Message string `json:"message"`
}
