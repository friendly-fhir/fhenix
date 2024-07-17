package cmd

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"atomicgo.dev/cursor"
	"github.com/friendly-fhir/fhenix/driver"
	"github.com/friendly-fhir/fhenix/internal/ansi"
	"github.com/friendly-fhir/fhenix/internal/snek"
)

var spinner = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}

type TTYListener struct {
	offset    int
	downloads map[string]*download
	m         sync.Mutex

	verbose bool
	out     io.Writer
	driver.BaseListener
}

func NewProgressListener(ctx context.Context, verbose bool) *TTYListener {
	out := snek.CommandOut(ctx)
	return &TTYListener{
		out:     out,
		verbose: verbose,
	}
}

func keyOf(registry, pkg, version string) string {
	return fmt.Sprintf("%v::%s@%s", registry, pkg, version)
}

func (l *TTYListener) BeforeFetch(registry, pkg, version string) {
	key := keyOf(registry, pkg, version)
	l.m.Lock()
	defer l.m.Unlock()
	offset := l.offset
	l.offset++
	if l.downloads == nil {
		l.downloads = make(map[string]*download)
	}
	l.downloads[key] = &download{
		Offset: offset,
	}
	name := fmt.Sprintf("%s@%s", pkg, version)

	cursor.DownAndClear(offset)
	fmt.Fprint(l.out, progress(name, 0, 0, 0, nil))
	cursor.Up(offset + 1)
}

func (l *TTYListener) OnFetch(registry, pkg, version string, data int64) {
	key := keyOf(registry, pkg, version)
	l.m.Lock()
	defer l.m.Unlock()
	download := l.downloads[key]
	download.TotalBytes = data

	name := fmt.Sprintf("%s@%s", pkg, version)

	cursor.DownAndClear(download.Offset)
	fmt.Fprint(l.out, progress(name, 0, data, download.State, nil))
	cursor.Up(download.Offset + 1)
}

func (l *TTYListener) OnFetchWrite(registry, pkg, version string, data []byte) {
	key := keyOf(registry, pkg, version)
	l.m.Lock()
	defer l.m.Unlock()
	download := l.downloads[key]
	download.Current += int64(len(data))
	download.State++

	name := fmt.Sprintf("%s@%s", pkg, version)

	cursor.DownAndClear(download.Offset)
	fmt.Fprint(l.out, progress(name, download.Current, download.TotalBytes, download.State, nil))
	cursor.Up(download.Offset + 1)
}

func (l *TTYListener) AfterFetch(registry, pkg, version string, err error) {
	key := keyOf(registry, pkg, version)
	l.m.Lock()
	defer l.m.Unlock()
	download := l.downloads[key]

	name := fmt.Sprintf("%s@%s", pkg, version)

	cursor.DownAndClear(download.Offset)
	fmt.Fprint(l.out, progress(name, download.Current, download.TotalBytes, download.State, nil))
	cursor.Up(download.Offset + 1)
}

func (l *TTYListener) OnCacheHit(registry, pkg, version string) {
	key := keyOf(registry, pkg, version)
	l.m.Lock()
	defer l.m.Unlock()
	if l.downloads == nil {
		l.downloads = make(map[string]*download)
	}
	dl, ok := l.downloads[key]
	if !ok {
		offset := l.offset
		l.offset++
		dl = &download{
			Offset: offset,
		}
		l.downloads[key] = dl
	}
	if dl.Current != 0 {
		return
	}
	name := fmt.Sprintf("%s@%s", pkg, version)

	cursor.DownAndClear(dl.Offset)
	fmt.Fprint(l.out, cacheProgress(name))
	cursor.Up(dl.Offset + 1)
}

var _ driver.Listener = (*TTYListener)(nil)

type download struct {
	TotalBytes int64
	Current    int64
	Offset     int
	State      int
}

func cacheProgress(name string) string {
	const (
		maxLength = 50
	)
	var (
		sb strings.Builder
	)
	state := ansi.FGYellow.Format("✓")
	fmt.Fprintf(&sb, "[%s] %s", state, name)
	beginning := ansi.StripFormat(sb.String())
	if len(beginning) < maxLength {
		sb.WriteString(strings.Repeat(".", maxLength-len(beginning)))
	}
	fmt.Fprintf(&sb, " cached\n")
	return sb.String()
}

func progress(name string, from, to int64, seq int, err error) string {
	const (
		maxLength = 50
	)
	var (
		sb strings.Builder
	)
	state := spinner[seq%len(spinner)]
	if from == 0 {
		state = ansi.FGYellow.Format("-")
	} else if from == to {
		state = ansi.FGGreen.Format("✓")
	}
	if err != nil {
		state = ansi.FGRed.Format("x")
	}
	fmt.Fprintf(&sb, "[%s] %s", state, name)
	beginning := ansi.StripFormat(sb.String())
	if err != nil {
		fmt.Fprintf(&sb, " error: %v", err)
		return sb.String()
	}

	if len(beginning) < maxLength {
		sb.WriteString(strings.Repeat(".", maxLength-len(beginning)))
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
	if len(toStr) > len(fromStr) {
		fromStr = strings.Repeat(" ", len(toStr)-len(fromStr)) + fromStr
	}
	fmt.Fprintf(&sb, " % 4d%% (%v / %v)\n", percent, fromStr, toStr)
	return sb.String()
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
