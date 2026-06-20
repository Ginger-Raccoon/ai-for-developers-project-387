package main_test

// Integration tests for the Call Booking API.
//
// Each test gets its own isolated in-memory server so tests are fully
// independent and can run in parallel.
//
// User scenarios covered:
//  1. Admin creates an event type; guest lists and reads it.
//  2. Guest views available slots for an event type (full window + date filter).
//  3. Guest creates a booking; admin sees it in the upcoming list.
//  4. Double-booking is rejected with 409.
//  5. Validation errors return the correct 4xx codes in the correct order.
//  6. CORS headers are present on every response; OPTIONS returns 204.

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"booking-service/internal/domain"
	"booking-service/internal/handler"
	"booking-service/internal/service"
	"booking-service/internal/store"
)

// Wire types match the JSON shape returned by the API (dates as RFC3339 strings).
type bookingWire struct {
	ID             string `json:"id"`
	EventTypeID    string `json:"eventTypeId"`
	EventTypeTitle string `json:"eventTypeTitle"`
	GuestName      string `json:"guestName"`
	GuestEmail     string `json:"guestEmail"`
	Start          string `json:"start"`
	End            string `json:"end"`
	CreatedAt      string `json:"createdAt"`
}

type slotWire struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// ── Test helpers ──────────────────────────────────────────────────────────────

// newTestServer wires up fresh in-memory stores and returns an isolated server.
// The server is automatically closed when the test ends.
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	ets := store.NewEventTypeStore()
	bs := store.NewBookingStore()

	etSvc := service.NewEventTypeService(ets)
	slSvc := service.NewSlotService(ets, bs)
	bSvc := service.NewBookingService(ets, bs)

	etH := handler.NewEventTypeHandler(etSvc)
	slH := handler.NewSlotHandler(slSvc)
	bH := handler.NewBookingHandler(bSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/event-types", etH.List)
	mux.HandleFunc("POST /api/v1/event-types", etH.Create)
	mux.HandleFunc("GET /api/v1/event-types/{id}", etH.Read)
	mux.HandleFunc("GET /api/v1/event-types/{eventTypeId}/slots", slH.List)
	mux.HandleFunc("GET /api/v1/bookings", bH.List)
	mux.HandleFunc("POST /api/v1/bookings", bH.Create)

	srv := httptest.NewServer(handler.CORSMiddleware(mux))
	t.Cleanup(srv.Close)
	return srv
}

func request(t *testing.T, srv *httptest.Server, method, path, body string) *http.Response {
	t.Helper()
	var r io.Reader
	if body != "" {
		r = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, srv.URL+path, r)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	return resp
}

func decode[T any](t *testing.T, resp *http.Response) T {
	t.Helper()
	defer resp.Body.Close()
	var v T
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	return v
}

func assertStatus(t *testing.T, resp *http.Response, want int) {
	t.Helper()
	if resp.StatusCode != want {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		t.Fatalf("expected HTTP %d, got %d: %s", want, resp.StatusCode, body)
	}
}

// tomorrowMidnightUTC returns tomorrow at 00:00 UTC formatted as RFC3339.
// Using midnight ensures the slot is always within the 14-day window,
// always in the future, and always grid-aligned for any slot duration.
func tomorrowMidnightUTC() string {
	now := time.Now().UTC()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
	return tomorrow.Format(time.RFC3339)
}

// tomorrowDate returns tomorrow's date as YYYY-MM-DD.
func tomorrowDate() string {
	return time.Now().UTC().AddDate(0, 0, 1).Format("2006-01-02")
}

const (
	quickCallBody = `{"id":"30min","title":"Quick call","description":"A 30-minute intro","duration":30}`
)

// ── Scenario 1: Admin creates event type; guest lists and reads it ─────────────

func TestEventTypes_CreateAndList(t *testing.T) {
	srv := newTestServer(t)

	// Initially empty
	resp := request(t, srv, "GET", "/api/v1/event-types", "")
	assertStatus(t, resp, http.StatusOK)
	list := decode[[]domain.EventType](t, resp)
	if len(list) != 0 {
		t.Fatalf("expected empty list, got %d items", len(list))
	}

	// Create event type
	resp = request(t, srv, "POST", "/api/v1/event-types", quickCallBody)
	assertStatus(t, resp, http.StatusCreated)
	et := decode[domain.EventType](t, resp)
	if et.ID != "30min" || et.Title != "Quick call" || et.Duration != 30 {
		t.Fatalf("unexpected event type in response: %+v", et)
	}

	// List returns the created type
	resp = request(t, srv, "GET", "/api/v1/event-types", "")
	assertStatus(t, resp, http.StatusOK)
	list = decode[[]domain.EventType](t, resp)
	if len(list) != 1 || list[0].ID != "30min" {
		t.Fatalf("expected exactly 1 event type, got: %+v", list)
	}
}

