package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/friendly-fhir/fhenix/internal/cobra/cmd"
	"github.com/friendly-fhir/fhenix/internal/snek"
)

var (
	version = "snapshot"
	date    string
)

const (
	defaultTimeout = time.Minute * 2
)

func main() {
	date, err := time.Parse(time.RFC3339, date)
	if err != nil {
		date = time.Now()
	}
	root := cmd.RootCommand{}
	app := snek.NewApplication(&root, &snek.AppInfo{
		Website:        "https://friendly-fhir.org",
		DocsURL:        "https://friendly-fhir.org/fhenix",
		ReportIssueURL: "https://github.com/friendly-fhir/fhenix/issues",
		KeyTerms:       []string{"fhenix.yaml"},
		Variables:      []string{"fhenix config", "fhir registry"},
		Version:        version,
		Date:           date,
	})

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	app.Execute(ctx).Exit()
}
