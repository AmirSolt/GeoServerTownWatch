package main

import (
	"fmt"
	"townwatch/base"
	"townwatch/services/events"
)

func main() {

	b := base.LoadBase()
	defer b.Kill()

	events.LoadRoutes(b)

	fmt.Println("=======")
	fmt.Println(b.DOMAIN)
	fmt.Println("=======")

	b.Run()
}