func TestEventTypes_ReadByID(t *testing.T) {
	srv := newTestServer(t)
	request(t, srv, "POST", "/api/v1/event-types", quickCallBody).Body.Close()

	resp := request(t, srv, "GET", "/api/v1/event-types/30min", "")
	assertStatus(t, resp, http.StatusOK)
	et := decode[domain.EventType](t, resp)
	if et.ID != "30min" || et.Duration != 30 {
		t.Fatalf("unexpected event type: %+v", et)
	}
}

// ── Scenario 2: Guest views available slots ───────────────────────────────────

func TestSlots_FullWindow(t *testing.T) {
	srv := newTestServer(t)
	request(t, srv, "POST", "/api/v1/event-types", quickCallBody).Body.Close()

	resp := request(t, srv, "GET", "/api/v1/event-types/30min/slots", "")
	assertStatus(t, resp, http.StatusOK)
	slots := decode[[]slotWire](t, resp)
	// 14 days × 48 slots/day = 672 for a 30-min event type
	if len(slots) == 0 {
		t.Fatal("expected non-empty slots over 14-day window")
	}
}

func TestSlots_DateFilter(t *testing.T) {
	srv := newTestServer(t)
	request(t, srv, "POST", "/api/v1/event-types", quickCallBody).Body.Close()

	resp := request(t, srv, "GET", "/api/v1/event-types/30min/slots?date="+tomorrowDate(), "")
	assertStatus(t, resp, http.StatusOK)
	slots := decode[[]slotWire](t, resp)
	// 24h / 30min = 48 slots for one day
	if len(slots) != 48 {
		t.Fatalf("expected 48 slots for tomorrow, got %d", len(slots))
	}
}

func TestSlots_SlotBecomesUnavailableAfterBooking(t *testing.T) {
	srv := newTestServer(t)
	request(t, srv, "POST", "/api/v1/event-types", quickCallBody).Body.Close()

	start := tomorrowMidnightUTC()

	// Count free slots before booking
	resp := request(t, srv, "GET", "/api/v1/event-types/30min/slots?date="+tomorrowDate(), "")
	assertStatus(t, resp, http.StatusOK)
	before := decode[[]slotWire](t, resp)

	// Book the first slot
	request(t, srv, "POST", "/api/v1/bookings",
		fmt.Sprintf(`{"eventTypeId":"30min","guestName":"Alice","guestEmail":"alice@example.com","start":%q}`, start),
	).Body.Close()

	// One fewer slot available
	resp = request(t, srv, "GET", "/api/v1/event-types/30min/slots?date="+tomorrowDate(), "")
	assertStatus(t, resp, http.StatusOK)
	after := decode[[]slotWire](t, resp)

	if len(after) != len(before)-1 {
		t.Fatalf("expected %d slots after booking, got %d", len(before)-1, len(after))
	}
}

// ── Scenario 3: Guest creates a booking; admin sees it ───────────────────────

func TestBooking_Create(t *testing.T) {
	srv := newTestServer(t)
	request(t, srv, "POST", "/api/v1/event-types", quickCallBody).Body.Close()

	start := tomorrowMidnightUTC()
	resp := request(t, srv, "POST", "/api/v1/bookings",
		fmt.Sprintf(`{"eventTypeId":"30min","guestName":"Ivan","guestEmail":"ivan@example.com","start":%q}`, start))
	assertStatus(t, resp, http.StatusCreated)

	b := decode[bookingWire](t, resp)
	if b.ID == "" {
		t.Fatal("expected non-empty booking ID")
	}
	if b.EventTypeID != "30min" || b.EventTypeTitle != "Quick call" {
		t.Fatalf("unexpected event type fields: %+v", b)
	}
	if b.GuestName != "Ivan" || b.GuestEmail != "ivan@example.com" {
		t.Fatalf("unexpected guest fields: %+v", b)
	}
	if b.Start == "" || b.End == "" || b.CreatedAt == "" {
		t.Fatalf("expected all time fields to be set: %+v", b)
	}
}

func TestBookings_ListShowsUpcoming(t *testing.T) {
	srv := newTestServer(t)
	request(t, srv, "POST", "/api/v1/event-types", quickCallBody).Body.Close()

	start := tomorrowMidnightUTC()
	request(t, srv, "POST", "/api/v1/bookings",
		fmt.Sprintf(`{"eventTypeId":"30min","guestName":"Ivan","guestEmail":"ivan@example.com","start":%q}`, start),
	).Body.Close()

	resp := request(t, srv, "GET", "/api/v1/bookings", "")
	assertStatus(t, resp, http.StatusOK)
	bookings := decode[[]bookingWire](t, resp)
	if len(bookings) != 1 || bookings[0].GuestName != "Ivan" {
		t.Fatalf("expected 1 booking for Ivan, got: %+v", bookings)
	}
}

