import React, { useEffect, useRef } from 'react'
import { PlusOutlined, MinusCircleOutlined } from '@ant-design/icons'
import { Button, Input, Form } from 'antd';
import PropTypes from 'prop-types'

import constants from '../../constants'
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

const DomainAddForm = (props) => {
  const formRef = useRef(null)

  const reset = () => {
    formRef.current.setFieldsValue({
      domains: ['']
    })
  }

  const { accessToken, onSubmitted, onError } = props

  const handleFinish = (values) => {
    let { domains } = values
    api.domain.add(accessToken, { domains })
      .then(() => {
        reset()
        onSubmitted()
      })
      .catch(err => {
        onError(err)
      })
  }

  useEffect(() => {
    reset()
  })

  return (
    <Form ref={formRef} onFinish={handleFinish} layout="vertical">
      <Form.List name="domains">
        {(fields, { add, remove }) => {
          return (
            <div>
              {fields.map((field, index) => (
                <Form.Item
                  {...field}
                  label={index === 0 ? 'Domains': ''}
                  validateTrigger={['onChange', 'onBlur']}
                  rules={[
                    {
                      required: true,
                      whitespace: true,
                      pattern: constants.regexp.domain
                    }
                  ]}
                >
                  <Input placeholder="Domain Name" suffix={
                    fields.length > 1 ? (
                      <MinusCircleOutlined onClick={() => remove(field.name)}/>
                    ) : null
                  }/>
                </Form.Item>
              ))}
              <Form.Item>
                <Button type="dashed" onClick={() => add()} block>
                  <PlusOutlined /> Add One
                </Button>
              </Form.Item>
              <Form.Item>
                <Button type="primary" htmlType="submit" block>
                  Submit
                </Button>
              </Form.Item>
            </div>
          )
        }}
      </Form.List>
    </Form>
  );
}

DomainAddForm.propTypes = {
  onSubmitted: PropTypes.func.isRequired,
  onError: PropTypes.func
}

export default connect(mapStateToProps)(DomainAddForm)
