package domain

import (
	"encoding/json"
	"time"
)

// slotJSON is the wire representation of a Slot.
type slotJSON struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// MarshalJSON implements json.Marshaler for Slot, using RFC3339 UTC.
func (s Slot) MarshalJSON() ([]byte, error) {
	return json.Marshal(slotJSON{
		Start: s.Start.UTC().Format(time.RFC3339),
		End:   s.End.UTC().Format(time.RFC3339),
	})
}

// bookingJSON is the wire representation of a Booking.
type bookingJSON struct {
	ID             string `json:"id"`
	EventTypeID    string `json:"eventTypeId"`
	EventTypeTitle string `json:"eventTypeTitle"`
	GuestName      string `json:"guestName"`
	GuestEmail     string `json:"guestEmail"`
	Start          string `json:"start"`
	End            string `json:"end"`
	CreatedAt      string `json:"createdAt"`
}

// MarshalJSON implements json.Marshaler for Booking, using RFC3339 UTC.
func (b Booking) MarshalJSON() ([]byte, error) {
	return json.Marshal(bookingJSON{
		ID:             b.ID,
		EventTypeID:    b.EventTypeID,
		EventTypeTitle: b.EventTypeTitle,
		GuestName:      b.GuestName,
		GuestEmail:     b.GuestEmail,
		Start:          b.Start.UTC().Format(time.RFC3339),
		End:            b.End.UTC().Format(time.RFC3339),
		CreatedAt:      b.CreatedAt.UTC().Format(time.RFC3339),
	})
}
