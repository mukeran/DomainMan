import React, { useState } from 'react'
import { Input, Skeleton } from 'antd'
import { connect } from 'react-redux'

import WhoisBox from '../../components/WhoisBox'
import './style.less'
import api from '../../api'
import showErrorMessage from '../../actions/showErrorMessage'

const { Search } = Input

const mapStateToProps = state => {
  const { auth } = state
  const { accessToken } = auth || {
    accessToken: ''
  }
  return {
    accessToken
  }
}

const WhoisQuery = (props) => {
  const { accessToken } = props
  const [loading, setLoading] = useState(false)
  const [whois, setWhois] = useState(null)

  const handleQuery = (domain) => {
    setLoading(true)
    api.whois.query(accessToken, { domain, forceUpdate: false })
      .then(data => {
        setWhois(data.whois)
      })
      .catch(err => {
        showErrorMessage(err, 'Failed to load domain whois')
      })
      .finally(() => {
        setLoading(false)
      })
  }
  
  const handleRefresh = (domain) => {
    setLoading(true)
    api.whois.query(accessToken, { domain, forceUpdate: true })
      .then(data => {
        setWhois(data.whois)
      })
      .catch(err => {
        showErrorMessage(err, 'Failed to load domain whois')
      })
      .finally(() => {
        setLoading(false)
      })
  }

  return (
    <div>
      <div className="whois-query-row">
        <Search
          className="whois-query-input"
          placeholder="Domain Name"
          enterButton="Query"
          onSearch={handleQuery}
          size="large"
        />
      </div>
      <div className="whois-result-row">
        <Skeleton
          loading={loading}
          paragraph={{
            rows: 10
          }}
          active>
          {whois !== null ?
            (<WhoisBox
              {...whois}
              onRefresh={handleRefresh}
              simple
            />)
            : null}
        </Skeleton>
      </div>
    </div>
  )
}

export default connect(mapStateToProps)(WhoisQuery)
