package reports

import (
	"townwatch/base"

	"github.com/getsentry/sentry-go"
	"github.com/robfig/cron"
)

type NotifCreateManyParams struct {
	params []NotifCreateParams
}

type NotifCreateParams struct {
	UserID   string `db:"user_id" json:"user_id"`
	Subject  string `db:"subject" json:"subject"`
	BodyHTML string `db:"body_html" json:"body_html"`
}

func LoadCronJobs(b *base.Base, c *cron.Cron) {
	err := c.AddFunc("0 30 * * * ", func() {
		CreateGlobalReports(b)
	})

	if err != nil {
		sentry.CaptureException(err)
	}
}
