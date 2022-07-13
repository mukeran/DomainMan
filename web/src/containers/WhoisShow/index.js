import React, { useEffect, useState } from 'react'
import { connect } from 'react-redux'

import WhoisBox from '../../components/WhoisBox'
import { Skeleton } from 'antd'
import { useParams } from 'react-router-dom'
import { useNavigate } from 'react-router-dom'
import api from '../../api'
import showErrorMessage from '../../actions/showErrorMessage'

const mapStateToProps = state => {
  const { auth } = state
  const { accessToken } = auth || {
    accessToken: ''
  }
  return {
    accessToken
  }
}

const WhoisShow = (props) => {
  const params = useParams()

  const [loading, setLoading] = useState(true)
  const [whois, setWhois] = useState(null)
  const { accessToken } = props


  useEffect(() => {
    setLoading(true)
    api.whois.show(accessToken, params.whoisID)
      .then(data => {
        setWhois(data.whois)
      })
      .catch(err => {
        showErrorMessage(err, 'Failed to load domain whois')
      })
      .finally(() => {
        setLoading(false)
      })
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [params.whoisID])

  const navigate = useNavigate()

  const handleRefresh = (domain) => {
    setLoading(true)
    api.whois.query(accessToken, { domain, forceUpdate: true })
      .then(data => {
        navigate(`/whois/${data.whois.id}`)
      })
      .catch(err => {
        showErrorMessage(err, 'Failed to refresh domain whois')
      })
      .finally(() => {
        setLoading(false)
      })
  }

  return (
    <div>
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
          />)
          : null}
      </Skeleton>
    </div>
  )
}

export default connect(mapStateToProps)(WhoisShow)
