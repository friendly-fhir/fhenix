module github.com/friendly-fhir/fhenix

go 1.22.3

// friendly-fhir internal dependencies
require github.com/friendly-fhir/go-fhir v0.0.0-20240627035249-eacfb3386af5

// golang.org official dependencies
require (
	golang.org/x/oauth2 v0.21.0
	golang.org/x/sync v0.7.0
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0
)

// github.com dependencies
require (
	github.com/google/go-cmp v0.6.0
	github.com/iancoleman/strcase v0.3.0
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/cobra v1.8.0
	github.com/spf13/pflag v1.0.5
)

// gopkg.in dependencies
require (
	golang.org/x/term v0.22.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1
)

// atomicgo.dev dependencies
require atomicgo.dev/cursor v0.2.0
