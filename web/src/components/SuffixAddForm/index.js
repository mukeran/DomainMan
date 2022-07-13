import React, { useRef } from 'react'
import { Button, Input, Form } from 'antd';
import PropTypes from 'prop-types'

import api from '../../api';
import { connect } from 'react-redux';

const mapStateToProps = state => {
  const { auth } = state
  const { accessToken } = auth || {
    accessToken: ''
  }
  return {
    accessToken
  }
}

const SuffixAddForm = (props) => {
  const formRef = useRef(null)

  const { accessToken, onSubmitted, onError } = props

  const handleFinish = (values) => {
    api.suffix.add(accessToken, { values })
      .then(() => {
        formRef.current.resetFields()
        onSubmitted()
      })
      .catch(err => {
        onError(err)
      })
  }

  return (
    <Form ref={formRef} onFinish={handleFinish} layout="vertical">
      <Form.Item
        name="name"
        label="Suffix"
        validateTrigger={['onChange', 'onBlur']}
        rules={[
          {
            required: true,
            whitespace: true
          }
        ]}
      >
        <Input placeholder="Suffix"/>
      </Form.Item>
      <Form.Item name="memo" label="Memo">
        <Input.TextArea placeholder="Memo"/>
      </Form.Item>
      <Form.Item name="description" label="Description">
        <Input.TextArea placeholder="Description"/>
      </Form.Item>
      <Form.Item
        name="whoisServer"
        label="Whois Server"
        validateTrigger={['onChange', 'onBlur']}
        rules={[
          {
            required: true,
            whitespace: true
          }
        ]}
      >
        <Input placeholder="Whois Server"/>
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit" block>
          Submit
        </Button>
      </Form.Item>
    </Form>
  )
}

SuffixAddForm.propTypes = {
  onSubmitted: PropTypes.func.isRequired,
  onError: PropTypes.func
}

export default connect(mapStateToProps)(SuffixAddForm)
