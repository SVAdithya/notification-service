package models

import "errors"

// Application errors
var (
	ErrInvalidNotificationID = errors.New("invalid notification ID")
	ErrInvalidRecipient      = errors.New("invalid recipient email address")
	ErrMissingContent        = errors.New("missing email content")
	ErrInvalidSubject        = errors.New("invalid email subject")
	ErrSMTPError             = errors.New("SMTP server error")
	ErrInvalidEmailFormat    = errors.New("invalid email format")
	ErrMarshalFailed         = errors.New("failed to marshal data")
	ErrUnmarshalFailed       = errors.New("failed to unmarshal data")
	ErrNetworkError          = errors.New("network error")
	ErrTimeout               = errors.New("request timeout")
	ErrAuthenticationFailed  = errors.New("SMTP authentication failed")
)