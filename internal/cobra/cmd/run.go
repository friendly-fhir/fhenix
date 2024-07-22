package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/friendly-fhir/fhenix/config"
	"github.com/friendly-fhir/fhenix/driver"
	"github.com/friendly-fhir/fhenix/internal/snek"
	"github.com/friendly-fhir/fhenix/internal/snek/terminal"
	"github.com/friendly-fhir/fhenix/registry"
)

type RunCommand struct {
	Output string

	Parallel  int
	Force     bool
	RM        bool
	FHIRCache string
	Verbose   bool
	Timeout   time.Duration

	NoProgress bool
	Log        string
	snek.BaseCommand
}

func (rc *RunCommand) Info() *snek.CommandInfo {
	return &snek.CommandInfo{
		Use:     "run <fhenix config>",
		Summary: "Run generation",
		Description: snek.Lines(
			fmt.Sprintf("Run the generation process against the specified %v file", snek.FormatKeyword.Format("fhenix config")),
			"",
			"This command will download the relevant FHIR definitions if it is not already cached",
			"and generate the code based on the configuration provided.",
		),
		Examples: snek.Examples(
			"fhenix run fhenix.yaml --rm --output ./destination",
			"fhenix run fhenix.yaml --fhir-cache ~/.fhir --timeout 5m",
			"fhenix run fhenix.yaml --parallel 4",
		),
	}
}

func (rc *RunCommand) PositionalArgs() snek.PositionalArgs {
	return snek.ExactArgs(1)
}

func (rc *RunCommand) Flags() []*snek.FlagSet {
	communication := snek.NewFlagSet("Communication")
	communication.DurationP(&rc.Timeout, "timeout", "t", 0, "Timeout for the download")
	communication.BoolP(&rc.Force, "force", "f", false, "Force download of FHIR IGs")
	communication.Int(&rc.Parallel, "parallel", runtime.NumCPU(), "The number of parallel workers to use")

	output := snek.NewFlagSet("Output")
	output.Bool(&rc.RM, "rm", false, "Remove all contents from the output directory prior to writing")
	output.StringP(&rc.Output, "output", "o", "", "The output directory to write the generated code to")
	output.String(&rc.FHIRCache, "fhir-cache", "", "The configuration path to download the FHIR IGs to")
	output.BoolP(&rc.Verbose, "verbose", "v", false, "Enable verbose output")
	output.Bool(&rc.NoProgress, "no-progress", false, "Disable progress output")
	output.String(&rc.Log, "log", "", "The log file to write the output to")

	return []*snek.FlagSet{
		output,
		communication,
	}
}

func (rc *RunCommand) Run(ctx context.Context, args []string) error {
	if len(args) != 1 {
		return snek.UsageError("expected exactly one argument")
	}

	cfg, err := config.FromFile(args[0])
	if err != nil {
		return err
	}

	if rc.RM {
		if err := os.RemoveAll(cfg.OutputDir); err != nil {
			return err
		}
	}

	var listeners []driver.Listener
	if term, ok := terminal.New(snek.CommandOut(ctx)); !rc.NoProgress && ok {
		term.HideCursor()

		listeners = append(listeners, NewProgressListener(term, rc.Verbose))

		defer term.Close()
		defer term.Bottom()
	} else {
		listeners = append(listeners, NewLogListener(snek.CommandOut(ctx), rc.Verbose))
	}

	if rc.Log != "" {
		log, err := os.Create(rc.Log)
		if err != nil {
			return err
		}
		defer log.Close()
		listeners = append(listeners, NewLogListener(log, true))
	}

	cache := registry.DefaultCache()
	if rc.FHIRCache != "" {
		cache = registry.NewCache(rc.FHIRCache)
	}

	if timeout := rc.Timeout; timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	reporter := func(cause error) {
		snek.Warningf(ctx, "template error: %v", cause)
	}

	opts := []driver.Option{
		driver.ForceDownload(rc.Force),
		driver.Parallel(rc.Parallel),
		driver.Cache(cache),
		driver.Listeners(listeners...),
		driver.TemplateReportFunc(reporter),
	}
	driver, err := driver.New(cfg, opts...)
	if err != nil {
		return err
	}

	return driver.Run(ctx)
}

var _ snek.Command = (*RunCommand)(nil)
