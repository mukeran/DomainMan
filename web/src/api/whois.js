import axios from 'axios'
import config from '../config'
import { processRequest } from './processRequest'

const whois = {
  list: (accessToken, { query, offset, limit }) =>
    processRequest(axios.get(`${config.apiRoot}/whois`, {
      params: { query, offset, limit },
      headers: { 'Access-Token': accessToken }
    })),
  query: (accessToken, { domain, forceUpdate }) =>
    processRequest(axios.post(`${config.apiRoot}/whois`, { domain, forceUpdate }, {
      headers: { 'Access-Token': accessToken }
    })),
  show: (accessToken, id) =>
    processRequest(axios.get(`${config.apiRoot}/whois/${id}`, {
      headers: { 'Access-Token': accessToken }
    })),
  delete: (accessToken, id) =>
    processRequest(axios.delete(`${config.apiRoot}/whois/${id}`, {
      headers: { 'Access-Token': accessToken }
    }))
}


export default whois