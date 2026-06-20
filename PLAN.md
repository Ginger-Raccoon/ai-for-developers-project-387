# Development Plan

## Bugs (Priority 1)

### B-1. Weak email validation
**File:** `backend/internal/handler/booking_handler.go:53`  
Current check passes `"@test"`, `"user@"`, `"a@.b"`. Replace with a proper regex.

### B-2. Race condition on concurrent bookings
**File:** `backend/internal/service/booking_service.go:70`  
`HasOverlap` check and record creation are not atomic: two simultaneous requests can both pass the check and create a duplicate booking.

### B-3. ListFuture skips bookings scheduled for exactly "now"
**File:** `backend/internal/store/booking_store.go:31`  
`b.Start.After(now)` excludes a slot whose start time equals the current moment. Use `!b.Start.Before(now)` instead.

---

## New Features (Priority 2)

### F-1. Booking cancellation
No `DELETE /api/v1/bookings/{id}` in API or UI. The service cannot be considered complete without it.

### F-2. Health check endpoint
No `/health` route — required for Render, Docker health checks, and future Kubernetes deployments.

### F-3. Auto-refresh of bookings list in admin panel
`useBookings` fetches data once on mount. New bookings don't appear without a full page reload.  
**File:** `frontend/src/hooks/useBookings.ts`

---

## Tech Debt (Priority 3)

### T-1. Duplicated validation logic
Email and name validation is spread across handler, service, and frontend in three different variations.

### T-2. Ignored JSON encoding errors
`//nolint:errcheck` in `backend/internal/handler/middleware.go:29` silently swallows encoding failures.

### T-3. No structured logging on the backend
In production there is zero visibility into what the server is doing.
