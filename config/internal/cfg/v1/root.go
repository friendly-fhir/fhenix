package cfg

// Root is the root configuration node of the YAML file.
type Root struct {
	// Version is the schema version of the configuration file.
	// This will always be '1'
	Version int `yaml:"version"`

	// Mode is the mode type of template output that will be used for the
	// generated output. May be one of 'text' or 'html'.
	Mode Mode `yaml:"mode"`

	// RootDir is the root directory of the project.
	//
	// Relative paths will be translated to coherent roots in the 'config' package,
	// by assuming the root is relative to where the configuration file was
	// parsed (if specified), or relative to the current working directory if not.
	RootDir string `yaml:"root-dir"`

	// OutputDir is the directory where the generated output will be written.
	//
	// Relative paths will be translated to coherent roots in the 'config' package,
	// by assuming the root is relative to where the configuration file was
	// parsed (if specified), or relative to the current working directory if not.
	OutputDir string `yaml:"output-dir"`

	// Default is a configuration node that specifies default values to use for
	// transforms. This just helps to reduce the boilerplate when several
	// transformations use the same set of templates.
	Default Transform `yaml:"default"`

	// Transforms is a list of transformations to apply to the input data.
	Transforms []*Transform `yaml:"transforms"`
}
