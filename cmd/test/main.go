package main

import (
	"townwatch/base"
	"townwatch/services/events"
)

func main() {

	b := base.LoadBase()
	defer b.Kill()

	events.LoadInit(b)

	events.LoadRoutes(b)

	b.Run()
}
