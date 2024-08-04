# ![Fhenix](./docs/phenix-logo.png)

![Continuous Integration](https://img.shields.io/github/actions/workflow/status/friendly-fhir/fhenix/.github%2Fworkflows%2Fcontinuous-integration.yaml?logo=github)
[![GitHub Release](https://img.shields.io/github/v/release/friendly-fhir/fhenix?include_prereleases)][github-releases]
[![Gitter Channel](https://img.shields.io/badge/matrix-%23friendly--fhir-darkcyan?logo=gitter)][gitter-channel]
[![readthedocs](https://img.shields.io/badge/docs-readthedocs-blue?logo=readthedocs&logoColor=white)][docs]
[![Godocs](https://img.shields.io/badge/docs-reference-blue?logo=go&logoColor=white)][go-docs]

Fhenix is a flexible and lightweight tool for generating content from a modeling
of a [FHIR] IG's definitional entries.

This leverages packages as defined in the [Simplifier registry], parses their
relevant entities ([StructureDefinition], [CodeSystem], and
[ValueSet]), constructs a model of these entities, and then feeds it into Go
templates to generate content.

Check out our [examples](./examples)!

[FHIR]: https://www.hl7.org/fhir/
[ValueSet]: https://www.hl7.org/fhir/valueset.html
[CodeSystem]: https://www.hl7.org/fhir/codesystem.html
[StructureDefinition]: https://www.hl7.org/fhir/structuredefinition.html

[gitter-channel]: https://matrix.to/#/#friendly-fhir:gitter.im
[docs]: https://friendly-fhir.github.io/fhenix/
[go-docs]: https://pkg.go.dev/github.com/friendly-fhir/fhenix
[github-releases]: https://github.com/friendly-fhir/fhenix/releases
[Simplifier registry]: https://simplifier.net

## Quick Start

* [📚 Getting Starts][docs]
* [🚂 Examples](./examples)
* [🚀 Use](#use)
* [📦 Getting Started][getting_started]

[getting_started]: https://friendly-fhir.org/fhenix/user-manual/getting-started/

## Use

Fhenix generates content from [FHIR Conformance Modules] defined in registries
like [Simplifier].

With Fhenix, you can do:

* Code generation, to create custom bindings and logic for handling your
  organization's profiles
* Documentation generation, to create custom documentation for your IGs
* Data generation, to create test data based on profiles
* Schema generation, to create SQL tables of FHIR resource data

Some practical projects leveraging this within [Friendly FHIR]:

* [go-fhir](https://github.com/friendly-fhir/go-fhir):
  A Go library for working with FHIR resources
* [rust-fhir](https://github.com/friendly-fhir/rust-fhir):
  A Rust library for working with FHIR resources
* [sqlite-fhir-server](https://github.com/friendly-fhir/sqlite-fhirstore)
  A FHIR server that uses [SQLite] as a backend for storage and emitting
  REST APIs. Generates SQLite tables from `StructureDefinition` resources.

[SQLite]: https://sqlite.org/index.html
[Friendly FHIR]: https://github.com/friendly-fhir
[FHIR Conformance Modules]: https://www.hl7.org/fhir/conformance-module.html
