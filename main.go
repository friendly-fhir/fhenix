package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/friendly-fhir/fhenix/internal/cobra/cmd"
	"github.com/friendly-fhir/fhenix/internal/snek"
)

const (
	defaultTimeout = time.Minute * 2
)

func main() {
	root := cmd.RootCommand{}
	app := snek.NewApplication(&root, &snek.AppInfo{
		Website:        "https://friendly-fhir.org",
		DocsURL:        "https://friendly-fhir.org/fhenix",
		ReportIssueURL: "https://github.com/friendly-fhir/fhenix/issues",
		KeyTerms:       []string{"fhenix.yaml"},
		Variables:      []string{"fhenix config", "fhir registry"},
	})

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	app.Execute(ctx).Exit()
}
