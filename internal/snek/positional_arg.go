package snek

import "github.com/spf13/cobra"

// PositionalArgs is used to represent positional argument requirements for a
// command.
type PositionalArgs interface {
	positionArg() cobra.PositionalArgs
}

type positionalArgs cobra.PositionalArgs

func (pa positionalArgs) positionArg() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		err := cobra.PositionalArgs(pa)(cmd, args)
		if err == nil {
			return nil
		}
		Errorf(cmd.Context(), "%v", err)
		_ = cmd.Usage()
		return &usageError{Message: err.Error()}
	}
}

var (
	// ArbitraryArgs is used to represent a command that can take any number of
	// positional arguments.
	ArbitraryArgs = positionalArgs(cobra.ArbitraryArgs)

	// NoArgs is used to represent a command that takes no positional arguments.
	NoArgs = positionalArgs(cobra.NoArgs)
)

// MinimumNArgs returns a PositionalArgs that requires at least n arguments.
func MinimumNArgs(n int) PositionalArgs {
	return positionalArgs(cobra.MinimumNArgs(n))
}

// MaximumNArgs returns a PositionalArgs that requires at most n arguments.
func MaximumNArgs(n int) PositionalArgs {
	return positionalArgs(cobra.MaximumNArgs(n))
}

// ExactArgs returns a PositionalArgs that requires exactly n arguments.
func ExactArgs(n int) PositionalArgs {
	return positionalArgs(cobra.ExactArgs(n))
}

// RangeArgs returns a PositionalArgs that requires between min and max
func RangeArgs(min, max int) PositionalArgs {
	return positionalArgs(cobra.RangeArgs(min, max))
}

// MinimumNArgs returns a PositionalArgs that requires at least n arguments.
func MatchAll(pargs ...PositionalArgs) PositionalArgs {
	return positionalArgs(cobra.MatchAll(func(cmd *cobra.Command, args []string) error {
		for _, parg := range pargs {
			if err := parg.positionArg()(cmd, args); err != nil {
				return err
			}
		}
		return nil
	}))
}

// Condition returns a [PositionalArgs] that requires the condition to be met.
func Condition(cond func(args []string) error) PositionalArgs {
	return positionalArgs(func(cmd *cobra.Command, args []string) error {
		return cond(args)
	})
}
