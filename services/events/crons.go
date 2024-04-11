package events

import (
	"context"
	"fmt"
	"time"
	"townwatch/base"

	"github.com/getsentry/sentry-go"
	"github.com/robfig/cron"
)

func LoadCronJobs(b *base.Base, c *cron.Cron) {
	err := c.AddFunc("0 0 * * *", func() {

		sentry.CaptureMessage(fmt.Sprintf("Events cron started at: %s", time.Now().Format(time.RFC1123)))
		_, err := fetchAndStoreEvents(b, context.Background(), time.Now().Add(-time.Duration(24*4)*time.Hour).UTC(), time.Now().UTC())
		if err != nil {
			return
		}
	})

	if err != nil {
		sentry.CaptureException(err)
	}
}
