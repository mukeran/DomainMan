import React, { useEffect, useState } from 'react'
import { connect } from 'react-redux'

import api from '../../api'
import { Skeleton, Empty } from 'antd'
import DomainDetail from '../../components/DomainDetail'
import { useParams } from 'react-router-dom'
import { useNavigate } from 'react-router-dom'
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

const DomainShow = (props) => {
  const [domain, setDomain] = useState({})
  const [whois, setWhois] = useState({})
  const [domainLoading, setDomainLoading] = useState(true)
  const [whoisLoading, setWhoisLoading] = useState(true)

  const params = useParams()

  const { accessToken } = props

  const loadDomain = () => {
    setDomainLoading(true)
    api.domain.show(accessToken, params.domainID)
      .then(data => {
        setDomain(data.domain)
        setDomainLoading(false)
        setWhoisLoading(true)
        api.whois.query(accessToken, { domain: data.domain.name, forceUpdate: false })
          .then(data => {
            setWhois(data.whois)
          })
          .catch(err => {
            showErrorMessage(err, 'Failed to load domain whois')
          })
          .finally(() => {
            setWhoisLoading(false)
          })
      })
      .catch(err => {
        showErrorMessage(err, 'Failed to load domain')
        setDomainLoading(false)
      })
  }

  useEffect(() => {
    loadDomain()
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [params.domainID])

  const navigate = useNavigate()

  return (
    <div>
      <Skeleton
        loading={domainLoading}
        paragraph={{
          rows: 10
        }}
        active>
        { domain !== null ?
          <DomainDetail
            domain={domain}
            whois={whois}
            loading={whoisLoading}
            onBack={() => navigate('/domain')}
          />
          : <Empty/>
        }
      </Skeleton>
    </div>
  )
}

export default connect(mapStateToProps)(DomainShow)
