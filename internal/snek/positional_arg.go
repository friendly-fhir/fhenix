package snek

import "github.com/spf13/cobra"

type PositionalArgs interface {
	positionArg() cobra.PositionalArgs
}

type positionalArgs cobra.PositionalArgs

func (pa positionalArgs) positionArg() cobra.PositionalArgs {
	return cobra.PositionalArgs(pa)
}

var (
	ArbitraryArgs = positionalArgs(cobra.ArbitraryArgs)

	NoArgs = positionalArgs(cobra.NoArgs)
)

func MinimumNArgs(n int) PositionalArgs {
	return positionalArgs(cobra.MinimumNArgs(n))
}

func MaximumNArgs(n int) PositionalArgs {
	return positionalArgs(cobra.MaximumNArgs(n))
}

func ExactArgs(n int) PositionalArgs {
	return positionalArgs(cobra.ExactArgs(n))
}

func RangeArgs(min, max int) PositionalArgs {
	return positionalArgs(cobra.RangeArgs(min, max))
}

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
