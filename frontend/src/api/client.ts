import { ApiException } from './types'

const BASE_URL = ''

export async function apiFetch<T>(path: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${BASE_URL}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
  })

  if (!response.ok) {
    let message = `HTTP ${response.status}`
    try {
      const body = await response.json()
      if (body?.message) message = body.message
    } catch {
      // ignore parse error, keep default message
    }
    throw new ApiException(response.status, message)
  }

  if (response.status === 204) return undefined as T
  return response.json() as Promise<T>
}
