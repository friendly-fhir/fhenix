package config

import (
	"bytes"

	"github.com/friendly-fhir/fhenix/config/internal/cfg/v1"
	"github.com/friendly-fhir/fhenix/config/internal/opts"
	"gopkg.in/yaml.v3"
)

func fromV1(data []byte, opts *opts.Options) (*Config, error) {
	var cfg cfg.Root
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	var result Config
	// These are mapped 1:1
	if cfg.Mode != "" {
		result.Mode = Mode(cfg.Mode)
	} else {
		result.Mode = ModeText
	}

	var err error
	if opts.OutputDir != "" {
		result.OutputDir, err = opts.RootPath(opts.OutputDir)
	} else if cfg.OutputDir != "" {
		result.OutputDir, err = opts.RootPath(cfg.OutputDir)
	} else {
		result.OutputDir, err = opts.OutputPath("dist")
	}
	if err != nil {
		return nil, err
	}

	base, err := fromV1Transform(opts, &cfg.Default)
	if err != nil {
		return nil, err
	}

	for _, pkg := range cfg.Input.Packages {
		result.Input = append(result.Input, &Package{
			Name:    pkg.Name,
			Version: pkg.Version,
			Path:    pkg.Path,
		})
	}

	result.Transforms = make([]*Transform, len(cfg.Transforms))
	for i := range result.Transforms {
		result.Transforms[i], err = mergeV1Transforms(opts, base, cfg.Transforms[i])
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func fromV1Transform(opts *opts.Options, transform *cfg.Transform) (*Transform, error) {
	var err error
	var result Transform
	result.Funcs = make(map[string]string, len(transform.Funcs))
	for name, path := range transform.Funcs {
		result.Funcs[name], err = opts.RootPath(path)
		if err != nil {
			return nil, err
		}
	}
	if transform.Templates == nil {
		result.Templates = map[string]string{}
	} else {
		result.Templates = make(map[string]string, 6+len(transform.Templates.Partials))
		entries := []struct {
			name   string
			member string
		}{
			{"main", transform.Templates.Main},
			{"header", transform.Templates.Header},
			{"footer", transform.Templates.Footer},
			{"code-system", transform.Templates.CodeSystem},
			{"value-set", transform.Templates.ValueSet},
			{"type", transform.Templates.Type},
		}

		for _, entry := range entries {
			if entry.member == "" {
				continue
			}
			result.Templates[entry.name], err = opts.RootPath(entry.member)
			if err != nil {
				return nil, err
			}
		}
		for name, path := range transform.Templates.Partials {
			result.Templates[name], err = opts.RootPath(path)
			if err != nil {
				return nil, err
			}
		}
	}

	result.Include = fromV1Filters(transform.Include)
	result.Exclude = fromV1Filters(transform.Exclude)
	result.OutputPath = transform.OutputPath

	return &result, nil
}

func mergeV1Transforms(opts *opts.Options, base *Transform, transformv1 *cfg.Transform) (*Transform, error) {
	result, err := fromV1Transform(opts, transformv1)
	if err != nil {
		return nil, err
	}
	if result.OutputPath == "" {
		result.OutputPath = base.OutputPath
	}

	for name, path := range base.Funcs {
		if _, ok := result.Funcs[name]; !ok {
			result.Funcs[name] = path
		}
	}
	for name, path := range base.Templates {
		if _, ok := result.Templates[name]; !ok {
			result.Templates[name] = path
		}
	}

	result.Include = append(result.Include, base.Include...)
	result.Exclude = append(result.Exclude, base.Exclude...)
	return result, nil
}

func fromV1Filters(filters []*cfg.TransformFilter) []*TransformFilter {
	result := make([]*TransformFilter, len(filters))
	for i, filter := range filters {
		result[i] = fromV1Filter(filter)
	}
	return result
}

func fromV1Filter(filter *cfg.TransformFilter) *TransformFilter {
	return &TransformFilter{
		Name:      filter.Name,
		Type:      filter.Type,
		URL:       filter.URL,
		Package:   filter.Package,
		Source:    filter.Source,
		Condition: filter.Condition,
	}
}
