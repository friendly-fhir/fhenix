package main

import (
	"os"

	"github.com/friendly-fhir/fhenix/internal/cobra/cmd"
)

func main() {
	if err := cmd.Root.Execute(); err != nil {
		// The error is already logged by Cobra
		os.Exit(1)
	}
}
