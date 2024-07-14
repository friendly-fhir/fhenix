package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/friendly-fhir/fhenix/registry"
	"github.com/spf13/cobra"
)

var DownloadFlags struct {
	Timeout time.Duration
	Force   bool
	Verbose bool

	CacheDir  string
	Registry  string
	AuthToken string
}

type downloadListener struct {
	verbose bool
	registry.BaseCacheListener
}

func (l *downloadListener) BeforeFetch(registry, pkg, version string) {
	fmt.Printf("connecting to %s@%s (from %s)\n", pkg, version, registry)
}

func (l *downloadListener) OnFetch(registry, pkg, version string, data int64) {
	fmt.Printf("downloading %s@%s (from %s): %d bytes\n", pkg, version, registry, data)
}

func (l *downloadListener) OnFetchWrite(registry, pkg, version string, data []byte) {
	if l.verbose {
		fmt.Printf("* [%s@%s] (from %s): %d bytes...\n", pkg, version, registry, len(data))
	}
}

func (l *downloadListener) AfterFetch(registry, pkg, version string, err error) {
	if err != nil {
		fmt.Printf("download error: %v\n", err)
	} else {
		fmt.Printf("downloaded %s@%s (from %s)\n", pkg, version, registry)
	}
}

func (l *downloadListener) OnCacheHit(registry, pkg, version string) {
	fmt.Printf("cache-hit: %s@%s (from %s)\n", pkg, version, registry)
}

func (l *downloadListener) OnUnpack(registry, pkg, version, file string, data int64) {
	if l.verbose {
		fmt.Printf("[%s@%s] %s (from %s): %d bytes\n", pkg, version, file, registry, data)
	}
}

var _ registry.CacheListener = (*downloadListener)(nil)

var Download = &cobra.Command{
	Use:   "download <package> <version>",
	Short: "Download FHIR IGs",
	Long:  "Download FHIR Implementation Guides (IGs) from the web",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return cmd.Usage()
		}
		pkg, version := args[0], args[1]

		ctx := context.Background()
		if DownloadFlags.Timeout != 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(cmd.Context(), DownloadFlags.Timeout)
			defer cancel()
		}
		var cache *registry.Cache
		if dir := DownloadFlags.CacheDir; dir != "" {
			cache = registry.NewCache(dir)
		} else {
			cache = registry.DefaultCache()
		}
		listener := &downloadListener{
			verbose: DownloadFlags.Verbose,
		}
		cache.AddListener(listener)

		var opts []registry.Option
		opts = append(opts, registry.URL(DownloadFlags.Registry))
		if token := DownloadFlags.AuthToken; token != "" {
			opts = append(opts, registry.Auth(registry.StaticTokenSource(token)))
		}

		client, err := registry.NewClient(ctx, opts...)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "error: %v\n", err)
			os.Exit(1)
		}
		cache.AddClient("default", client)

		fetch := cache.Fetch
		if DownloadFlags.Force {
			fetch = cache.ForceFetch
		}
		if err := fetch(ctx, "default", pkg, version); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "error: %v\n", err)
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	Root.AddCommand(Download)
	flags := Download.Flags()
	flags.DurationVarP(&DownloadFlags.Timeout, "timeout", "t", 0, "timeout for the download")
	flags.BoolVarP(&DownloadFlags.Force, "force", "f", false, "force download even if the package is already cached")
	flags.StringVar(&DownloadFlags.CacheDir, "fhir-cache", "", "directory to store the downloaded packages")
	flags.StringVarP(&DownloadFlags.Registry, "registry", "r", "https://packages.simplifier.net", "registry to download the package from")
	flags.StringVarP(&DownloadFlags.AuthToken, "auth-token", "T", "", "auth token for the registry")
	flags.BoolVarP(&DownloadFlags.Verbose, "verbose", "v", false, "enable verbose output")
}
