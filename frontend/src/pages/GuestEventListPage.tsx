import { Alert, Empty, Row, Col, Spin, Typography } from 'antd'
import { useEventTypes } from '@/hooks/useEventTypes'
import { EventTypeCard } from '@/components/EventTypeCard'

export function GuestEventListPage() {
  const { data, loading, error } = useEventTypes()

  if (loading) return <Spin size="large" style={{ display: 'block', marginTop: 80 }} />
  if (error) return <Alert type="error" message={error} />

  return (
    <>
      <Typography.Title level={2}>Выберите тип встречи</Typography.Title>
      {data.length === 0 ? (
        <Empty description="Нет доступных типов событий" />
      ) : (
        <Row gutter={[24, 24]}>
          {data.map((et) => (
            <Col key={et.id} xs={24} sm={12} md={8}>
              <EventTypeCard eventType={et} />
            </Col>
          ))}
        </Row>
      )}
    </>
  )
}
