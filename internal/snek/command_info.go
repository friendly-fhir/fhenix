package snek

type CommandInfo struct {
	// Use is an example usage of the command.
	// This must be prefixed with the command name.
	Use string

	Aliases     []string
	Summary     string
	Description string
	Examples    []string

	// KeyTerms are words that are important in the command, and will be
	// highlighted in the help output. (Optional)
	KeyTerms []string

	// Variables are words that are variable inputs in the command, and will be
	// highlighted in the help output. (Optional)
	Variables []string

	Annotations Annotations

	Version string

	Hidden     bool
	ShowCursor bool
}
