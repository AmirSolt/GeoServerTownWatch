package base

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

type CError struct {
	*sentry.EventID
	UserMsg string
	DevMsg  string
	Error   error
}

func (base *Base) loadLogging() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              base.GLITCHTIP_DSN,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		Debug:            !base.IS_PROD,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	// base.Engine.Use(SentryGinNew(SentryGinOptions{}))

}

func (base *Base) killLogging() {
	sentry.Flush(time.Second * 5)
}
