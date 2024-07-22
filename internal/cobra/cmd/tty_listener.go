package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/friendly-fhir/fhenix/data"
	"github.com/friendly-fhir/fhenix/driver"
	"github.com/friendly-fhir/fhenix/internal/ansi"
	"github.com/friendly-fhir/fhenix/internal/snek/spinner"
	"github.com/friendly-fhir/fhenix/internal/snek/terminal"
	"github.com/friendly-fhir/fhenix/registry"
)

type TTYListener struct {
	offset     int
	stage      int
	downloads  map[string]*download
	loaderPkgs map[string]*loadPackage
	transforms map[int]*transform

	m sync.Mutex

	terminal *terminal.Terminal
	spinner  *spinner.Spinner

	verbose bool
	driver.BaseListener
}

func NewProgressListener(term *terminal.Terminal, verbose bool) *TTYListener {
	return &TTYListener{
		verbose:  verbose,
		terminal: term,
		spinner:  spinner.Dots(time.Second),
	}
}

func keyOf(registry, pkg, version string) string {
	return fmt.Sprintf("%v::%s@%s", registry, pkg, version)
}

func (l *TTYListener) download(key string) *download {
	if l.downloads == nil {
		l.downloads = make(map[string]*download)
	}
	d, ok := l.downloads[key]
	if !ok {
		offset := l.offset
		l.offset++

		d = &download{
			Line:    l.terminal.Line(offset),
			Spinner: l.spinner.Clone(),
		}
		l.downloads[key] = d
	}
	return d
}

func (l *TTYListener) transform(n int) *transform {
	if l.transforms == nil {
		l.transforms = make(map[int]*transform)
	}
	t, ok := l.transforms[n]
	if !ok {
		offset := l.offset
		l.offset++

		t = &transform{
			Line: l.terminal.Line(offset),
		}
		l.transforms[n] = t
	}
	return t
}

func (l *TTYListener) loadPackage(key string) *loadPackage {
	if l.loaderPkgs == nil {
		l.loaderPkgs = make(map[string]*loadPackage)
	}
	p, ok := l.loaderPkgs[key]
	if !ok {
		offset := l.offset
		l.offset++

		p = &loadPackage{
			Line: l.terminal.Line(offset),
		}
		l.loaderPkgs[key] = p
	}
	return p
}

func (l *TTYListener) BeforeStage(s driver.Stage) {
	l.m.Lock()
	defer l.m.Unlock()

	offset := l.offset
	l.offset++
	line := l.terminal.Line(offset)
	l.stage++
	stage := l.stage

	prefix := fmt.Sprintf("[%d/5]", stage)
	switch s {
	case driver.StageDownload:
		line.Printf("%s %s\n", prefix, "Downloading FHIR packages")
	case driver.StageLoadTransform:
		line.Printf("%s %s\n", prefix, "Loading transformations")
	case driver.StageLoadConformance:
		line.Printf("%s %s\n", prefix, "Loading conformance module")
	case driver.StageLoadModel:
		line.Printf("%s %s\n", prefix, "Loading model")
	case driver.StageTransform:
		line.Printf("%s %s\n", prefix, "Transforming Outputs")
	}
}

func (l *TTYListener) BeforeFetch(registry, pkg, version string) {
	key := keyOf(registry, pkg, version)
	l.m.Lock()
	defer l.m.Unlock()
	download := l.download(key)

	name := fmt.Sprintf("%s%s%s", ansi.FGBrightWhite.Format(pkg), ansi.FGGray.Format("@"), version)

	download.Line.Print(l.progress(name, 0, 0, nil))
}

func (l *TTYListener) OnFetch(registry, pkg, version string, data int64) {
	key := keyOf(registry, pkg, version)
	l.m.Lock()
	defer l.m.Unlock()
	download := l.download(key)
	download.TotalBytes = data

	name := fmt.Sprintf("%s%s%s", ansi.FGBrightWhite.Format(pkg), ansi.FGGray.Format("@"), version)

	download.Line.Print(l.progress(name, 0, data, nil))
}

func (l *TTYListener) OnFetchWrite(registry, pkg, version string, data []byte) {
	key := keyOf(registry, pkg, version)
	l.m.Lock()
	defer l.m.Unlock()
	download := l.download(key)
	download.Current += int64(len(data))

	name := fmt.Sprintf("%s%s%s", ansi.FGBrightWhite.Format(pkg), ansi.FGGray.Format("@"), version)

	download.Line.Print(l.progress(name, download.Current, download.TotalBytes, nil))
}

func (l *TTYListener) AfterFetch(registry, pkg, version string, err error) {
	key := keyOf(registry, pkg, version)
	l.m.Lock()
	defer l.m.Unlock()
	download := l.download(key)

	name := fmt.Sprintf("%s%s%s", ansi.FGBrightWhite.Format(pkg), ansi.FGGray.Format("@"), version)

	if err != nil {
		download.Line.Print(l.valueProgress(ansi.FGRed.Format("x"), name, "error"))
	} else {
		download.Line.Print(l.progress(name, download.Current, download.TotalBytes, nil))
	}
}

