package main

import (
	"log"

	"harbor-img-clear/cmd"
)

func main() {
	newClean := cmd.NewHarborCleanCommand()
	if err := newClean.Execute(); err != nil {
		log.Fatal(err)
	}

}
