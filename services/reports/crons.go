package reports

import (
	"context"
	"townwatch/base"

	"github.com/getsentry/sentry-go"
	"github.com/robfig/cron"
)

func LoadCronJobs(b *base.Base, c *cron.Cron) {
	err := c.AddFunc("0 30 * * * ", func() {
		CreateGlobalReports(b, context.Background())
	})

	if err != nil {
		sentry.CaptureException(err)
	}
}
