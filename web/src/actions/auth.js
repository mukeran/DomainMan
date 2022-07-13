export const UPDATE_ACCESS_TOKEN = 'UPDATE_ACCESS_TOKEN'

export function updateAccessToken (accessToken) {
  return {
    type: UPDATE_ACCESS_TOKEN,
    accessToken
  }
}
