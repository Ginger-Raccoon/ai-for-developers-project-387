import { useEffect } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { Result, Button, Descriptions } from 'antd'
import dayjs from 'dayjs'
import type { Booking } from '@/api/types'

export function GuestConfirmationPage() {
  const location = useLocation()
  const navigate = useNavigate()
  const booking = location.state as Booking | null

  useEffect(() => {
    if (!booking) navigate('/', { replace: true })
  }, [booking, navigate])

  if (!booking) return null

  return (
    <Result
      status="success"
      title="Вы записаны!"
      subTitle={`Встреча "${booking.eventTypeTitle}" подтверждена`}
      extra={[
        <Button type="primary" onClick={() => navigate('/')} key="home">
          На главную
        </Button>,
      ]}
    >
      <Descriptions bordered column={1} style={{ maxWidth: 480, margin: '0 auto' }}>
        <Descriptions.Item label="Встреча">{booking.eventTypeTitle}</Descriptions.Item>
        <Descriptions.Item label="Дата и время">
          {dayjs(booking.start).format('D MMMM YYYY, HH:mm')} —{' '}
          {dayjs(booking.end).format('HH:mm')}
        </Descriptions.Item>
        <Descriptions.Item label="Имя">{booking.guestName}</Descriptions.Item>
        <Descriptions.Item label="Email">{booking.guestEmail}</Descriptions.Item>
      </Descriptions>
    </Result>
  )
}
