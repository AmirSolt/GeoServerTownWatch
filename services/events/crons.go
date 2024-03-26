package events

import (
	"context"
	"time"
	"townwatch/base"

	"github.com/robfig/cron"
)

func LoadCronJobs(b *base.Base, c *cron.Cron) {
	c.AddFunc("0 30 * * * *", func() {
		_, err := FetchAndStoreTorontoEvents(b, context.Background(), time.Now().Add(-time.Duration(24*4)*time.Hour).UTC(), time.Now().UTC())
		if err != nil {
			return
		}
	})
}
