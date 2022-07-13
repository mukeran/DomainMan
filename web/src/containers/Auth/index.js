import { KeyOutlined } from '@ant-design/icons'
import { Row } from 'antd'
import { Card } from 'antd'
import { Button } from 'antd'
import { Input } from 'antd'
import { Col } from 'antd'
import { message } from 'antd'
import { Typography } from 'antd'
import React, { useRef, useState } from 'react'
import { useDispatch } from 'react-redux'
import { connect } from 'react-redux'
import { useNavigate } from 'react-router-dom'
import { updateAccessToken } from '../../actions/auth'
import api from '../../api'
import constants from '../../constants'

import './style.less'

const mapStateToProps = state => {
  const { auth } = state
  const { accessToken } = auth || {
    accessToken: ''
  }
  return {
    accessToken
  }
}

const Auth = (props) => {
  const dispatch = useDispatch()
  const navigate = useNavigate()
  const [isChecking, setIsChecking] = useState(false)
  const inputEl = useRef(null)

  const handleAuthClick = () => {
    const accessToken = inputEl.current.input.value
    setIsChecking(true)
    api.system.ping(accessToken)
      .then(() => {
        dispatch(updateAccessToken(accessToken))
        navigate('/')
      })
      .catch(err => {
        if (err.status === constants.status.connectionError) {
          message.error('Failed to connect to server')
        } else {
          message.error('Invalid access token')
        }
      })
      .finally(() => {
        setIsChecking(false)
      })
  }

  return (
    <Row>
      <Col sm={{ span: 20, offset: 2 }} md={{ span: 16, offset: 4 }} lg={{ span: 12, offset: 6 }} xl={{ span: 8, offset: 8 }}>
        <Typography.Title className="auth-title">DomainMan Web UI</Typography.Title>
        <Card title="Using Access Token to access Web API">
          <Input size="large" placeholder="Access Token" prefix={<KeyOutlined />} defaultValue={props.accessToken} ref={inputEl} disabled={isChecking} />
          <Button className="auth-button" type="primary" size="large" block onClick={handleAuthClick} disabled={isChecking}>Auth</Button>
        </Card>
      </Col>
    </Row>
  )
}

export default connect(mapStateToProps)(Auth)