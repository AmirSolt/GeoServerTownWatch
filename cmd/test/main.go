package main

import (
	"bytes"
	"fmt"
	"text/template"
)

type Pet struct {
	Name   string
	Sex    string
	Intact bool
	Age    string
	Breed  string
}

func main() {
	dogs := []Pet{
		{
			Name:   "Jujube",
			Sex:    "Female",
			Intact: false,
			Age:    "10 months",
			Breed:  "German Shepherd/Pitbull",
		},
		{
			Name:   "Zephyr",
			Sex:    "Male",
			Intact: true,
			Age:    "13 years, 3 months",
			Breed:  "German Shepherd/Border Collie",
		},
	}

	fileName := "pets.tmpl"
	filePath := fmt.Sprintf("./cmd/test/%s", fileName)
	tmpl, err := template.New(fileName).ParseFiles(filePath)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, dogs)
	if err != nil {
		panic(err)
	}

	fmt.Println(buf.String())
}
