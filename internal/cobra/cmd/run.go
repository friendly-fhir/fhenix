package cmd

import (
	"context"
	"os"
	"time"

	"github.com/friendly-fhir/fhenix/config"
	"github.com/friendly-fhir/fhenix/driver"
	"github.com/spf13/cobra"
)

var RunFlags struct {
	Output string

	Force       bool
	RM          bool
	FHIRIGCache string
	Timeout     time.Duration
}

var Run = &cobra.Command{
	Use:   "run <config file> [--rm] [--output <output directory>] [--fhirig-cache <cache path>] [--timeout <timeout>]",
	Short: "Run generation",
	Long:  "Run the generation process",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		cfg, err := config.FromFile(args[0])
		if err != nil {
			return err
		}

		if RunFlags.RM {
			if err := os.RemoveAll(cfg.OutputDir); err != nil {
				return err
			}
		}

		driver, err := driver.New(cfg, nil)
		if err != nil {
			return err
		}
		ctx := cmd.Context()
		if timeout := RunFlags.Timeout; timeout != 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		return driver.Run(ctx)
	},
}

func init() {
	Root.AddCommand(Run)
	flags := Run.Flags()
	flags.StringVarP(&RunFlags.Output, "output", "o", "", "The output directory to write the generated code to")
	flags.BoolVar(&RunFlags.RM, "rm", false, "Remove all contents from the output directory prior to writing")
	flags.StringVar(&RunFlags.FHIRIGCache, "fhirig-cache", "", "The configuration path to download the FHIR IGs to")
	flags.BoolVar(&RunFlags.Force, "force", false, "Force download of FHIR IGs")
	flags.DurationVar(&RunFlags.Timeout, "timeout", 0, "Timeout for the download")
}
