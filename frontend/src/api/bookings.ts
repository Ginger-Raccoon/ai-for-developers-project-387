import { apiFetch } from './client'
import type { Booking, CreateBookingRequest } from './types'

export function getBookings(): Promise<Booking[]> {
  return apiFetch('/api/v1/bookings')
}

export function createBooking(body: CreateBookingRequest): Promise<Booking> {
  return apiFetch('/api/v1/bookings', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}
