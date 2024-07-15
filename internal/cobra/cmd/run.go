package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/friendly-fhir/fhenix/config"
	"github.com/friendly-fhir/fhenix/driver"
	"github.com/friendly-fhir/fhenix/registry"
	"github.com/spf13/cobra"
)

type Listener struct {
	out io.Writer
	driver.BaseListener
}

func (l *Listener) BeforeDownload() {
	fmt.Fprintln(l.out, "[1] Downloading FHIR Packages...")
}

func (l *Listener) AfterDownload(err error) {
	if err != nil {
		fmt.Fprintf(l.out, "[1] Downloading FHIR Packages... Failed: %v\n", err)
	} else {
		fmt.Fprintln(l.out, "[1] Downloading FHIR Packages... Succeeded")
	}
}

func (l *Listener) BeforeLoadTransform() {
	fmt.Fprintln(l.out, "[4] Loading Transformations...")
}

func (l *Listener) AfterLoadTransform(err error) {
	if err != nil {
		fmt.Fprintf(l.out, "[4] Loading Transformations... Failed: %v\n", err)
	} else {
		fmt.Fprintln(l.out, "[4] Loading Transformations... Succeeded")
	}
}

func (l *Listener) BeforeLoadConformance() {
	fmt.Fprintln(l.out, "[2] Loading Conformance...")
}

func (l *Listener) AfterLoadConformance(err error) {
	if err != nil {
		fmt.Fprintf(l.out, "[2] Loading Conformance... Failed: %v\n", err)
	} else {
		fmt.Fprintln(l.out, "[2] Loading Conformance... Succeeded")
	}
}

func (l *Listener) BeforeLoadModel() {
	fmt.Fprintln(l.out, "[3] Loading Model...")
}

func (l *Listener) AfterLoadModel(err error) {
	if err != nil {
		fmt.Fprintf(l.out, "[3] Loading Model... Failed: %v\n", err)
	} else {
		fmt.Fprintln(l.out, "[3] Loading Model... Succeeded")
	}
}

func (l *Listener) BeforeTransform() {
	fmt.Fprintln(l.out, "[5] Transforming Model...")
}

func (l *Listener) OnTransform(output string) {
	fmt.Fprintf(l.out, "[5] Transforming Model... %s\n", output)
}

func (l *Listener) AfterTransform(jobs int, err error) {
	if err != nil {
		fmt.Fprintf(l.out, "[5] Transforming Model... Failed: %v\n", err)
	} else {
		fmt.Fprintln(l.out, "[5] Transforming Model... Succeeded")
	}
}

func (l *Listener) OnFetch(registry, pkg, value string, bytes int64) {
	fmt.Fprintf(l.out, "[1] ... Downloading %s::%s/%s (%d bytes)\n", registry, pkg, value, bytes)
}

var _ driver.Listener = (*Listener)(nil)

var RunFlags struct {
	Output string

	Parallel  int
	Force     bool
	RM        bool
	FHIRCache string
	Timeout   time.Duration
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
		cache := registry.DefaultCache()
		if RunFlags.FHIRCache != "" {
			cache = registry.NewCache(RunFlags.FHIRCache)
		}

		driver, err := driver.New(cfg,
			driver.ForceDownload(RunFlags.Force),
			driver.Parallel(RunFlags.Parallel),
			driver.Cache(cache),
			driver.Listeners(&Listener{out: cmd.OutOrStdout()}),
		)
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
	flags.StringVar(&RunFlags.FHIRCache, "fhir-cache", "", "The configuration path to download the FHIR IGs to")
	flags.BoolVar(&RunFlags.Force, "force", false, "Force download of FHIR IGs")
	flags.DurationVar(&RunFlags.Timeout, "timeout", 0, "Timeout for the download")
	flags.IntVar(&RunFlags.Parallel, "parallel", runtime.NumCPU(), "The number of parallel workers to use")
}
