import React, { useEffect, useState } from 'react'
import { Layout } from 'antd'
import { TransitionGroup, CSSTransition } from 'react-transition-group'
import { useLocation, Outlet } from 'react-router-dom'

import './style.less'
import Sidebar from '../../components/Sidebar'
import Footer from '../../components/Footer'
import { connect } from 'react-redux'
import api from '../../api'
import constants from '../../constants'
import { message } from 'antd'
import { useNavigate } from 'react-router-dom'

const mapStateToProps = state => {
  const { auth } = state
  const { accessToken } = auth || {
    accessToken: ''
  }
  return {
    accessToken
  }
}

const App = (props) => {
  const location = useLocation()
  const navigate = useNavigate()

  const [collapsed, setCollapsed] = useState(false)

  const onCollapse = value => {
    setCollapsed(value)
  }

  const { accessToken } = props
  useEffect(() => {
    api.system.ping(accessToken)
      .catch(err => {
        if (err.status === constants.status.connectionError) {
          message.error('Failed to connect to server')
        } else {
          message.error('Invalid access token')
          navigate('/auth')
        }
      })
  })

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Layout.Sider collapsible collapsed={collapsed} onCollapse={onCollapse}>
        <div className="logo" />
        <Sidebar/>
      </Layout.Sider>
      <Layout>
        <Layout.Header style={{
          padding: 0,
          background: '#fff',
        }}>
        </Layout.Header>
        <Layout.Content className="content">
          <TransitionGroup>
            <CSSTransition
              timeout={500}
              classNames="fade"
              key={location.pathname}
              mountOnEnter={false}
              unmountOnExit={true}
            >
              <Outlet />
            </CSSTransition>
          </TransitionGroup>
        </Layout.Content>
        <Layout.Footer><Footer/></Layout.Footer>
      </Layout>
    </Layout>
  )
}

export default connect(mapStateToProps)(App)
