package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"

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

	// Serve frontend static files if the ./static directory exists.
	if info, err := os.Stat("./static"); err == nil && info.IsDir() {
		staticFS := os.DirFS("./static")
		fileServer := http.FileServer(http.FS(staticFS))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Try to serve the file; fall back to index.html for SPA routing.
			if _, err := fs.Stat(staticFS, r.URL.Path[1:]); err == nil && r.URL.Path != "/" {
				fileServer.ServeHTTP(w, r)
				return
			}
			http.ServeFileFS(w, r, staticFS, "index.html")
		})
	}

	// CORSMiddleware wraps the mux; it intercepts OPTIONS preflights before
	// they reach any registered route handler.
	corsHandler := handler.CORSMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("Booking service listening on %s", addr)
	if err := http.ListenAndServe(addr, corsHandler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
