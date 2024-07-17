package cmd

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/friendly-fhir/fhenix/driver"
	"github.com/friendly-fhir/fhenix/internal/snek"
)

func lines(lines ...string) string {
	return strings.Join(lines, "\n")
}

type Listener struct {
	verbose bool
	out     io.Writer
	driver.BaseListener
}

func NewDriverListener(ctx context.Context, verbose bool) *Listener {
	return &Listener{
		out: snek.CommandOut(ctx),
	}
}

func (l *Listener) BeforeDownload() {
	fmt.Fprintln(l.out, "[1] Downloading FHIR Packages...")
}

func (l *Listener) AfterDownload(err error) {
	if err != nil {
		fmt.Fprintf(l.out, "[1] Downloading FHIR Packages... Failed: %v\n", err)
	} else {
		fmt.Fprintln(l.out, "[1] Downloading FHIR Packages... Succeeded")
	}
}

func (l *Listener) OnUnpack(registry, pkg, version, file string, data int64) {
	if l.verbose {
		fmt.Printf("[%s@%s] %s (from %s): %d bytes\n", pkg, version, file, registry, data)
	}
}
func (l *Listener) BeforeLoadTransform() {
	fmt.Fprintln(l.out, "[4] Loading Transformations...")
}

func (l *Listener) AfterLoadTransform(err error) {
	if err != nil {
		fmt.Fprintf(l.out, "[4] Loading Transformations... Failed: %v\n", err)
	} else {
		fmt.Fprintln(l.out, "[4] Loading Transformations... Succeeded")
	}
}

func (l *Listener) BeforeLoadConformance() {
	fmt.Fprintln(l.out, "[2] Loading Conformance...")
}

func (l *Listener) AfterLoadConformance(err error) {
	if err != nil {
		fmt.Fprintf(l.out, "[2] Loading Conformance... Failed: %v\n", err)
	} else {
		fmt.Fprintln(l.out, "[2] Loading Conformance... Succeeded")
	}
}

func (l *Listener) BeforeLoadModel() {
	fmt.Fprintln(l.out, "[3] Loading Model...")
}

func (l *Listener) AfterLoadModel(err error) {
	if err != nil {
		fmt.Fprintf(l.out, "[3] Loading Model... Failed: %v\n", err)
	} else {
		fmt.Fprintln(l.out, "[3] Loading Model... Succeeded")
	}
}

func (l *Listener) BeforeTransform() {
	fmt.Fprintln(l.out, "[5] Transforming Model...")
}

func (l *Listener) OnTransform(output string) {
	fmt.Fprintf(l.out, "[5] Transforming Model... %s\n", output)
}

func (l *Listener) AfterTransform(jobs int, err error) {
	if err != nil {
		fmt.Fprintf(l.out, "[5] Transforming Model... Failed: %v\n", err)
	} else {
		fmt.Fprintln(l.out, "[5] Transforming Model... Succeeded")
	}
}

func (l *Listener) BeforeFetch(registry, pkg, version string) {
	fmt.Printf("connecting to %s@%s (from %s)\n", pkg, version, registry)
}

func (l *Listener) OnFetch(registry, pkg, version string, data int64) {
	fmt.Printf("downloading %s@%s (from %s): %d bytes\n", pkg, version, registry, data)
}

func (l *Listener) OnFetchWrite(registry, pkg, version string, data []byte) {
	if l.verbose {
		fmt.Printf("* [%s@%s] (from %s): %d bytes...\n", pkg, version, registry, len(data))
	}
}

func (l *Listener) AfterFetch(registry, pkg, version string, err error) {
	if err != nil {
		fmt.Printf("download error: %v\n", err)
	} else {
		fmt.Printf("downloaded %s@%s (from %s)\n", pkg, version, registry)
	}
}

func (l *Listener) OnCacheHit(registry, pkg, version string) {
	fmt.Printf("cache-hit: %s@%s (from %s)\n", pkg, version, registry)
}

var _ driver.Listener = (*Listener)(nil)
