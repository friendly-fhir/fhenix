package config

import "github.com/friendly-fhir/fhenix/config/internal/opts"

// WithRootDir sets the root directory which all configuration paths will be
// considered relative to.
func WithRootDir(dir string) Option {
	return opts.OptionFunc(func(cfg *opts.Options) {
		cfg.RootDir = dir
	})
}

// WithOutputDir sets the output directory where generated output will be
// written.
func WithOutputDir(dir string) Option {
	return opts.OptionFunc(func(cfg *opts.Options) {
		cfg.OutputDir = dir
	})
}

// Option is an interface for composable optins that can be provided to
// configuration objects that can be read.
type Option = opts.Option
