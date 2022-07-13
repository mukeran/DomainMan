import { message } from 'antd'
import constants from '../constants'
const { status, errorMessage } = constants

const showErrorMessage = function (err, extraMessage) {
  if (err.status === status.ok)
    return
  let _message
  if (typeof errorMessage[err.status] === 'undefined')
    _message = ((typeof extraMessage !== 'undefined') ? `${extraMessage}: ` : '') + 'Unknown error'
  else
    _message = ((typeof extraMessage !== 'undefined') ? `${extraMessage}: ` : '') + errorMessage[err.status]
  message.error(_message)
}

export default showErrorMessage