# Config

> **Note:** This doc is a stub. It will be expanded soon!

Fhenix projects are defined by a configuration `yaml` file, along with a number
of [`text/template`](https://go.dev/pkg/text/template) files. The configuration
file is used to define the project's inputs and teh templates used for
generating the outputs, along with the path at which the outputs will live.

## Hosted Schema

The [JSON Schema] definition of this configuration format is hosted in two
accessible locations:

* In the [fhenix repository] under [jsonschema/fhenix-v0.schema.json]
* Hosted on [this docs site](https://friendly-fhir.github.io/fhenix/jsonschema/fhenix-v0.schema.json)

Either location may be used to point other tooling at to validate fhenix config
files. By default, configuration files will use the latter with a schema
validation comment:

```yaml
# yaml-language-server: $schema=https://friendly-fhir.github.io/fhenix/jsonschema/fhenix-v0.schema.json
...
```

[JSON Schema]: https://json-schema.org/
[fhenix repository]: https://github.com/friendly-fhir/fhenix
[jsonschema/fhenix-v0.schema.json]: https://github.com/friendly-fhir/fhenix/blob/master/jsonschema/fhenix-v0.schema.json

## Example

Below is an example configuration file for a project that generates Go code from
the HL7 FHIR R4 Core package.

```yaml
# yaml-language-server: $schema=https://friendly-fhir.github.io/fhenix/jsonschema/fhenix-v0.schema.json

version: 0

# The name and version of the FHIR Package, as found on simplifier.net
package:
  name: hl7.fhir.r4.core
  version: "4.0.1"

default:
  template:
    field-name: 'templates/field-name.go.tmpl'
    type-name:  'templates/type-name.go.tmpl'
    type-def:   'templates/type-def.go.tmpl'

transformations:
  - input:
      type: StructureDefinition
    output: 'r4/core/{{ .Name | snakecase }}.go'
    template: 'templates/structure-def.go.tmpl'
  - input:
      type: CodeSystem
    output: 'r4/core/codes/{{ .Name | snakecase }}.go'
    templates:
      default: 'templates/code-def.go.tmpl'
  - input:
      type: ValueSet
    output: 'r4/core/valueset/{{ .Name | snakecase }}.go'
    templates:
      default: 'templates/value-def.go.tmpl'
```
