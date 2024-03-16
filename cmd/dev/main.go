package main

import (
	"fmt"
	"townwatch/base"
	"townwatch/services/events"

	"github.com/robfig/cron"
)

func main() {

	b := base.LoadBase()
	defer b.Kill()

	events.LoadRoutes(b)

	go func() {
		c := cron.New()
		events.LoadCronJobs(b, c)
		c.Start()
		defer c.Stop()
	}()

	fmt.Println("=======")
	fmt.Println(b.DOMAIN)
	fmt.Println("=======")

	b.Run()
}