func TestBookings_ListIsSortedAscending(t *testing.T) {
	srv := newTestServer(t)
	request(t, srv, "POST", "/api/v1/event-types", quickCallBody).Body.Close()

	now := time.Now().UTC()
	day1 := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
	day2 := day1.Add(30 * time.Minute)

	request(t, srv, "POST", "/api/v1/bookings",
		fmt.Sprintf(`{"eventTypeId":"30min","guestName":"Bob","guestEmail":"bob@example.com","start":%q}`, day2.Format(time.RFC3339)),
	).Body.Close()
	request(t, srv, "POST", "/api/v1/bookings",
		fmt.Sprintf(`{"eventTypeId":"30min","guestName":"Alice","guestEmail":"alice@example.com","start":%q}`, day1.Format(time.RFC3339)),
	).Body.Close()

	resp := request(t, srv, "GET", "/api/v1/bookings", "")
	assertStatus(t, resp, http.StatusOK)
	bookings := decode[[]bookingWire](t, resp)
	if len(bookings) != 2 {
		t.Fatalf("expected 2 bookings, got %d", len(bookings))
	}
	if bookings[0].GuestName != "Alice" || bookings[1].GuestName != "Bob" {
		t.Fatalf("expected ascending order Alice→Bob, got: %v, %v", bookings[0].GuestName, bookings[1].GuestName)
	}
}

// ── Scenario 4: Double booking is rejected ────────────────────────────────────

func TestBooking_Conflict(t *testing.T) {
	srv := newTestServer(t)
	request(t, srv, "POST", "/api/v1/event-types", quickCallBody).Body.Close()

	bookBody := fmt.Sprintf(
		`{"eventTypeId":"30min","guestName":"Ivan","guestEmail":"ivan@example.com","start":%q}`,
		tomorrowMidnightUTC(),
	)

	resp := request(t, srv, "POST", "/api/v1/bookings", bookBody)
	assertStatus(t, resp, http.StatusCreated)
	resp.Body.Close()

	// Second booking at the same slot → 409
	resp = request(t, srv, "POST", "/api/v1/bookings", bookBody)
	assertStatus(t, resp, http.StatusConflict)
	errBody := decode[domain.ApiError](t, resp)
	if errBody.Message == "" {
		t.Fatal("expected non-empty error message in 409 response")
	}
}


// ── Validation: event types ───────────────────────────────────────────────────

func TestEventType_DuplicateID_Returns409(t *testing.T) {
	srv := newTestServer(t)
	request(t, srv, "POST", "/api/v1/event-types", quickCallBody).Body.Close()

	resp := request(t, srv, "POST", "/api/v1/event-types", quickCallBody)
	assertStatus(t, resp, http.StatusConflict)
	resp.Body.Close()
}

