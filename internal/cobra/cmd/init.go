package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

func lines(lines ...string) string {
	return strings.Join(lines, "\n")
}

var InitCommand = &cobra.Command{
	Use:   "init [name]",
	Short: "Initializes a new fhenix project",
	Long: lines(
		"Initializes a new fhenix project",
		"",
		"Creates a new fhenix project at the location specified by [name], if provided, ",
		"or in the current directory if not.",
	),
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
}
