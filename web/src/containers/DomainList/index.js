import React from 'react'

import ItemList from '../../components/ItemList'
import DomainAddForm from '../../components/DomainAddForm'
import { useNavigate } from 'react-router-dom'
import api from '../../api'
import { connect } from 'react-redux'

const mapStateToProps = state => {
  const { auth } = state
  const { accessToken } = auth || {
    accessToken: ''
  }
  return {
    accessToken
  }
}

const DomainList = (props) => {
  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id'
    },
    {
      title: 'Domain Name',
      dataIndex: 'name',
      key: 'name'
    },
    {
      title: 'Created At',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: createdAt => (new Date(createdAt)).toLocaleString()
    },
    {
      title: 'Updated At',
      dataIndex: 'updatedAt',
      key: 'updatedAt',
      render: updatedAt => (new Date(updatedAt)).toLocaleString()
    }
  ]
  const navigate = useNavigate()
  const { accessToken } = props

  const handleLoad = ({ offset, limit }) => {
    return api.domain.list(accessToken, { offset, limit })
      .then(data => {
        return {
          items: data.domains,
          total: data.total
        }
      })
  }

  const handleDelete = (id) => {
    return api.domain.delete(accessToken, id)
  }

  return <ItemList
    name="domain"
    onLoad={handleLoad}
    onDelete={handleDelete}
    queryable
    tableRowKey={domain => domain.id}
    tableColumns={columns}
    addForm={DomainAddForm}
    onRowClick={(e, record) => {
      navigate(`/domain/${record.id}`)
    }}
    showHandCursorInRow
  />
}

export default connect(mapStateToProps)(DomainList)