func (l *TTYListener) OnCacheHit(registry, pkg, version string) {
	key := keyOf(registry, pkg, version)
	l.m.Lock()
	defer l.m.Unlock()
	if l.downloads == nil {
		l.downloads = make(map[string]*download)
	}
	download := l.download(key)
	if download.Current != 0 {
		return
	}
	name := fmt.Sprintf("%s%s%s", ansi.FGBrightWhite.Format(pkg), ansi.FGGray.Format("@"), version)

	download.Line.Print(l.valueProgress(ansi.FGYellow.Format("✓"), name, "cache"))
}

func (l *TTYListener) BeforeLoadPackage(ref registry.PackageRef) {
	l.m.Lock()
	defer l.m.Unlock()

	pkg := l.loadPackage(ref.String())
	name := fmt.Sprintf("%s%s%s", ansi.FGBrightWhite.Format(ref.Name()), ansi.FGGray.Format("@"), ref.Version())
	pkg.Line.Print(l.valueProgress(ansi.FGYellow.Format("-"), name, "loading"))
}

func (l *TTYListener) AfterLoadPackage(ref registry.PackageRef, err error) {
	l.m.Lock()
	defer l.m.Unlock()

	name := fmt.Sprintf("%s%s%s", ansi.FGBrightWhite.Format(ref.Name()), ansi.FGGray.Format("@"), ref.Version())
	pkg := l.loadPackage(ref.String())
	if err != nil {
		pkg.Line.Println(l.valueProgress(ansi.FGRed.Format("x"), name, "error"))
	} else {
		pkg.Line.Println(l.valueProgress(ansi.FGGreen.Format("✓"), name, "loaded"))
	}
}

func (l *TTYListener) BeforeTransform(i int) {
	l.m.Lock()
	defer l.m.Unlock()

	content := fmt.Sprintf("transform %d", i)
	transform := l.transform(i)
	transform.Line.Print(l.valueProgress(ansi.FGYellow.Format("-"), content, "transforming"))
}

func (l *TTYListener) OnTransform(i int, output string) {
	l.m.Lock()
	defer l.m.Unlock()

	content := l.transformPrefix(i, output)
	transform := l.transform(i)
	transform.Line.Print(l.valueProgress(ansi.FGYellow.Format("-"), content, "transforming"))
}

func (l *TTYListener) transformPrefix(i int, output string) string {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	path, err := filepath.Rel(cwd, output)
	if err != nil {
		path = output
	}
	return fmt.Sprintf("transform %d (%s)", i, ansi.FGGray.Format(path))
}

func (l *TTYListener) AfterTransformOutput(i int, output string, err error) {
	l.m.Lock()
	defer l.m.Unlock()

	content := l.transformPrefix(i, output)
	transform := l.transform(i)
	if err != nil {
		transform.Line.Print(l.valueProgress(ansi.FGRed.Format("x"), content, "error"))
	} else {
		transform.Line.Print(l.valueProgress(ansi.FGGreen.Format("✓"), content, "done"))
	}
}

var _ driver.Listener = (*TTYListener)(nil)

type download struct {
	TotalBytes int64
	Current    int64
	Line       *terminal.Line
	Spinner    *spinner.Spinner
}

type loadPackage struct {
	Line *terminal.Line
}

type transform struct {
	Line *terminal.Line
}

func (l *TTYListener) valueProgress(state, name, suffix string) string {
	var (
		sb strings.Builder
	)
	const (
		minLength = 59
		maxLength = 119
	)
	length := max(min(l.terminal.Width()-1, maxLength), minLength)
	fmt.Fprintf(&sb, "[%s] %s", state, name)
	prefixRunes := utf8.RuneCountInString(ansi.StripFormat(sb.String()))
	suffixRunes := utf8.RuneCountInString(ansi.StripFormat(suffix))
	reserved := prefixRunes + suffixRunes
	dots := length - reserved
	if reserved < 0 {
		dots = 0
	}

	if dots != 0 {
		dots := " " + strings.Repeat(".", dots)
		sb.WriteString(ansi.FGGray.Format(dots))
	}
	fmt.Fprintf(&sb, " %v\n", suffix)
	return sb.String()
}

func (l *TTYListener) progress(name string, from, to int64, err error) string {
	state := l.spinner.Update()
	if from == 0 {
		state = ansi.FGYellow.Format("-")
	} else if from == to {
		state = ansi.FGGreen.Format("✓")
	}
	if err != nil {
		state = ansi.FGRed.Format("x")
	}

	percent := 0
	if to > 0 {
		percent = int((float64(from) / float64(to)) * 100.00)
	}
	fromStr, toStr := "0", "?"
	if from > 0 {
		fromStr = data.Quantity(from).String()
	}
	if to > 0 {
		toStr = data.Quantity(to).String()
	}

	suffix := fmt.Sprintf("(%v / %v)", fromStr, toStr)
	value := fmt.Sprintf("%3d%% %-28s", percent, suffix)
	return l.valueProgress(state, name, value)
}
