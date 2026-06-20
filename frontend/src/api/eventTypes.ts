import { apiFetch } from './client'
import type { EventType, Slot, CreateEventTypeRequest } from './types'

export function getEventTypes(): Promise<EventType[]> {
  return apiFetch('/api/v1/event-types')
}

export function getEventType(id: string): Promise<EventType> {
  return apiFetch(`/api/v1/event-types/${id}`)
}

export function getSlots(eventTypeId: string, date?: string): Promise<Slot[]> {
  const params = date ? `?date=${date}` : ''
  return apiFetch(`/api/v1/event-types/${eventTypeId}/slots${params}`)
}

export function createEventType(body: CreateEventTypeRequest): Promise<EventType> {
  return apiFetch('/api/v1/event-types', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}
