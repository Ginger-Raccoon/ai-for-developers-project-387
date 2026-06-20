import { Alert, Button, Spin, Table, Typography } from 'antd'
import { useNavigate } from 'react-router-dom'
import dayjs from 'dayjs'
import { useBookings } from '@/hooks/useBookings'
import type { Booking } from '@/api/types'
import type { ColumnsType } from 'antd/es/table'

const columns: ColumnsType<Booking> = [
  {
    title: 'Тип встречи',
    dataIndex: 'eventTypeTitle',
    key: 'eventTypeTitle',
  },
  {
    title: 'Дата и время',
    key: 'start',
    render: (_, r) =>
      `${dayjs(r.start).format('D MMM YYYY, HH:mm')} — ${dayjs(r.end).format('HH:mm')}`,
  },
  {
    title: 'Гость',
    dataIndex: 'guestName',
    key: 'guestName',
  },
  {
    title: 'Email',
    dataIndex: 'guestEmail',
    key: 'guestEmail',
  },
]

export function AdminDashboardPage() {
  const { data, loading, error } = useBookings()
  const navigate = useNavigate()

  if (loading) return <Spin size="large" style={{ display: 'block', marginTop: 80 }} />
  if (error) return <Alert type="error" message={error} />

  return (
    <>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <Typography.Title level={2} style={{ margin: 0 }}>
          Предстоящие встречи
        </Typography.Title>
        <Button type="primary" onClick={() => navigate('/admin/event-types/new')}>
          + Создать тип события
        </Button>
      </div>

      <Table
        rowKey="id"
        columns={columns}
        dataSource={data}
        locale={{ emptyText: 'Нет предстоящих встреч' }}
        pagination={false}
      />
    </>
  )
}
