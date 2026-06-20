import { useState, useEffect } from 'react'
import { getEventTypes } from '@/api/eventTypes'
import type { EventType } from '@/api/types'

export function useEventTypes() {
  const [data, setData] = useState<EventType[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    getEventTypes()
      .then(setData)
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false))
  }, [])

  return { data, loading, error }
}
