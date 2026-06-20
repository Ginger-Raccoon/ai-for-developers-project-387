import { Layout, Menu } from 'antd'
import { Link, useLocation } from 'react-router-dom'

const { Header, Content } = Layout

export function AppLayout({ children }: { children: React.ReactNode }) {
  const location = useLocation()

  const selectedKey = location.pathname.startsWith('/admin') ? 'admin' : 'guest'

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ display: 'flex', alignItems: 'center', gap: 32 }}>
        <span style={{ color: '#fff', fontWeight: 700, fontSize: 18, whiteSpace: 'nowrap' }}>
          Запись на звонок
        </span>
        <Menu
          theme="dark"
          mode="horizontal"
          selectedKeys={[selectedKey]}
          style={{ flex: 1, minWidth: 0 }}
          items={[
            { key: 'guest', label: <Link to="/">Записаться</Link> },
            { key: 'admin', label: <Link to="/admin">Для организатора</Link> },
          ]}
        />
      </Header>
      <Content style={{ padding: '32px 48px', maxWidth: 960, margin: '0 auto', width: '100%' }}>
        {children}
      </Content>
    </Layout>
  )
}