func TestEventType_EmptyTitle_Returns400(t *testing.T) {
	srv := newTestServer(t)
	resp := request(t, srv, "POST", "/api/v1/event-types",
		`{"id":"test","title":"","description":"","duration":30}`)
	assertStatus(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestEventType_ZeroDuration_Returns400(t *testing.T) {
	srv := newTestServer(t)
	resp := request(t, srv, "POST", "/api/v1/event-types",
		`{"id":"test","title":"Test","description":"","duration":0}`)
	assertStatus(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestEventType_NotFound_Returns404(t *testing.T) {
	srv := newTestServer(t)
	resp := request(t, srv, "GET", "/api/v1/event-types/nonexistent", "")
	assertStatus(t, resp, http.StatusNotFound)
	resp.Body.Close()
}

func TestSlots_UnknownEventType_Returns404(t *testing.T) {
	srv := newTestServer(t)
	resp := request(t, srv, "GET", "/api/v1/event-types/ghost/slots", "")
	assertStatus(t, resp, http.StatusNotFound)
	resp.Body.Close()
}

// ── Validation: bookings — strict order matters ───────────────────────────────

func TestBooking_EmptyName_Returns400(t *testing.T) {
	srv := newTestServer(t)
	request(t, srv, "POST", "/api/v1/event-types", quickCallBody).Body.Close()

	resp := request(t, srv, "POST", "/api/v1/bookings",
		fmt.Sprintf(`{"eventTypeId":"30min","guestName":"","guestEmail":"a@b.com","start":%q}`,
			tomorrowMidnightUTC()))
	assertStatus(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestBooking_InvalidEmail_Returns400(t *testing.T) {
	srv := newTestServer(t)
	request(t, srv, "POST", "/api/v1/event-types", quickCallBody).Body.Close()

	resp := request(t, srv, "POST", "/api/v1/bookings",
		fmt.Sprintf(`{"eventTypeId":"30min","guestName":"Ivan","guestEmail":"notanemail","start":%q}`,
			tomorrowMidnightUTC()))
	assertStatus(t, resp, http.StatusBadRequest)
	resp.Body.Close()
}

func TestBooking_UnknownEventType_Returns404(t *testing.T) {
	srv := newTestServer(t)
	resp := request(t, srv, "POST", "/api/v1/bookings",
		fmt.Sprintf(`{"eventTypeId":"ghost","guestName":"Ivan","guestEmail":"a@b.com","start":%q}`,
			tomorrowMidnightUTC()))
	assertStatus(t, resp, http.StatusNotFound)
	resp.Body.Close()
}

// ValidationOrder: empty name → 400 even when eventTypeId doesn't exist (name check runs first)
func TestBooking_ValidationOrder_NameBeforeEventType(t *testing.T) {
	srv := newTestServer(t)
	resp := request(t, srv, "POST", "/api/v1/bookings",
		fmt.Sprintf(`{"eventTypeId":"ghost","guestName":"","guestEmail":"a@b.com","start":%q}`,
			tomorrowMidnightUTC()))
	assertStatus(t, resp, http.StatusBadRequest) // 400, not 404
	resp.Body.Close()
}

func TestBooking_SlotOutsideWindow_Returns422(t *testing.T) {
	srv := newTestServer(t)
	request(t, srv, "POST", "/api/v1/event-types", quickCallBody).Body.Close()

	farFuture := time.Now().UTC().AddDate(0, 0, 30).Truncate(30 * time.Minute).Format(time.RFC3339)
	resp := request(t, srv, "POST", "/api/v1/bookings",
		fmt.Sprintf(`{"eventTypeId":"30min","guestName":"Ivan","guestEmail":"a@b.com","start":%q}`, farFuture))
	assertStatus(t, resp, http.StatusUnprocessableEntity)
	resp.Body.Close()
}

func TestBooking_SlotNotOnGrid_Returns422(t *testing.T) {
	srv := newTestServer(t)
	request(t, srv, "POST", "/api/v1/event-types", quickCallBody).Body.Close()

	// tomorrow at 00:15 — not aligned to 30-min grid
	now := time.Now().UTC()
	offGrid := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 15, 0, 0, time.UTC).Format(time.RFC3339)
	resp := request(t, srv, "POST", "/api/v1/bookings",
		fmt.Sprintf(`{"eventTypeId":"30min","guestName":"Ivan","guestEmail":"a@b.com","start":%q}`, offGrid))
	assertStatus(t, resp, http.StatusUnprocessableEntity)
	resp.Body.Close()
}

// ── CORS ──────────────────────────────────────────────────────────────────────

func TestCORS_Options_Returns204(t *testing.T) {
	srv := newTestServer(t)
	resp := request(t, srv, "OPTIONS", "/api/v1/event-types", "")
	defer resp.Body.Close()
	assertStatus(t, resp, http.StatusNoContent)
	if got := resp.Header.Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("expected CORS origin *, got %q", got)
	}
}

func TestCORS_HeadersPresentOnAllResponses(t *testing.T) {
	srv := newTestServer(t)
	for _, path := range []string{
		"/api/v1/event-types",
		"/api/v1/bookings",
	} {
		resp := request(t, srv, "GET", path, "")
		defer resp.Body.Close()
		if got := resp.Header.Get("Access-Control-Allow-Origin"); got != "*" {
			t.Fatalf("%s: expected CORS header *, got %q", path, got)
		}
	}
}

// ── JSON response shape ───────────────────────────────────────────────────────

func TestEventTypes_EmptyListIsArray_NotNull(t *testing.T) {
	srv := newTestServer(t)
	resp := request(t, srv, "GET", "/api/v1/event-types", "")
	assertStatus(t, resp, http.StatusOK)
	defer resp.Body.Close()
	var raw json.RawMessage
	json.NewDecoder(resp.Body).Decode(&raw) //nolint:errcheck
	if string(raw) == "null" {
		t.Fatal("empty list must be [] not null")
	}
}

func TestBookings_EmptyListIsArray_NotNull(t *testing.T) {
	srv := newTestServer(t)
	resp := request(t, srv, "GET", "/api/v1/bookings", "")
	assertStatus(t, resp, http.StatusOK)
	defer resp.Body.Close()
	var raw json.RawMessage
	json.NewDecoder(resp.Body).Decode(&raw) //nolint:errcheck
	if string(raw) == "null" {
		t.Fatal("empty list must be [] not null")
	}
}
