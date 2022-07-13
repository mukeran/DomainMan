import {
  UPDATE_ACCESS_TOKEN
} from '../actions/auth'

function auth (
  state = {
    accessToken: ''
  },
  action
) {
  switch (action.type) {
    case UPDATE_ACCESS_TOKEN:
      return {
        ...state,
        accessToken: action.accessToken
      }
    default:
      return state
  }
}


export default auth