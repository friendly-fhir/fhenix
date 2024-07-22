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

func (l *Listener) OnUnpack(registry, pkg, version, file string, data int64) {
	if l.verbose {
		l.out.Printf("[%s@%s] %s (from %s): %d bytes", pkg, version, file, registry, data)
	}
}

func (l *Listener) BeforeStage(s driver.Stage) {
	switch s {
	case driver.StageDownload:
		l.out.Println("[1] Downloading FHIR Packages...")
	case driver.StageLoadConformance:
		l.out.Println("[2] Loading Conformance...")
	case driver.StageLoadModel:
		l.out.Println("[3] Loading Model...")
	case driver.StageLoadTransform:
		l.out.Println("[4] Loading Transformations...")
	case driver.StageTransform:
		l.out.Println("[5] Transforming Model...")
	}
}

func (l *Listener) AfterStage(s driver.Stage, err error) {
	suffix := "succeeded"
	if err != nil {
		suffix = "failed: " + err.Error()
	}
	switch s {
	case driver.StageDownload:
		l.out.Printf("[1] Downloading FHIR Packages... %s", suffix)
	case driver.StageLoadConformance:
		l.out.Printf("[2] Loading Conformance... %s", suffix)
	case driver.StageLoadModel:
		l.out.Printf("[3] Loading Model... %s", suffix)
	case driver.StageLoadTransform:
		l.out.Printf("[4] Loading Transformations... %s", suffix)
	case driver.StageTransform:
		l.out.Printf("[5] Transforming Model... %s", suffix)
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
