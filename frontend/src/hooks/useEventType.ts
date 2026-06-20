import { useState, useEffect } from 'react'
import { getEventType } from '@/api/eventTypes'
import type { EventType } from '@/api/types'

export function useEventType(id: string) {
  const [data, setData] = useState<EventType | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    setLoading(true)
    setError(null)
    getEventType(id)
      .then(setData)
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false))
  }, [id])

  return { data, loading, error }
}
