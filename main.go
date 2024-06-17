package main

import (
	"harbor-img/cmd"
	"log"
)

func main() {
	app := cmd.NewHarborCmd()
	if err := app.Execute(); err != nil {
		log.Fatal(err)
	}
}
