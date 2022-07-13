package status

const (
	OK = iota
	AccessDenied
	BadParameter
	ServerError
	NotFound
	ConnectionError
	FormatError
	MessageQueueFull
	UnsupportedSuffix
	BadWhoisFormatOrNotRegistered
)
