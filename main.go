package main

import (
	"github.com/naftulikay/golang-snakes/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Failed to execute command: %s\n", err)
	}
}

