import { useState, useEffect } from 'react'
import { getBookings } from '@/api/bookings'
import type { Booking } from '@/api/types'

export function useBookings() {
  const [data, setData] = useState<Booking[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    getBookings()
      .then(setData)
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false))
  }, [])

  return { data, loading, error }
}
