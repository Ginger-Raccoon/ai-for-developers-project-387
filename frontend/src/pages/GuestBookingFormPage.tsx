import { useEffect, useState } from 'react'
import { useLocation, useNavigate, useParams } from 'react-router-dom'
import { Form, Input, Button, Typography, Descriptions, Alert, message } from 'antd'
import dayjs from 'dayjs'
import { createBooking } from '@/api/bookings'
import { ApiException } from '@/api/types'
import type { Slot, EventType } from '@/api/types'

interface LocationState {
  slot: Slot
  eventType: EventType
}

export function GuestBookingFormPage() {
  const { eventTypeId = '' } = useParams()
  const navigate = useNavigate()
  const location = useLocation()
  const state = location.state as LocationState | null

  const [submitting, setSubmitting] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => {
    if (!state?.slot || !state?.eventType) {
      navigate(`/book/${eventTypeId}`, { replace: true })
    }
  }, [state, eventTypeId, navigate])

  if (!state?.slot || !state?.eventType) return null

  const { slot, eventType } = state

  async function handleSubmit(values: { guestName: string; guestEmail: string }) {
    setSubmitting(true)
    try {
      const booking = await createBooking({
        eventTypeId,
        guestName: values.guestName,
        guestEmail: values.guestEmail,
        start: slot.start,
      })
      navigate('/booking/success', { state: booking })
    } catch (e) {
      if (e instanceof ApiException && e.status === 409) {
        message.error('Этот слот уже занят. Выберите другое время.')
      } else if (e instanceof ApiException && e.status === 400) {
        message.error(e.message)
      } else {
        message.error('Произошла ошибка. Попробуйте ещё раз.')
      }
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <>
      <Typography.Title level={2}>Подтверждение записи</Typography.Title>

      <Descriptions bordered column={1} style={{ marginBottom: 32, maxWidth: 480 }}>
        <Descriptions.Item label="Встреча">{eventType.title}</Descriptions.Item>
        <Descriptions.Item label="Дата и время">
          {dayjs(slot.start).format('D MMMM YYYY, HH:mm')} — {dayjs(slot.end).format('HH:mm')}
        </Descriptions.Item>
        <Descriptions.Item label="Длительность">{eventType.duration} мин</Descriptions.Item>
      </Descriptions>

      <Alert
        type="info"
        message="Укажите ваши данные для бронирования"
        style={{ marginBottom: 24, maxWidth: 480 }}
      />

      <Form form={form} layout="vertical" onFinish={handleSubmit} style={{ maxWidth: 480 }}>
        <Form.Item
          name="guestName"
          label="Ваше имя"
          rules={[{ required: true, message: 'Введите ваше имя' }]}
        >
          <Input placeholder="Иван Иванов" />
        </Form.Item>

        <Form.Item
          name="guestEmail"
          label="Email"
          rules={[
            { required: true, message: 'Введите email' },
            { type: 'email', message: 'Введите корректный email' },
          ]}
        >
          <Input placeholder="ivan@example.com" />
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" loading={submitting} size="large">
            Записаться
          </Button>
        </Form.Item>
      </Form>
    </>
  )
}
