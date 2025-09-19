package models

import "errors"

// Application errors
var (
	ErrInvalidNotificationID = errors.New("invalid notification ID")
	ErrInvalidRecipient      = errors.New("invalid recipient")
	ErrMissingContent        = errors.New("missing message content")
	ErrInvalidMediaType      = errors.New("invalid media type")
	ErrAPICallFailed         = errors.New("WhatsApp API call failed")
	ErrInvalidTemplate       = errors.New("invalid template")
	ErrInvalidPhoneNumber    = errors.New("invalid phone number")
	ErrMarshalFailed         = errors.New("failed to marshal data")
	ErrUnmarshalFailed       = errors.New("failed to unmarshal data")
	ErrNetworkError          = errors.New("network error")
	ErrTimeout               = errors.New("request timeout")
)