package snek

import (
	"context"
	"strings"
)

// Command is the primary entry-point interface for use in the snek package.
// This will be translated into a cobra.Command for use in the CLI.
type Command interface {
	// Info returns information about the command.
	Info() *CommandInfo

	// Run executes the command with the given arguments.
	Run(ctx context.Context, args []string) error

	// Complete returns a list of possible completions for the current set of args
	// and argument being completed.
	Complete(args []string, toComplete string) Completion

	// PositionalArgs returns information about the positional arguments for this
	// command.
	PositionalArgs() PositionalArgs

	// Flags returns a list of flags that can be used with this command.
	Flags() []*FlagSet

	// Commands returns a list of sub-commands that can be run under this
	// command.
	Commands() Commands

	command()
}

// BaseCommand is a simple implementation of the Command interface that must be
// embedded in all command implementations.
type BaseCommand struct{}

func (c *BaseCommand) Info() *CommandInfo {
	return nil
}

func (c *BaseCommand) Run(ctx context.Context, args []string) error {
	if len(args) > 0 {
		return UsageErrorf("invalid input %q", strings.Join(args, " "))
	}
	return errNotImplemented
}

func (c *BaseCommand) Complete(args []string, toComplete string) Completion {
	return NoCompletion
}

func (c *BaseCommand) PositionalArgs() PositionalArgs {
	return ArbitraryArgs
}

func (c *BaseCommand) Flags() []*FlagSet {
	return nil
}

func (c *BaseCommand) Commands() Commands {
	return Commands{}
}

func (c *BaseCommand) command() {}

var _ Command = (*BaseCommand)(nil)
