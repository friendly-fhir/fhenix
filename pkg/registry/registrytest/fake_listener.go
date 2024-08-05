package registrytest

import (
	"bytes"
	"fmt"
	"io"

	"github.com/friendly-fhir/fhenix/pkg/registry"
)

// CacheListener is a fake implementation of the [registry.CacheListener] interface.
type CacheListener struct {
	onFetch       map[string]*fetch
	onFetchWrite  map[string]*write
	onUnpack      map[string]*fetch
	onUnpackWrite map[string]*write
	onCacheHit    map[string]int
	onDelete      map[string]int

	registry.BaseCacheListener
}

func NewCacheListener() *CacheListener {
	return &CacheListener{
		onFetch:       make(map[string]*fetch),
		onFetchWrite:  make(map[string]*write),
		onUnpack:      make(map[string]*fetch),
		onUnpackWrite: make(map[string]*write),
		onCacheHit:    make(map[string]int),
		onDelete:      make(map[string]int),
	}
}

// FetchCalls returns the number of times [OnFetch] was called for the given
// package.
func (l *CacheListener) FetchCalls(registry, pkg, version string) int {
	key := l.prepare(registry, pkg, version)
	return l.onFetch[key].calls
}

// FetchBytes returns the total number of bytes fetched for the given package.
func (l *CacheListener) FetchBytes(registry, pkg, version string) int64 {
	key := l.prepare(registry, pkg, version)
	return l.onFetch[key].bytes
}

// FetchWriteCalls returns the number of times [OnFetchWrite] was called for the
// given package.
func (l *CacheListener) FetchWriteCalls(registry, pkg, version string) int {
	key := l.prepare(registry, pkg, version)
	return l.onFetchWrite[key].calls
}

// FetchWriteTotal returns the total number of bytes written during fetch for the
// given package.
func (l *CacheListener) FetchWriteTotal(registry, pkg, version string) int64 {
	key := l.prepare(registry, pkg, version)
	return l.onFetchWrite[key].total
}

// FetchWriteContent returns a reader that reads the bytes written during fetch
// for the given package.
func (l *CacheListener) FetchWriteContent(registry, pkg, version string) io.Reader {
	key := l.prepare(registry, pkg, version)
	fetch := l.onFetchWrite[key]
	var readers []io.Reader
	for _, chunk := range fetch.chunks {
		readers = append(readers, bytes.NewReader(chunk))
	}
	return io.MultiReader(readers...)
}

// UnpackCalls returns the number of times [OnUnpack] was called for the given
// package.
func (l *CacheListener) UnpackCalls(registry, pkg, version string) int {
	key := l.prepare(registry, pkg, version)
	return l.onUnpack[key].calls
}

// UnpackTotal returns the total number of bytes unpacked for the given package.
func (l *CacheListener) UnpackTotal(registry, pkg, version string) int64 {
	key := l.prepare(registry, pkg, version)
	return l.onUnpack[key].bytes
}

// UnpackFileCalls returns the number of times [OnUnpack] was called for the given
// package and file.
func (l *CacheListener) UnpackFileCalls(registry, pkg, version, file string) int {
	key := l.prepareUnpack(registry, pkg, version, file)
	return l.onUnpack[key].calls
}

// UnpackFileTotal returns the total number of bytes unpacked for the given file
// in the package.
func (l *CacheListener) UnpackFileTotal(registry, pkg, version, file string) int64 {
	key := l.prepareUnpack(registry, pkg, version, file)
	return l.onUnpack[key].bytes
}

// UnpackWriteCalls returns the number of times [OnUnpackWrite] was called for the
// given package.
func (l *CacheListener) UnpackWriteCalls(registry, pkg, version string) int {
	key := l.prepare(registry, pkg, version)
	return l.onUnpackWrite[key].calls
}

// UnpackWriteTotal returns the total number of bytes written during unpack for
// the given package.
func (l *CacheListener) UnpackWriteTotal(registry, pkg, version string) int64 {
	key := l.prepare(registry, pkg, version)
	return l.onUnpackWrite[key].total
}

// UnpackWriteContent returns a reader that reads the bytes written during unpack
// for the given package.
func (l *CacheListener) UnpackWriteContent(registry, pkg, version string) io.Reader {
	key := l.prepare(registry, pkg, version)
	unpack := l.onUnpackWrite[key]
	var readers []io.Reader
	for _, chunk := range unpack.chunks {
		readers = append(readers, bytes.NewReader(chunk))
	}
	return io.MultiReader(readers...)
}

