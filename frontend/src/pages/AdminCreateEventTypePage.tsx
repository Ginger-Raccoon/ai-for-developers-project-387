import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Form, Input, InputNumber, Button, Typography, message, Alert } from 'antd'
import { createEventType } from '@/api/eventTypes'
import { ApiException } from '@/api/types'

export function AdminCreateEventTypePage() {
  const navigate = useNavigate()
  const [submitting, setSubmitting] = useState(false)
  const [form] = Form.useForm()

  async function handleSubmit(values: {
    id: string
    title: string
    description: string
    duration: number
  }) {
    setSubmitting(true)
    try {
      await createEventType(values)
      message.success('Тип события создан')
      navigate('/admin')
    } catch (e) {
      if (e instanceof ApiException && e.status === 409) {
        form.setFields([{ name: 'id', errors: ['Тип события с таким ID уже существует'] }])
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
      <Typography.Title level={2}>Новый тип события</Typography.Title>

      <Alert
        type="info"
        message='ID используется в URL, например "30min" или "quick-chat"'
        style={{ marginBottom: 24, maxWidth: 520 }}
      />

      <Form form={form} layout="vertical" onFinish={handleSubmit} style={{ maxWidth: 520 }}>
        <Form.Item
          name="id"
          label="ID (slug)"
          rules={[
            { required: true, message: 'Введите ID' },
            {
              pattern: /^[a-z0-9-]+$/,
              message: 'Только строчные буквы, цифры и дефис',
            },
          ]}
        >
          <Input placeholder="30min" />
        </Form.Item>

        <Form.Item
          name="title"
          label="Название"
          rules={[{ required: true, message: 'Введите название' }]}
        >
          <Input placeholder="Быстрый звонок" />
        </Form.Item>

        <Form.Item name="description" label="Описание">
          <Input.TextArea rows={3} placeholder="Краткое описание встречи" />
        </Form.Item>

        <Form.Item
          name="duration"
          label="Длительность (минут)"
          rules={[{ required: true, message: 'Введите длительность' }]}
        >
          <InputNumber min={1} max={480} style={{ width: '100%' }} placeholder="30" />
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" loading={submitting} size="large">
            Создать
          </Button>
          <Button style={{ marginLeft: 12 }} onClick={() => navigate('/admin')}>
            Отмена
          </Button>
        </Form.Item>
      </Form>
    </>
  )
}
