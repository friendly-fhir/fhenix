package cmd

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/friendly-fhir/fhenix/driver"
)

type Listener struct {
	verbose bool
	out     *log.Logger
	driver.BaseListener
}

func NewLogListener(w io.Writer, verbose bool) *Listener {
	flags := 0
	if verbose {
		flags = log.LstdFlags
	}
	return &Listener{
		out: log.New(w, "", flags),
	}
}

func (l *Listener) BeforeDownload() {
	l.out.Println("[1] Downloading FHIR Packages...")
}

func (l *Listener) AfterDownload(err error) {
	if err != nil {
		l.out.Printf("[1] Downloading FHIR Packages... Failed: %v", err)
	} else {
		l.out.Println("[1] Downloading FHIR Packages... Succeeded")
	}
}

func (l *Listener) OnUnpack(registry, pkg, version, file string, data int64) {
	if l.verbose {
		l.out.Printf("[%s@%s] %s (from %s): %d bytes", pkg, version, file, registry, data)
	}
}
func (l *Listener) BeforeLoadTransform() {
	l.out.Println("[4] Loading Transformations...")
}

func (l *Listener) AfterLoadTransform(err error) {
	if err != nil {
		l.out.Printf("[4] Loading Transformations... Failed: %v", err)
	} else {
		l.out.Println("[4] Loading Transformations... Succeeded")
	}
}

func (l *Listener) BeforeLoadConformance() {
	l.out.Println("[2] Loading Conformance...")
}

func (l *Listener) AfterLoadConformance(err error) {
	if err != nil {
		l.out.Printf("[2] Loading Conformance... Failed: %v", err)
	} else {
		l.out.Println("[2] Loading Conformance... Succeeded")
	}
}

func (l *Listener) BeforeLoadModel() {
	l.out.Println("[3] Loading Model...")
}

func (l *Listener) AfterLoadModel(err error) {
	if err != nil {
		l.out.Printf("[3] Loading Model... Failed: %v", err)
	} else {
		l.out.Println("[3] Loading Model... Succeeded")
	}
}

func (l *Listener) BeforeTransformStage() {
	l.out.Println("[5] Transforming Model...")
}

func (l *Listener) OnTransformStage(output string) {
	l.out.Printf("[5] Transforming Model... %s", output)
}

func (l *Listener) AfterTransformStage(jobs int, err error) {
	if err != nil {
		l.out.Printf("[5] Transforming Model... Failed: %v", err)
	} else {
		l.out.Println("[5] Transforming Model... Succeeded")
	}
}

func (l *Listener) BeforeFetch(registry, pkg, version string) {
	l.out.Printf("%s@%s: checking for package", pkg, version)
}

func (l *Listener) OnFetch(registry, pkg, version string, data int64) {
	l.out.Printf("%s@%s: reading %d bytes from %s registry", pkg, version, data, registry)
}

func (l *Listener) OnFetchWrite(registry, pkg, version string, data []byte) {
	if l.verbose {
		l.out.Printf("%s@%s: %d bytes...", pkg, version, len(data))
	}
}

func (l *Listener) AfterFetch(registry, pkg, version string, err error) {
	if err != nil {
		l.out.Printf("%s@%s: error: %v", pkg, version, err)
	} else {
		l.out.Printf("%s@%s: fetch complete", pkg, version)
	}
}

func (l *Listener) OnCacheHit(registry, pkg, version string) {
	l.out.Printf("%s@%s in cache", pkg, version)
}

func (l *Listener) OnTransformOutput(n int, output string) {
	if l.verbose {
		l.out.Printf("transform(%d): %s -- created", n, l.shortPath(output))
	}
}

func (l *Listener) BeforeTransform(n int) {
	if l.verbose {
		l.out.Printf("transform(%d): starting", n)
	}
}

func (l *Listener) AfterTransformOutput(n int, output string, err error) {
	if err != nil {
		l.out.Printf("transform(%d): %v -- error: %v", n, l.shortPath(output), err)
	} else {
		l.out.Printf("transform(%d): %v -- generated", n, l.shortPath(output))
	}
}

func (l *Listener) shortPath(output string) string {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	path, err := filepath.Rel(cwd, output)
	if err != nil {
		path = output
	}
	return path
}

var _ driver.Listener = (*Listener)(nil)
