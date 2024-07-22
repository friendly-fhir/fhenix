package cmd

import (
	"context"
	"os"
	"runtime"
	"time"

	"github.com/friendly-fhir/fhenix/driver"
	"github.com/friendly-fhir/fhenix/internal/snek"
	"github.com/friendly-fhir/fhenix/internal/snek/terminal"
	"github.com/friendly-fhir/fhenix/registry"
)

type DownloadCommand struct {
	Timeout time.Duration
	Force   bool
	Verbose bool

	ExcludeDependencies bool
	Parallel            int

	CacheDir   string
	Registry   string
	AuthToken  string
	NoProgress bool

	Log string

	snek.BaseCommand
}

func (dc *DownloadCommand) Info() *snek.CommandInfo {
	return &snek.CommandInfo{
		Use:     "download <package> <version>",
		Summary: "Download FHIR IGs",
		Description: snek.Lines(
			"Download FHIR Implementation Guides (IGs) from the web",
		),
		Examples: snek.Examples(
			"fhenix download hl7.fhir.r4.core 4.0.1",
			"fhenix download hl7.fhir.us.core 6.1.0 --force --timeout 1m",
		),
	}
}

func (dc *DownloadCommand) Run(ctx context.Context, args []string) error {
	if len(args) != 2 {
		return snek.UsageError("expected exactly two arguments")
	}

	pkg, version := args[0], args[1]

	if dc.Timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, dc.Timeout)
		defer cancel()
	}
	var cache *registry.Cache
	if dir := dc.CacheDir; dir != "" {
		cache = registry.NewCache(dir)
	} else {
		cache = registry.DefaultCache()
	}
	var listener driver.Listener
	listener = NewLogListener(snek.CommandOut(ctx), dc.Verbose)
	if term, ok := terminal.New(snek.CommandOut(ctx)); !dc.NoProgress && ok {
		term.HideCursor()

		listener = NewProgressListener(term, dc.Verbose)

		defer term.Close()
		defer term.Bottom()
	}

	if dc.Log != "" {
		log, err := os.Create(dc.Log)
		if err != nil {
			return err
		}
		defer log.Close()
		cache.AddListener(NewLogListener(log, true))
	}

	cache.AddListener(listener)

	var opts []registry.Option
	opts = append(opts, registry.URL(dc.Registry))
	if token := dc.AuthToken; token != "" {
		opts = append(opts, registry.Auth(registry.StaticTokenSource(token)))
	}

	client, err := registry.NewClient(ctx, opts...)
	if err != nil {
		return err
	}
	cache.AddClient("default", client)

	downloader := registry.NewDownloader(cache).Force(dc.Force).Workers(dc.Parallel)

	downloader.Add("default", pkg, version, !dc.ExcludeDependencies)
	if err := downloader.Start(ctx); err != nil {
		return err
	}
	return nil
}

func (dc *DownloadCommand) PositionalArgs() snek.PositionalArgs {
	return snek.ExactArgs(2)
}

func (dc *DownloadCommand) Flags() []*snek.FlagSet {
	communication := snek.NewFlagSet("Communication")
	communication.DurationP(&dc.Timeout, "timeout", "t", 0, "timeout for the download")
	communication.BoolP(&dc.Force, "force", "f", false, "force download even if the package is already cached")
	communication.StringP(&dc.Registry, "registry", "r", "https://packages.simplifier.net", "registry to download the package from")
	communication.StringP(&dc.AuthToken, "auth-token", "T", "", "auth token for the registry")
	communication.IntP(&dc.Parallel, "parallel", "p", runtime.NumCPU(), "number of parallel downloads")

	output := snek.NewFlagSet("Output")
	output.String(&dc.CacheDir, "fhir-cache", "", "directory to store the downloaded packages")
	output.BoolP(&dc.Verbose, "verbose", "v", false, "enable verbose output")
	output.Bool(&dc.ExcludeDependencies, "exclude-dependencies", false, "include dependencies when downloading the package")
	output.Bool(&dc.NoProgress, "no-progress", false, "disable progress bar")
	output.String(&dc.Log, "log", "", "log file to write the download progress to")
	return []*snek.FlagSet{communication, output}
}

var _ snek.Command = (*DownloadCommand)(nil)
