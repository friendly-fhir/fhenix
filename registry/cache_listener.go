package registry

// CacheListener is a listener for events occurring when populating the cache.
type CacheListener interface {

	// OnFetch is an event handler invoked when a fetch operation is initiated.
	OnFetch(registry, pkg, version string, size int64)

	// OnFetchWrite is an event handler invoked when bytes written during a fetch
	// operation.
	OnFetchWrite(registry, pkg, version string, bytes []byte)

	// OnUnpack is an event handler invoked when a file is unpacked.
	OnUnpack(registry, pkg, version, file string, size int64)

	// OnUnpackWrite is an event handler invoked when bytes written during a
	// file unpack operation.
	OnUnpackWrite(registry, pkg, version, file string, bytes []byte)

	// OnDelete is an event handler invoked when a package is deleted from the
	// cache.
	OnDelete(registry, pkg, version string)

	// OnCacheHit is an event handler invoked when a cache hit occurs.
	OnCacheHit(registry, pkg, version string)

	listener()
}

// BaseCacheListener is a base implementation of the [CacheListener] interface.
type BaseCacheListener struct{}

// OnFetch is a no-op implementation of the [CacheListener] interface.
func (BaseCacheListener) OnFetch(registry, pkg, version string, size int64) {}

// OnFetchWrite is a no-op implementation of the [CacheListener] interface.
func (BaseCacheListener) OnFetchWrite(registry, pkg, version string, bytes []byte) {}

// OnUnpack is a no-op implementation of the [CacheListener] interface.
func (BaseCacheListener) OnUnpack(registry, pkg, version, file string, size int64) {}

// OnUnpackWrite is a no-op implementation of the [CacheListener] interface.
func (BaseCacheListener) OnUnpackWrite(registry, pkg, version, file string, bytes []byte) {}

// OnCacheHit is a no-op implementation of the [CacheListener] interface.
func (BaseCacheListener) OnCacheHit(registry, pkg, version string) {}

// OnDelete is a no-op implementation of the [CacheListener] interface.
func (BaseCacheListener) OnDelete(registry, pkg, version string) {}

func (BaseCacheListener) listener() {}

var _ CacheListener = (*BaseCacheListener)(nil)

// Listeners is a collection of [CacheListener]s that itself can be treated as
// a single listener.
type Listeners []CacheListener

// OnFetch is an event handler invoked when a fetch operation is initiated.
func (l Listeners) OnFetch(registry, pkg, version string, size int64) {
	for _, listener := range l {
		listener.OnFetch(registry, pkg, version, size)
	}
}

// OnFetchWrite is an event handler invoked when bytes written during a fetch
// operation.
func (l Listeners) OnFetchWrite(registry, pkg, version string, bytes []byte) {
	for _, listener := range l {
		listener.OnFetchWrite(registry, pkg, version, bytes)
	}
}

// OnUnpack is an event handler invoked when a file is unpacked.
func (l Listeners) OnUnpack(registry, pkg, version, file string, size int64) {
	for _, listener := range l {
		listener.OnUnpack(registry, pkg, version, file, size)
	}
}

// OnUnpackWrite is an event handler invoked when bytes written during a
// file unpack operation.
func (l Listeners) OnUnpackWrite(registry, pkg, version, file string, bytes []byte) {
	for _, listener := range l {
		listener.OnUnpackWrite(registry, pkg, version, file, bytes)
	}
}

// OnDelete is an event handler invoked when a package is deleted from the
// cache.
func (l Listeners) OnDelete(registry, pkg, version string) {
	for _, listener := range l {
		listener.OnDelete(registry, pkg, version)
	}
}

// OnCacheHit is an event handler invoked when a cache hit occurs.
func (l Listeners) OnCacheHit(registry, pkg, version string) {
	for _, listener := range l {
		listener.OnCacheHit(registry, pkg, version)
	}
}

func (l Listeners) listener() {}

var _ CacheListener = (*Listeners)(nil)
