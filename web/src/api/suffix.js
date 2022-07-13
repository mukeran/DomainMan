import axios from 'axios'
import config from '../config'
import { processRequest } from './processRequest'

const suffix = {
  list: (accessToken, { query, offset, limit }) =>
    processRequest(axios.get(`${config.apiRoot}/suffix`, {
      params: { query, offset, limit },
      headers: { 'Access-Token': accessToken }
    })),
  add: (accessToken, { name, memo, description, whoisServer }) =>
    processRequest(axios.post(`${config.apiRoot}/suffix`, { name, memo, description, whoisServer }, {
      headers: { 'Access-Token': accessToken }
    })),
  show: (accessToken, id) =>
    processRequest(axios.get(`${config.apiRoot}/suffix/${id}`, {
      headers: { 'Access-Token': accessToken }
    })),
  modify: (accessToken, id, { memo, description, whoisServer }) =>
    processRequest(axios.patch(`${config.apiRoot}/suffix/${id}`, { memo, description, whoisServer }, {
      headers: { 'Access-Token': accessToken }
    })),
  delete: (accessToken, id) =>
    processRequest(axios.delete(`${config.apiRoot}/suffix/${id}`, {
      headers: { 'Access-Token': accessToken }
    }))
}


export default suffix