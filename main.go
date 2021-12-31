package main

import (
	"log"

	"github.com/osamaadam/gomanz/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}