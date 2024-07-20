package cmd

import (
	"github.com/friendly-fhir/fhenix/internal/snek"
)

type InitCommand struct {
	snek.BaseCommand
}

func (ic *InitCommand) Info() *snek.CommandInfo {
	return &snek.CommandInfo{
		Use:     "init [name]",
		Summary: "Initializes a new fhenix project",
		Description: snek.Lines(
			"Initializes a new fhenix project",
			"",
			"Creates a new fhenix project at the location specified by [name], if provided, ",
			"or in the current directory if not.",
		),
	}
}

func (ic *InitCommand) PositionalArgs() snek.PositionalArgs {
	return snek.MaximumNArgs(1)
}
