import status from './status'

const errorMessage = {
  [status.accessDenied]: 'Access denied',
  [status.badParameter]: 'Bad parameter',
  [status.connectionError]: 'Failed to contact server',
  [status.serverError]: 'Server error',
  [status.notFound]: 'Not found',
  [status.serverConnectionError]: 'Server side connection error',
  [status.formatError]: 'Format error',
  [status.messageQueueFull]: 'Message queue full',
  [status.unsupportedSuffix]: 'Unsupported suffix',
  [status.badWhoisFormatOrNotRegistered]: 'Bad whois format or not registered',
}

export default errorMessage