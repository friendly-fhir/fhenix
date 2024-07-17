package main

import (
	"context"

	"github.com/friendly-fhir/fhenix/internal/cobra/cmd"
	"github.com/friendly-fhir/fhenix/internal/snek"
)

func main() {
	app := snek.NewApplication(&cmd.RootCommand{}, &snek.AppInfo{
		Website:  "https://friendly-fhir.org",
		IssueURL: "https://github.com/friendly-fhir/fhenix/issues",
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = app.Execute(ctx)
}
