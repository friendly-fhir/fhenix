package cmd

import (
	"github.com/friendly-fhir/fhenix/internal/snek"
)

type RootCommand struct {
	Verbose bool
	NoColor bool
	Output  string

	Hidden bool
	snek.BaseCommand
}

func (r *RootCommand) Info() *snek.CommandInfo {
	return &snek.CommandInfo{
		Use:     "fhenix <command>",
		Summary: "Fhenix is a lightweight tool for generating code from FHIR StructureDefinitions",
		Description: lines(
			"Fhenix is a lightweight tool for generating code from FHIR StructureDefinitions",
		),
	}
}

func (r *RootCommand) Commands() snek.Commands {
	commands := snek.Commands{}
	communication := commands.Group("Communication")
	communication.Add(&DownloadCommand{})

	generation := commands.Group("Generation")
	generation.Add(&InitCommand{})
	generation.Add(&RunCommand{})
	return commands
}

var _ snek.Command = (*RootCommand)(nil)
