import React, { useEffect, useState } from 'react'
import PropTypes from 'prop-types'
import { Button, Table, Drawer, message } from 'antd';
import { PlusOutlined, ReloadOutlined } from '@ant-design/icons'
import showErrorMessage from '../../actions/showErrorMessage';

const ItemList = (props) => {
  const [isNewDrawerVisible, setIsNewDrawerVisible] = useState(false)
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
  })
  const [showHandCursor, setShowHandCursor] = useState(false)
  const [items, setItems] = useState([])
  const { name, onLoad } = props
  const [loading, setLoading] = useState(true)

  const load = (pagination) => {
    setLoading(true)
    onLoad({ offset: (pagination.current - 1) * pagination.pageSize, limit: pagination.pageSize })
      .then(data => {
        setPagination({
          ...pagination,
          total: data.total
        })
        setItems(data.items)
      })
      .catch(err => {
        showErrorMessage(err)
      })
      .finally(() => {
        setLoading(false)
      })
  }

  useEffect(() => {
    load(pagination)
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const handleRefreshClick = () => {
    load(pagination)
  }

  const handleTableChange = (newPagination) => {
    load(newPagination)
  }

  const handleAddClick = () => {
    setIsNewDrawerVisible(true)
  }

  const handleDrawerClose = () => {
    setIsNewDrawerVisible(false)
  }

  const handleAddSubmitted = () => {
    message.success(`Successfully added ${props.name}(s)`)
    load(pagination)
    setIsNewDrawerVisible(false)
  }

  const handleAddError = () => {
    message.error(`Failed to add ${props.name}(s)`)
  }

  const { tableRowKey, tableColumns, onDelete, addForm: AddForm, onRowClick, showHandCursorInRow } = props
  let extendedTableColumns = tableColumns.map(column => {
    return {
      ...column,
      onCell: record => ({
        onClick: e => onRowClick && onRowClick(e, record),
        onMouseEnter: e => { showHandCursorInRow && setShowHandCursor(true) },
        onMouseLeave: e => { showHandCursorInRow && setShowHandCursor(false) },
      })
    }
  })
  if (typeof onDelete !== 'undefined') {
    const handleDeleteClick = (id) => {
      onDelete(id)
        .then(() => {
          message.success(`Successfully deleted ${name} ${id}`)
          load(pagination)
        })
        .catch(err => {
          showErrorMessage(err, `Error occurs when deleting ${name} ${id}.`)
        })
    }
    extendedTableColumns.push({
      title: 'Action',
      key: 'action',
      fixed: 'right',
      render: (text, record) => (
        <span>
          <Button type="danger"
                  onClick={(e) => {
                    e.preventDefault();
                    handleDeleteClick(record.id)
                  }}
                  disabled={record.isDeleting}>Delete</Button>
        </span>
      )
    })
  }
  return (
    <div>
      <div style={{ marginBottom: '10px' }}>
        <Button.Group>
          <Button onClick={handleAddClick}>
            <PlusOutlined /> Add
          </Button>
          <Button type="primary" onClick={handleRefreshClick}>
            <ReloadOutlined /> Refresh
          </Button>
        </Button.Group>
      </div>
      <Table
        columns={extendedTableColumns}
        dataSource={items}
        rowKey={tableRowKey}
        pagination={pagination}
        loading={loading}
        onChange={handleTableChange}
        style={{
          cursor: showHandCursor ? 'pointer' : 'default'
        }}
      />
      {typeof AddForm !== 'undefined' ?
        <Drawer
          title={`Add ${name}(s)`}
          onClose={handleDrawerClose}
          visible={isNewDrawerVisible}
        >
          <AddForm onSubmitted={handleAddSubmitted} onError={handleAddError}/>
        </Drawer>
        : null}
    </div>
  );
}

ItemList.propTypes = {
  name: PropTypes.string.isRequired,
  onLoad: PropTypes.func.isRequired,
  onDelete: PropTypes.func,
  onBatchDelete: PropTypes.func,
  queryable: PropTypes.bool,
  tableRowKey: PropTypes.func.isRequired,
  tableColumns: PropTypes.arrayOf(PropTypes.object).isRequired,
  addForm: PropTypes.elementType,
  onRowClick: PropTypes.func,
  showHandCursorInRow: PropTypes.bool
}

ItemList.defaultProps = {
  queryable: false
}

export default ItemList
