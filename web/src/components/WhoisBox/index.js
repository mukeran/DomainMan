import React from 'react'
import PropTypes from 'prop-types'
import { Descriptions, Button } from 'antd'
import { Link } from 'react-router-dom'
import { ReloadOutlined } from '@ant-design/icons'

import constants from '../../constants'
import whoisStatus from '../../constants/whoisStatus'

class WhoisBox extends React.Component {
  static propTypes = {
    domainName: PropTypes.string.isRequired,
    createdAt: PropTypes.string.isRequired,
    raw: PropTypes.string.isRequired,
    registrant: PropTypes.string.isRequired,
    registrantEmail: PropTypes.string,
    registrar: PropTypes.string.isRequired,
    updatedDate: PropTypes.string,
    registrationDate: PropTypes.string.isRequired,
    expirationDate: PropTypes.string.isRequired,
    status: PropTypes.number.isRequired,
    nameServer: PropTypes.arrayOf(PropTypes.string).isRequired,
    dnssec: PropTypes.string,
    dsData: PropTypes.string,
    id: PropTypes.number,
    simple: PropTypes.bool,
    showRefresh: PropTypes.bool,
    onRefresh: PropTypes.func
  }

  renderStatus (status) {
    let str = []
    for (let i = 0; i < 23; i++)
      if ((status & (1 << i)) !== 0)
        str.push(whoisStatus[i])
    return str.join('\n')
  }

  render () {
    return (
      <div>
        <Descriptions
          title={
            this.props.showRefresh ?
              <div>
                <span>Queried at {new Date(this.props.createdAt).toLocaleString()}</span>
                &nbsp;<Button onClick={() => this.props.onRefresh(this.props.domainName)}>
                  <ReloadOutlined/>Refresh
                </Button>
              </div> : null
          }
          bordered>
          <Descriptions.Item label="Domain Name">{this.props.domainName}</Descriptions.Item>
          <Descriptions.Item label="Registrant">{this.props.registrant}</Descriptions.Item>
          <Descriptions.Item label="Registrant Email">{this.props.registrantEmail}</Descriptions.Item>
          <Descriptions.Item label="Registrar" span={3}>{this.props.registrar}</Descriptions.Item>
          <Descriptions.Item
            label="Registration Date">{new Date(this.props.registrationDate).toLocaleString()}</Descriptions.Item>
          <Descriptions.Item
            label="Expiration Date">{new Date(this.props.expirationDate).toLocaleString()}</Descriptions.Item>
          {this.props.updatedDate !== constants.invalidTime &&
            <Descriptions.Item label="Updated Date">
              {new Date(this.props.updatedDate).toLocaleString()}
            </Descriptions.Item>
          }
          <Descriptions.Item
            label="Status"
            span={3}
            style={{ whiteSpace: 'pre-wrap' }}
          >
            {this.renderStatus(this.props.status)}
          </Descriptions.Item>
          <Descriptions.Item
            label="Name Server"
            span={3}
            style={{ whiteSpace: 'pre-wrap' }}
          >{(this.props.nameServer || []).join('\n')}</Descriptions.Item>
        </Descriptions>
        {this.props.simple && typeof this.props.id !== 'undefined' ?
          <Link to={{ pathname: `/whois/${this.props.id}`, state: { noReload: true } }}>Show more...</Link> :
          <div>
            <span>Raw Information</span>
            <pre>{this.props.raw}</pre>
          </div>
        }
      </div>
    )
  }
}

export default WhoisBox
