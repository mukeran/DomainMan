package errors

import "errors"

var (
	ErrUnsupportedSuffix     = errors.New("unsupported suffix")
	ErrMessageQueueCancelled = errors.New("the message queue has been cancelled")
	ErrMessageQueueFull      = errors.New("the message queue is full")
	ErrBadFormat             = errors.New("bad format")
)
