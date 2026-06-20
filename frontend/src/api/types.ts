export interface EventType {
  id: string
  title: string
  description: string
  duration: number
}

export interface Slot {
  start: string
  end: string
}

export interface Booking {
  id: string
  eventTypeId: string
  eventTypeTitle: string
  guestName: string
  guestEmail: string
  start: string
  end: string
  createdAt: string
}

export interface CreateEventTypeRequest {
  id: string
  title: string
  description: string
  duration: number
}

export interface CreateBookingRequest {
  eventTypeId: string
  guestName: string
  guestEmail: string
  start: string
}

export interface ApiError {
  message: string
}

export class ApiException extends Error {
  readonly status: number

  constructor(status: number, message: string) {
    super(message)
    this.name = 'ApiException'
    this.status = status
  }
}
