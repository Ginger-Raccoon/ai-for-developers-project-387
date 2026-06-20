import { useState, useEffect } from 'react'
import { getSlots } from '@/api/eventTypes'
import type { Slot } from '@/api/types'

export function useSlots(eventTypeId: string, date: string | null) {
  const [data, setData] = useState<Slot[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!date) {
      setData([])
      return
    }
    setLoading(true)
    setError(null)
    getSlots(eventTypeId, date)
      .then(setData)
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false))
  }, [eventTypeId, date])

  return { data, loading, error }
}
