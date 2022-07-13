import React from 'react'
import PropTypes from 'prop-types'

import { PageHeader, Button, Descriptions, Row, Col, Card, Table, Skeleton } from 'antd'
import { ReloadOutlined, DeleteOutlined, SettingOutlined } from '@ant-design/icons'
import WhoisBox from '../WhoisBox'
const { Column } = Table

const DomainDetail = (props) => {
  const { domain, whois, loading, onBack } = props
  const whoisList = [
    {
      id: 1,
      createdAt: new Date().toISOString()
    },
    {
      id: 2,
      createAt: new Date(123).toISOString()
    }
  ]
  return (
    <div>
      <PageHeader
        ghost={false}
        title={domain.name}
        subTitle={domain.id}
        onBack={onBack}
        extra={[
          <Button type="primary"><ReloadOutlined />Refresh</Button>,
          <Button><SettingOutlined />Settings</Button>,
          <Button type="danger"><DeleteOutlined />Delete</Button>
        ]}
      >
        <Descriptions size="small" column={2}>
          <Descriptions.Item label="Created At">{ new Date(domain.createdAt).toLocaleString() }</Descriptions.Item>
          <Descriptions.Item label="Updated At">{ new Date(domain.updatedAt).toLocaleString() }</Descriptions.Item>
        </Descriptions>
      </PageHeader>
      <Row style={{ marginTop: '10px' }}>
        <Col span={16}>
          <Card
            style={{ width: '100%' }}
            bodyStyle={{ minHeight: '300px' }}
            size="small"
            title="Whois Information"
          >
            <Skeleton
              loading={loading}
              paragraph={{
                rows: 10
              }}
              active>
              {whois !== null ?
                (<WhoisBox
                  {...whois}
                />)
                : null}
            </Skeleton>
          </Card>
        </Col>
        <Col offset={1} span={7}>
          <Table
            dataSource={whoisList}
            pagination={{ pageSize: 6 }}
            size="small"
          >
            <Column
              title="History WHOIS"
              dataIndex="createdAt"
              key="createdAt"
              render={createdAt => (new Date(createdAt).toLocaleString())}
            />
          </Table>
        </Col>
      </Row>
    </div>
  )
}

DomainDetail.propTypes = {
  domain: PropTypes.object,
  whois: PropTypes.object,
  loading: PropTypes.bool,
  onBack: PropTypes.func
}

export default DomainDetail
