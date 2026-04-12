package processor

import "errors"

var (
	ErrMissingFields             = errors.New("missing fields")
	ErrUnsupportedMessageType    = errors.New("unsupported message type")
	ErrUnsupportedProcessingCode = errors.New("unsupported processing code")
	ErrInvalidAccountNumber      = errors.New("invalid account number")
	ErrInvalidAccountCombination = errors.New("invalid account combination")
	ErrInvalidID                 = errors.New("invalid id")
	ErrInvalidCurrencyCode       = errors.New("invalid currency code")
)
