package main

import (
	"fmt"
	"townwatch/base"
	"townwatch/services/events"

	"github.com/getsentry/sentry-go"
	"github.com/robfig/cron"
)

func main() {

	b := base.LoadBase()
	defer b.Kill()

	events.LoadInit(b)

	events.LoadRoutes(b)

	go func() {
		c := cron.New()

		defer func() {
			if r := recover(); r != nil {
				sentry.CaptureException(fmt.Errorf("recovered panic: %v", r))
			}
		}()

		events.LoadCronJobs(b, c)
		c.Start()
		// defer c.Stop()
	}()

	b.Run()
}