// DeleteCalls returns the number of times [OnDelete] was called for the given
// package.
func (l *CacheListener) DeleteCalls(registry, pkg, version string) int {
	key := l.prepare(registry, pkg, version)
	return l.onDelete[key]
}

// CacheHits returns the number of times [OnCacheHit] was called for the given
// package.
func (l *CacheListener) CacheHits(registry, pkg, version string) int {
	key := l.prepare(registry, pkg, version)
	return l.onCacheHit[key]
}

// OnFetch is a no-op implementation of the [CacheListener] interface.
func (l *CacheListener) OnFetch(registry, pkg, version string, bytes int64) {
	key := l.prepare(registry, pkg, version)
	l.onFetch[key].calls++
	l.onFetch[key].bytes += bytes
}

// OnFetchWrite is a no-op implementation of the [CacheListener] interface.
func (l *CacheListener) OnFetchWrite(registry, pkg, version string, bytes []byte) {
	key := l.prepare(registry, pkg, version)
	l.onFetchWrite[key].calls++
	l.onFetchWrite[key].total += int64(len(bytes))
	l.onFetchWrite[key].chunks = append(l.onFetchWrite[key].chunks, bytes)
}

// OnUnpack is a no-op implementation of the [CacheListener] interface.
func (l *CacheListener) OnUnpack(registry, pkg, version, file string, size int64) {
	key := l.prepareUnpack(registry, pkg, version, file)
	l.onUnpack[key].calls++
	l.onUnpack[key].bytes += size
	key = l.prepare(registry, pkg, version)
	l.onUnpack[key].calls++
	l.onUnpack[key].bytes += size
}

// OnUnpackWrite is a no-op implementation of the [CacheListener] interface.
func (l *CacheListener) OnUnpackWrite(registry, pkg, version, file string, bytes []byte) {
	key := l.prepareUnpack(registry, pkg, version, file)
	l.onUnpackWrite[key].calls++
	l.onUnpackWrite[key].total += int64(len(bytes))
	l.onUnpackWrite[key].chunks = append(l.onUnpackWrite[key].chunks, bytes)
	key = l.prepare(registry, pkg, version)
	l.onUnpackWrite[key].calls++
	l.onUnpackWrite[key].total += int64(len(bytes))
	l.onUnpackWrite[key].chunks = append(l.onUnpackWrite[key].chunks, bytes)
}

// OnCacheHit is a no-op implementation of the [CacheListener] interface.
func (l *CacheListener) OnCacheHit(registry, pkg, version string) {
	key := l.prepare(registry, pkg, version)
	l.onCacheHit[key]++
}

// OnDelete is a no-op implementation of the [CacheListener] interface.
func (l *CacheListener) OnDelete(registry, pkg, version string) {
	key := l.prepare(registry, pkg, version)
	l.onDelete[key]++
}

func (l *CacheListener) key(registry, pkg, version string) string {
	return registry + "/" + pkg + "@" + version
}

func (l *CacheListener) prepare(registry, pkg, version string) string {
	key := l.key(registry, pkg, version)
	if _, ok := l.onFetch[key]; !ok {
		l.onFetch[key] = &fetch{}
	}
	if _, ok := l.onFetchWrite[key]; !ok {
		l.onFetchWrite[key] = &write{}
	}
	if _, ok := l.onDelete[key]; !ok {
		l.onDelete[key] = 0
	}
	if _, ok := l.onCacheHit[key]; !ok {
		l.onCacheHit[key] = 0
	}
	return key
}

func (l *CacheListener) prepareUnpack(registry, pkg, version, file string) string {
	key := l.key(registry, pkg, version)
	fileKey := fmt.Sprintf("%s/%s", key, file)
	if _, ok := l.onUnpack[key]; !ok {
		l.onUnpack[key] = &fetch{}
	}
	if _, ok := l.onUnpackWrite[key]; !ok {
		l.onUnpackWrite[key] = &write{}
	}
	if _, ok := l.onUnpack[fileKey]; !ok {
		l.onUnpack[fileKey] = &fetch{}
	}
	if _, ok := l.onUnpackWrite[fileKey]; !ok {
		l.onUnpackWrite[fileKey] = &write{}
	}
	return fileKey
}

var _ registry.CacheListener = (*CacheListener)(nil)

type fetch struct {
	calls int
	bytes int64
}

type write struct {
	calls  int
	total  int64
	chunks [][]byte
}
