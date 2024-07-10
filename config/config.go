package config

// Mode is the type of template system being used for the output.
type Mode string

const (
	// ModeText is the mode for text-based output.
	ModeText Mode = "text"

	// ModeHTML is the mode for HTML-based output.
	ModeHTML Mode = "html"
)

// Config represents the configuration that is used for generation within Fhenix.
// This is not a "raw" representation of the config, but rather an object
// representation capable of performing actions independently of the underlying
// configuration version.
type Config struct {
	// Mode is the mode type of template system being used for the output.
	Mode Mode

	// OutputDir is the output directory where generated output will be written.
	OutputDir string

	// Input is a package input.
	Input *Package

	// Transforms contains a list of transforms to be applied to the input
	// definitions.
	Transforms []*Transform
}

// Package is the source package that will be used as input for the generation.
type Package struct {
	// Name is the name of the package (mandatory).
	Name string

	// Version is a version string for the package version (mandatory).
	Version string

	// Path is an optional path to specify to where the package is located.
	// If specified, this will override the package being fetched from the
	// package registry.
	Path string

	// IncludeDependencies is a flag to indicate whether dependencies of the
	// package should be included in the generation.
	IncludeDependencies bool
}

// Transform is a configuration for transforming input entities into templated
// output(s).
type Transform struct {
	// Include is a list of filters for conditions that an entity may satisfy to
	// be included in the transformation. At least one of these filters must be
	// satisfied for an entity to be included.
	Include []*TransformFilter

	// Exclude is a list of filters for conditions that an entity may satisfy to
	// be excluded from the transformation. If any of these filters are satisfied,
	// the entity will be excluded, even if it satisfies an 'include' filter.
	Exclude []*TransformFilter

	// OutputPath is a template for where the transformation will be output to.
	// The output path is always dynamic, and will be fed each individual entity
	// that has been matched by the filters to be transformed.
	OutputPath string

	// Funcs is mapping of template function-names to files containing templates
	// that will perform textual transformations. This enables more complex logic
	// to be included in template pipelines.
	Funcs map[string]string

	// Templates is a mapping of template names to files containing
	// templates that will be included in the main templates. This enables
	// re-use of common template snippets.
	//
	// Some special templates have implicit meaning in this configuration, but
	// this is otherwise freeform:
	//
	// - 'header': This template will be included at the top of the output file.
	// - 'footer': This template will be included at the bottom of the output file.
	// - 'type': This template will be called with each _individual entity_ 'type'
	//    that is matched by the filters.
	// - 'code-system': This template will be called with each _individual entity_
	//   'code-system' that is matched by the filters.
	// - 'main': This template will be called by the _list of all matched entities_.
	//   This template is provided by default by the implementation, which will call
	//   'header', followed by the appropriate intermediate template, followed by
	//   'footer'. Replacing this will replace all the above templates.
	Templates map[string]string
}

type TransformFilter struct {
	// Name is a filter on the name of the input entity.
	// This may be a regular expression.
	Name string

	// Type is the type of the input entity.
	Type string

	// URL is an exact-match filter on the URL of the input entity.
	URL string

	// Package is a filter on the package-name of the input entity.
	// This is only relevant if 'input.include-dependencies' is enabled.
	Package string

	// Source is a filter on the source-file that the input entity is defined in.
	// This may be a regular expression.
	Source string

	// Condition is a custom, template-filter that the input entity must satisfy.
	// This is meant as a back-door to allow for more complex conditions than
	// the other filters allow.
	Condition string
}
