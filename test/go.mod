module github.com/friendly-fhir/fhenix/test

go 1.22.3

// friendly-fhir internal dependencies
require github.com/friendly-fhir/fhenix v0.0.0 // replaced with local

// github.com dependencies
require (
	github.com/cucumber/gherkin/go/v26 v26.2.0 // indirect
	github.com/cucumber/godog v0.14.1
	github.com/cucumber/messages/go/v21 v21.0.1 // indirect
	github.com/gofrs/uuid v4.3.1+incompatible // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-memdb v1.3.4 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

// gopkg.in dependencies
require (
	golang.org/x/oauth2 v0.21.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
)

replace github.com/friendly-fhir/fhenix => ../
