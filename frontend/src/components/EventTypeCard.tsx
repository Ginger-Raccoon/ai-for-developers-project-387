import { Card, Tag, Button } from 'antd'
import { ClockCircleOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import type { EventType } from '@/api/types'

export function EventTypeCard({ eventType }: { eventType: EventType }) {
  const navigate = useNavigate()

  return (
    <Card
      title={eventType.title}
      extra={
        <Tag icon={<ClockCircleOutlined />} color="blue">
          {eventType.duration} мин
        </Tag>
      }
      actions={[
        <Button type="primary" onClick={() => navigate(`/book/${eventType.id}`)}>
          Записаться
        </Button>,
      ]}
    >
      <p style={{ color: '#666', margin: 0 }}>{eventType.description || '—'}</p>
    </Card>
  )
}
