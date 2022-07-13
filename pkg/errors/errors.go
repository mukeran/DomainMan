package errors

import "errors"

var (
	ErrUnsupportedSuffix             = errors.New("unsupported suffix")
	ErrUnsupportedWhoisMode          = errors.New("unsupported whois mode")
	ErrMessageQueueCancelled         = errors.New("the message queue has been cancelled")
	ErrMessageQueueFull              = errors.New("the message queue is full")
	ErrBadWhoisFormatOrNotRegistered = errors.New("bad WHOIS format or not registered")
	ErrUnsupportedDatabaseDialect    = errors.New("unsupported database dialect")
)

func Is(err error, target error) bool {
	return errors.Is(err, target)
}
