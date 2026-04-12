package rpc

import "errors"

var (
	ErrAccountNotFound  = errors.New("account not found with the given info")
	ErrTokenNotFound    = errors.New("token not found")
	ErrTokenAlreadyUsed = errors.New("token already used")
)
