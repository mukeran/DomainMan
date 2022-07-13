import React from 'react'
import { connect } from 'react-redux'

import ItemList from '../../components/ItemList'
import SuffixAddForm from '../../components/SuffixAddForm'
import api from '../../api'

const mapStateToProps = state => {
  const { auth } = state
  const { accessToken } = auth || {
    accessToken: ''
  }
  return {
    accessToken
  }
}

const SuffixList = (props) => {
  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id'
    },
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name'
    },
    {
      title: 'Memo',
      dataIndex: 'memo',
      key: 'memo'
    },
    {
      title: 'Whois Server',
      dataIndex: 'whoisServer',
      key: 'whoisServer'
    },
    {
      title: 'Created At',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: createdAt => (new Date(createdAt)).toLocaleString()
    },
    {
      title: 'Updated At',
      dataIndex: 'createdAt',
      key: 'updatedAt',
      render: updatedAt => (new Date(updatedAt)).toLocaleString()
    }
  ]
  const { accessToken } = props

  const handleLoad = ({ offset, limit }) => {
    return api.suffix.list(accessToken, { offset, limit })
      .then(data => {
        return {
          items: data.suffixes,
          total: data.total
        }
      })
  }

  const handleDelete = (id) => {
    return api.suffix.delete(accessToken, id)
  }

  return <ItemList
    {...props}
    name="suffix"
    onLoad={handleLoad}
    onDelete={handleDelete}
    queryable
    tableRowKey={suffix => suffix.id}
    tableColumns={columns}
    addForm={SuffixAddForm}
  />
}

export default connect(mapStateToProps)(SuffixList)
