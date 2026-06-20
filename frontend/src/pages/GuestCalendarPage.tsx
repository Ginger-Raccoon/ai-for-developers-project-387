import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Alert, Calendar, Spin, Typography, Button, Space, Empty, Badge } from 'antd'
import dayjs, { type Dayjs } from 'dayjs'
import { useEventType } from '@/hooks/useEventType'
import { useSlots } from '@/hooks/useSlots'
import type { Slot } from '@/api/types'

function formatTime(iso: string) {
  return dayjs(iso).format('HH:mm')
}

export function GuestCalendarPage() {
  const { eventTypeId = '' } = useParams()
  const navigate = useNavigate()
  const [selectedDate, setSelectedDate] = useState<string | null>(null)
  const [selectedSlot, setSelectedSlot] = useState<Slot | null>(null)

  const { data: eventType, loading: etLoading, error: etError } = useEventType(eventTypeId)
  const { data: slots, loading: slotsLoading } = useSlots(eventTypeId, selectedDate)

  const today = dayjs().startOf('day')
  const maxDate = today.add(14, 'day')

  function disabledDate(date: Dayjs) {
    return date.isBefore(today, 'day') || date.isAfter(maxDate, 'day')
  }

  function handleDateSelect(date: Dayjs) {
    const formatted = date.format('YYYY-MM-DD')
    setSelectedDate(formatted)
    setSelectedSlot(null)
  }

  function handleNext() {
    if (!selectedSlot || !eventType) return
    navigate(`/book/${eventTypeId}/confirm`, {
      state: { slot: selectedSlot, eventType },
    })
  }

  if (etLoading) return <Spin size="large" style={{ display: 'block', marginTop: 80 }} />
  if (etError) return <Alert type="error" message={etError} />
  if (!eventType) return null

  return (
    <>
      <Typography.Title level={2}>{eventType.title}</Typography.Title>
      <Typography.Paragraph type="secondary">
        Длительность: {eventType.duration} мин
      </Typography.Paragraph>

      <div style={{ display: 'flex', gap: 32, flexWrap: 'wrap', alignItems: 'flex-start' }}>
        <div style={{ flex: '0 0 320px' }}>
          <Calendar
            fullscreen={false}
            disabledDate={disabledDate}
            onSelect={handleDateSelect}
          />
        </div>

        <div style={{ flex: 1, minWidth: 200 }}>
          {!selectedDate && (
            <Typography.Text type="secondary">Выберите дату в календаре</Typography.Text>
          )}

          {selectedDate && slotsLoading && <Spin />}

          {selectedDate && !slotsLoading && slots.length === 0 && (
            <Empty description="Нет свободных слотов на этот день" />
          )}

          {selectedDate && !slotsLoading && slots.length > 0 && (
            <>
              <Typography.Text strong style={{ display: 'block', marginBottom: 12 }}>
                {dayjs(selectedDate).format('D MMMM YYYY')}
              </Typography.Text>
              <Space wrap>
                {slots.map((slot) => {
                  const isSelected = selectedSlot?.start === slot.start
                  return (
                    <Badge key={slot.start} dot={isSelected} color="blue">
                      <Button
                        type={isSelected ? 'primary' : 'default'}
                        onClick={() => setSelectedSlot(slot)}
                      >
                        {formatTime(slot.start)}
                      </Button>
                    </Badge>
                  )
                })}
              </Space>
            </>
          )}

          {selectedSlot && (
            <Button
              type="primary"
              size="large"
              style={{ marginTop: 24, display: 'block' }}
              onClick={handleNext}
            >
              Продолжить →
            </Button>
          )}
        </div>
      </div>
    </>
  )
}
