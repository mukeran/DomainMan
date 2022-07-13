import axios from 'axios'
import config from '../config'
import { processRequest } from './processRequest'

const domain = {
  list: (accessToken, { offset, limit }) =>
    processRequest(axios.get(`${config.apiRoot}/domain`, {
      params: { offset, limit },
      headers: { 'Access-Token': accessToken }
    })),
  add: (accessToken, { domains }) =>
    processRequest(axios.post(`${config.apiRoot}/domain`, { domains }, {
      headers: { 'Access-Token': accessToken }
    })),
  show: (accessToken, id) =>
    processRequest(axios.get(`${config.apiRoot}/domain/${id}`, {
      headers: { 'Access-Token': accessToken }
    })),
  delete: (accessToken, id) =>
    processRequest(axios.delete(`${config.apiRoot}/domain/${id}`, {
      headers: { 'Access-Token': accessToken }
    })),
}


export default domain