package main

import (
	"context"

	"github.com/friendly-fhir/fhenix/internal/cobra/cmd"
	"github.com/friendly-fhir/fhenix/internal/snek"
)

func main() {
	app := snek.NewApplication(&cmd.RootCommand{}, &snek.AppInfo{
		Website:   "https://friendly-fhir.org",
		IssueURL:  "https://github.com/friendly-fhir/fhenix/issues",
		KeyTerms:  []string{"fhenix.yaml"},
		Variables: []string{"fhenix config", "fhir registry"},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app.Execute(ctx).Exit()
}
