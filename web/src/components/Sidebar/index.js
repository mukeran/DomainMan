import React, { useLayoutEffect, useState } from 'react'
import { Menu } from 'antd'
import { Link } from 'react-router-dom'
import {
  DashboardOutlined,
  GlobalOutlined,
  NumberOutlined,
  QuestionOutlined,
  OrderedListOutlined,
  SearchOutlined
} from '@ant-design/icons'

import { useLocation } from 'react-router-dom'

const Sidebar = (props) => {
  const [selectedKeys, setSelectedKeys] = useState([])
  const [openKeys, setOpenKeys] = useState([])
  const handleOpenChange = openKeys => {
    setOpenKeys(openKeys)
  }

  const location = useLocation()
  useLayoutEffect(() => {
    const subMenu = {
      '/whois': ['whoisSubMenu'],
      '/whois/query': ['whoisSubMenu']
    }
    setSelectedKeys([location.pathname])
    setOpenKeys(openKeys => Array.from(new Set([...openKeys].concat(subMenu[location.pathname] || []))))
  }, [location.pathname])

  const items = [
    {
      key: '/',
      icon: <DashboardOutlined />,
      label: <Link to="/">Dashboard</Link>
    },
    {
      key: '/domain',
      icon: <GlobalOutlined />,
      label: <Link to="/domain">Domain</Link>
    },
    {
      key: '/suffix',
      icon: <NumberOutlined />,
      label: <Link to="/suffix">Suffix</Link>
    },
    {
      key: 'whoisSubMenu',
      icon: <QuestionOutlined />,
      label: 'Whois',
      children: [
        {
          key: '/whois',
          icon: <OrderedListOutlined />,
          label: <Link to="/whois">List</Link>
        },
        {
          key: '/whois/query',
          icon: <SearchOutlined />,
          label: <Link to="/whois/query">Query</Link>
        }
      ]
    }
  ]

  return (
    <Menu
      theme="dark"
      mode="inline"
      selectedKeys={selectedKeys}
      openKeys={openKeys}
      onOpenChange={handleOpenChange}
      items={items}
    />
  )
}

export default Sidebar