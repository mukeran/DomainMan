import constants from '../constants'

export function processRequest (axiosRequest) {
  return axiosRequest.then(
    ({ data }) => {
      if (data.status !== constants.status.ok)
        return Promise.reject({ status: data.status })
      return Promise.resolve(data)
    },
    ({ response }) => {
      if (typeof response === 'undefined' || typeof response.data === 'undefined')
        return Promise.reject({ status: constants.status.connectionError, httpStatus: 0 })
      if (typeof response.data.status !== 'undefined')
        return Promise.reject({ status: response.data.status })
      return Promise.reject({ status: constants.status.connectionError, httpStatus: response.status })
    }
  )
}
