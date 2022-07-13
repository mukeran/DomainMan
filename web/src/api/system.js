import axios from 'axios'
import config from '../config'
import { processRequest } from './processRequest'

const system = {
  ping: (accessToken) =>
    processRequest(axios.get(`${config.apiRoot}/system/ping`, {
      headers: { 'Access-Token': accessToken }
    })),
}


export default system