package snek

type CommandInfo struct {
	// Use is an example usage of the command.
	// This must be prefixed with the command name.
	Use string

	Aliases     []string
	Summary     string
	Description string
	Examples    []string

	Annotations Annotations

	Version string

	Hidden bool
}
