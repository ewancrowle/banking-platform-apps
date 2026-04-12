package rpc

import "errors"

var (
	ErrEmailExists      = errors.New("account with email already exists")
	ErrIDNotFound       = errors.New("no account found with the given id")
	ErrEmailNotFound    = errors.New("no account found with the given email")
	ErrPasswordMismatch = errors.New("password mismatch")
)
