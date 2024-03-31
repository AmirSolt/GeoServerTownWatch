package main

import (
	"fmt"
	"townwatch/base"
	"townwatch/services/areas"
	"townwatch/services/events"
	"townwatch/services/reports"

	"github.com/getsentry/sentry-go"
	"github.com/robfig/cron"
)

func main() {

	b := base.LoadBase()
	defer b.Kill()

	events.LoadInit(b)

	events.LoadRoutes(b)
	areas.LoadRoutes(b)
	reports.LoadRoutes(b)

	go func() {
		c := cron.New()

		defer func() {
			if r := recover(); r != nil {
				sentry.CaptureException(fmt.Errorf("recovered panic: %v", r))
			}
		}()

		events.LoadCronJobs(b, c)
		reports.LoadCronJobs(b, c)
		c.Start()
		// defer c.Stop()
	}()

	fmt.Println("=======")
	fmt.Println(b.DOMAIN)
	fmt.Println("=======")

	b.Run()
}
