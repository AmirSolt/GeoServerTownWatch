package main

import (
	"fmt"
	"townwatch/base"
)

func main() {
	base := base.Base{
		RootDir: "./",
	}

	base.LoadBase()
	defer base.Kill()

	fmt.Println("=======")
	fmt.Println(base.DOMAIN)
	fmt.Println("=======")

	base.Run()
}
