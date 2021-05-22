package main

import (
	"fmt"
	"github.com/naftulikay/golang-snakes/cmd"
	"os"
)

func main() {
	// this is obviously the main entrypoint to the go application
	// we call the public Execute() method from the cmd module, and if it returns an error, we report that error to
	// standard error and exit 1
	if err := cmd.Execute(); err != nil {
		fmt.Errorf("ERROR: %s\n", err)
		os.Exit(1)
	}
}
