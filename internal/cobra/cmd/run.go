package cmd

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/friendly-fhir/fhenix/internal/config"
	"github.com/friendly-fhir/fhenix/internal/fhirig"
	"github.com/friendly-fhir/fhenix/internal/model"
	"github.com/friendly-fhir/fhenix/internal/model/raw"
	"github.com/friendly-fhir/fhenix/internal/template/engine"
	"github.com/spf13/cobra"
)

var Run = &cobra.Command{
	Use:   "run",
	Short: "Run generation",
	Long:  "Run the generation process",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}
		output, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}
		if output == "" {
			output, err = os.Getwd()
			if err != nil {
				return err
			}
			output = filepath.Join(output, "dist")
		}
		if rm, err := cmd.Flags().GetBool("rm"); err == nil && rm {
			os.RemoveAll(output)
		}
		cfg, err := config.FromFile(args[0])
		if err != nil {
			return err
		}
		name, version := cfg.Package.Name, cfg.Package.Version
		pkg := fhirig.NewPackage(name, version)
		cachePath, err := cmd.Flags().GetString("fhirig-cache")
		if err != nil {
			return err
		}
		if cachePath == "" {
			cachePath = fhirig.SystemCacheDir()
		}
		timeout, err := cmd.Flags().GetDuration("timeout")
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(cmd.Context(), timeout)
		defer cancel()
		cache := &fhirig.PackageCache{
			Root:     cachePath,
			Listener: &Listener{},
		}
		entries, err := cache.FetchAndGet(ctx, pkg)
		if err != nil {
			return err
		}
		fhirCache := model.DefaultFHIRCache()
		for _, entry := range entries {
			if strings.HasPrefix(filepath.Base(entry), "StructureDefinition-") {
				sd, err := raw.ReadStructureDefinition(entry)
				if err != nil {
					return err
				}
				fhirCache.AddStructureDefinition(pkg, filepath.Base(entry), sd)
			}
			if strings.HasPrefix(filepath.Base(entry), "CodeSystem-") {
				cs, err := raw.ReadCodeSystem(entry)
				if err != nil {
					return err
				}
				fhirCache.AddCodeSystem(pkg, filepath.Base(entry), cs)
			}
			// if strings.HasPrefix(filepath.Base(entry), "ValueSet-") {
			// 	vs, err := raw.ReadValueSet(entry)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	builder.AddValueSet(pkg, vs)
			// }
		}
		model := model.NewModel(fhirCache)
		if err := model.DefineAllTypes(); err != nil {
			return err
		}

		engine := engine.New(cfg, engine.Output(output))
		if err := engine.Run(model); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	Root.AddCommand(Run)
	flags := Run.Flags()
	flags.StringP("output", "o", "", "The output directory to write the generated code to")
	flags.Bool("rm", false, "Remove all contents from the output directory prior to writing")
	flags.String("fhirig-cache", "", "The configuration path to download the FHIR IGs to")
	flags.Bool("force", false, "Force download of FHIR IGs")
	flags.Duration("timeout", 0, "Timeout for the download")
}
