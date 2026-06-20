package main

import (
	"log"
	"net/http"

	"booking-service/internal/handler"
	"booking-service/internal/service"
	"booking-service/internal/store"
)

func main() {
	// Stores (in-memory)
	eventTypeStore := store.NewEventTypeStore()
	bookingStore := store.NewBookingStore()

	// Services
	eventTypeSvc := service.NewEventTypeService(eventTypeStore)
	slotSvc := service.NewSlotService(eventTypeStore, bookingStore)
	bookingSvc := service.NewBookingService(eventTypeStore, bookingStore)

	// Handlers
	eventTypeHandler := handler.NewEventTypeHandler(eventTypeSvc)
	slotHandler := handler.NewSlotHandler(slotSvc)
	bookingHandler := handler.NewBookingHandler(bookingSvc)

	// Router (Go 1.22+ pattern matching with method prefix)
	mux := http.NewServeMux()

	// EventType routes
	mux.HandleFunc("GET /api/v1/event-types", eventTypeHandler.List)
	mux.HandleFunc("POST /api/v1/event-types", eventTypeHandler.Create)
	mux.HandleFunc("GET /api/v1/event-types/{id}", eventTypeHandler.Read)

	// Slot routes
	mux.HandleFunc("GET /api/v1/event-types/{eventTypeId}/slots", slotHandler.List)

	// Booking routes
	mux.HandleFunc("GET /api/v1/bookings", bookingHandler.List)
	mux.HandleFunc("POST /api/v1/bookings", bookingHandler.Create)

	// CORSMiddleware wraps the mux; it intercepts OPTIONS preflights before
	// they reach any registered route handler.
	corsHandler := handler.CORSMiddleware(mux)

	addr := ":8080"
	log.Printf("Booking service listening on %s", addr)
	if err := http.ListenAndServe(addr, corsHandler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
