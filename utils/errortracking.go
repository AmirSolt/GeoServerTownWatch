package utils

import "github.com/getsentry/sentry-go"

type CError struct {
	*sentry.EventID `json:"event_id"`
	Message         string `json:"message"`
	Error           error  `json:"-"`
}
