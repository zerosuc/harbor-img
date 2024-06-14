package main

import (
	"harbor-img-clear/cmd"
	"log"
)

func main() {
	newClean := cmd.NewHarborCleanCommand()
	if err := newClean.Execute(); err != nil {
		log.Fatal(err)
	}
}
