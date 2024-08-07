package cmd

import (
	"github.com/friendly-fhir/fhenix/internal/snek"
)

type RootCommand struct {
	Verbose bool
	NoColor bool
	Output  string

	snek.BaseCommand
}

func (r *RootCommand) Info() *snek.CommandInfo {
	return &snek.CommandInfo{
		Use:     "fhenix <command>",
		Summary: "Fhenix is a lightweight tool for generating code from FHIR StructureDefinitions",
		Description: snek.Lines(
			"Fhenix is a lightweight tool for generating code from FHIR StructureDefinitions",
		),
		Examples: snek.Examples(
			"fhenix init",
			"fhenix download hl7.fhir.r4.core 4.0.1 --registry https://packages.simplifier.net",
			"fhenix run fhenix.yaml --parallel 4",
		),
	}
}

func (r *RootCommand) Commands() snek.Commands {
	commands := snek.Commands{}
	registry := commands.Group("Registry")
	registry.Add(&DownloadCommand{})

	generation := commands.Group("Generation")
	generation.Add(&InitCommand{})
	generation.Add(&RunCommand{})
	return commands
}

var _ snek.Command = (*RootCommand)(nil)
