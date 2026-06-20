import { RouterProvider } from 'react-router-dom'
import { ConfigProvider, App as AntApp } from 'antd'
import ruRU from 'antd/locale/ru_RU'
import dayjs from 'dayjs'
import 'dayjs/locale/ru'
import { router } from './router'

dayjs.locale('ru')

export function App() {
  return (
    <ConfigProvider locale={ruRU}>
      <AntApp>
        <RouterProvider router={router} />
      </AntApp>
    </ConfigProvider>
  )
}
