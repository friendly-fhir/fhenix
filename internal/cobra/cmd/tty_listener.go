package cmd

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/friendly-fhir/fhenix/driver"
	"github.com/friendly-fhir/fhenix/internal/ansi"
	"github.com/friendly-fhir/fhenix/internal/snek/spinner"
	"github.com/friendly-fhir/fhenix/internal/snek/terminal"
)

type TTYListener struct {
	offset    int
	downloads map[string]*download
	m         sync.Mutex

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

var _ driver.Listener = (*TTYListener)(nil)

type download struct {
	TotalBytes int64
	Current    int64
	Line       *terminal.Line
	Spinner    *spinner.Spinner
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
		fromStr = toDataUnit(from)
	}
	if to > 0 {
		toStr = toDataUnit(to)
	}

	suffix := fmt.Sprintf("(%v / %v)", fromStr, toStr)
	value := fmt.Sprintf("%3d%% %-28s", percent, suffix)
	return l.valueProgress(state, name, value)
}

func toDataUnit(units int64) string {
	if units < 1024 {
		return fmt.Sprintf("%d B", units)
	}
	if units < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(units)/1024)
	}
	if units < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(units)/1024/1024)
	}
	if units < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2f GB", float64(units)/1024/1024/1024)
	}
	// Should never reach here?
	return fmt.Sprintf("%.2f TB", float64(units)/1024/1024/1024/1024)
}
