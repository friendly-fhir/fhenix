# Getting Started

## Installation

Fhenix can be installed in a number of ways that are outlined in this section.
The exact relevant section will depend on system you are running and which
package managers you may have available.

### System-Specific Installation

#### GitHub Actions

Leverage the [fhenix-actions] repository, which can do all of this for you
automatically.

[fhenix-actions]: https://github.com/friendly-fhir/fhenix-actions

#### macOS

##### Homebrew

`fhenix` has a Homebrew formula defined in [friendly-fhir/homebrew-tap]. You can
install via `brew` either by tapping the repository and then installing, or by
installing directly from the URL.

```bash
# Tapping first
brew tap friendly-fhir/homebrew-tap # This is only needed once

brew install fhenix

# Directly
brew install friendly-fhir/homebrew-tap/fhenix
```

[friendly-fhir/homebrew-tap]: https://github.com/friendly-fhir/homebrew-tap

### System-Agnostic Installation

#### `install.sh`

The `install.sh` script is a simple script that will download the latest release
of Fhenix and install it to your system. This is the recommended way to install
Fhenix if you are not using Go.

```bash
curl -sSL https://raw.githubusercontent.com/friendly-fhir/fhenix/master/install.sh | bash
```

This will fetch the latest official release of Fhenix, and install it into
`./bin` in your current working directory, unless `BINDIR` is specified -- in
which case it will be installed into that directory.

#### `go install`

Installation can be done via a normal `go install`:

```bash
go install github.com/friendly-fhir/fhenix@v0.1.0
```

Or for the latest (development) version:

```bash
go install github.com/friendly-fhir/fhenix@latest
```

**Note:** `${GOPATH}/bin` must be in your `PATH` for this to work. In
Posix-compliance systems, this can be done with:

```bash
export PATH="${PATH}:$(go env GOPATH)/bin"
```
