package main

import (
	"fmt"
	"townwatch/base"
)

func main() {

	b := base.LoadBase()
	defer b.Kill()

	fmt.Println("=======")
	fmt.Println(b.DOMAIN)
	fmt.Println("=======")

	b.Run()
}
