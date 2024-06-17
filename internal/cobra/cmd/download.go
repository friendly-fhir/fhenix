package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/friendly-fhir/fhenix/internal/fhirig"
	"github.com/spf13/cobra"
)

type Listener struct {
	fhirig.BaseListener
}

func (*Listener) OnFetchStart(pkg *fhirig.Package) {
	fmt.Printf("Fetching package %s\n", pkg.String())
}
func (*Listener) OnFetchEnd(pkg *fhirig.Package, err error) {
	if err != nil {
		fmt.Printf("Package %s failed to be fetched: %v\n", pkg.String(), err)
		return
	}
	fmt.Printf("Package %s fetched successfully\n", pkg.String())
}
func (*Listener) OnCacheHit(pkg *fhirig.Package) {
	fmt.Printf("Package %s already downloaded\n", pkg.String())
}

var Download = &cobra.Command{
	Use:   "download <package> <version>",
	Short: "Download FHIR IGs",
	Long:  "Download FHIR Implementation Guides (IGs) from the web",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return cmd.Usage()
		}
		pkg, version := args[0], args[1]

		timeout, err := cmd.Flags().GetDuration("timeout")
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(cmd.Context(), timeout)
		defer cancel()

		dir, err := cmd.Flags().GetString("fhirig-cache")
		if err != nil {
			return err
		}
		if strings.TrimSpace(dir) == "" {
			dir = fhirig.SystemCacheDir()
		}

		cache := &fhirig.PackageCache{
			Root:     dir,
			Listener: &Listener{},
		}
		fetch := cache.Fetch
		if force, err := cmd.Flags().GetBool("force"); err == nil && force {
			fetch = cache.ForceFetch
		}

		if err := fetch(ctx, fhirig.NewPackage(pkg, version)); err != nil {
			return fmt.Errorf("fetching %v@%v: %w", pkg, version, err)
		}
		return nil
	},
}

func init() {
	Root.AddCommand(Download)
	flags := Download.Flags()
	flags.DurationP("timeout", "t", 0, "timeout for the download")
	flags.BoolP("force", "f", false, "force download even if the package is already cached")
}
