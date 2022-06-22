package model

import "errors"

var (
	ErrKeyNotFound = errors.New("key not found")
    ErrInvalidArgument = errors.New("invalid arguments")
    ErrTooManyRequests = errors.New("user has made too many requests")
    ErrInternalError = errors.New("internal error")
)
