import { createBrowserRouter } from 'react-router-dom'
import { AppLayout } from '@/components/AppLayout'
import { GuestEventListPage } from '@/pages/GuestEventListPage'
import { GuestCalendarPage } from '@/pages/GuestCalendarPage'
import { GuestBookingFormPage } from '@/pages/GuestBookingFormPage'
import { GuestConfirmationPage } from '@/pages/GuestConfirmationPage'
import { AdminDashboardPage } from '@/pages/AdminDashboardPage'
import { AdminCreateEventTypePage } from '@/pages/AdminCreateEventTypePage'

function wrap(element: React.ReactNode) {
  return <AppLayout>{element}</AppLayout>
}

export const router = createBrowserRouter([
  { path: '/', element: wrap(<GuestEventListPage />) },
  { path: '/book/:eventTypeId', element: wrap(<GuestCalendarPage />) },
  { path: '/book/:eventTypeId/confirm', element: wrap(<GuestBookingFormPage />) },
  { path: '/booking/success', element: wrap(<GuestConfirmationPage />) },
  { path: '/admin', element: wrap(<AdminDashboardPage />) },
  { path: '/admin/event-types/new', element: wrap(<AdminCreateEventTypePage />) },
])
