package main

import (
	"context"
	"strings"

	"github.com/friendly-fhir/fhenix/internal/snek"
)

type ExampleSubCommand struct {
	name        string
	description string
	Value       string
	snek.BaseCommand
}

func (esc *ExampleSubCommand) Info() *snek.CommandInfo {
	return &snek.CommandInfo{
		Use:     esc.name,
		Summary: esc.description,
	}
}

func (esc *ExampleSubCommand) Run(ctx context.Context, args []string) error {
	snek.Noticef(ctx, "example notice")
	snek.Errorf(ctx, "example error")
	snek.Warningf(ctx, "example warning")
	return nil
}

func lines(lines ...string) string {
	return strings.Join(lines, "\n")
}

type RootCommand struct {
	Verbose bool
	NoColor bool
	Output  string

	Hidden bool
	snek.BaseCommand
}

func (rc *RootCommand) Info() *snek.CommandInfo {
	return &snek.CommandInfo{
		Use:     "fhenix <command>",
		Summary: "Fhenix is a lightweight tool for generating code from FHIR StructureDefinitions",
		Description: lines(
			"Fhenix is a lightweight tool for generating code from FHIR StructureDefinitions",
			"",
			"This tool manages downloading FHIR definitions, building a semantic model",
			"of the content, and generating content from those definitions",
		),
	}
}

func (rc *RootCommand) Flags() []*snek.FlagSet {
	communication := snek.NewFlagSet("Communication")
	communication.BoolP(&rc.Verbose, "verbose", "v", false, "Enable verbose output")

	output := snek.NewFlagSet("Output")
	output.Bool(&rc.NoColor, "no-color", false, "Disable color output")
	output.StringP(&rc.Output, "output", "o", "", "The specified output path (default: stdout)")
	output.Bool(&rc.Hidden, "hidden", false, "A hidden flag").MarkHidden()
	return snek.FlagSets(communication, output)
}

func (rc *RootCommand) Run(ctx context.Context, args []string) error {
	snek.Noticef(ctx, "example notice")
	snek.Errorf(ctx, "example error")
	snek.Warningf(ctx, "example warning")
	return nil
}

func (rc *RootCommand) Commands() snek.Commands {
	commands := snek.Commands{}

	g1 := commands.Group("Project")
	g1.Add(&ExampleSubCommand{
		name:        "init",
		description: "Initialize a new project",
	})
	g1.Add(&ExampleSubCommand{
		name:        "run",
		description: "Generate code from the specified configurations",
	})

	g2 := commands.Group("Registry")
	g2.Add(&ExampleSubCommand{
		name:        "download",
		description: "Download packages from the specified registry",
	})
	g2.Add(&ExampleSubCommand{
		name:        "registry",
		description: "Manage the known registries",
	})
	return commands
}

var _ snek.Command = (*RootCommand)(nil)

func main() {
	app := snek.NewApplication("snek-cli", &RootCommand{})

	app.Execute(context.Background())
}
