package events

import (
	"context"
	"net/http"
	"townwatch/base"

	"github.com/carlmjohnson/requests"
	"github.com/getsentry/sentry-go"
	"github.com/robfig/cron"
)

func LoadCronJobs(b *base.Base, c *cron.Cron) {
	c.AddFunc("0 30 * * * *", func() {
		err := requests.
			URL("/api/events/fetch").
			Method(http.MethodGet).
			Header(base.HeaderSecretKeyName, b.SECRET_API_KEY).
			Fetch(context.Background())
		if err != nil {
			sentry.CaptureException(err)
		}
	})
}
