package registry

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/friendly-fhir/fhenix/registry/internal/archive"
)

const (
	Local   = "local"
	Default = "default"
)

// Cache is a cache of downloaded and used FHIR packages from a given registry.
type Cache struct {
	outputPath string

	clients   map[string]*Client
	listeners Listeners

	// localPackages is a map of local packages that are not fetched from a
	// remote registry, but are explicitly added to the cache.
	localPackages map[string]string
}

// NewCache creates a new cache with the specified output path.
func NewCache(outputPath string) *Cache {
	return &Cache{
		outputPath: filepath.FromSlash(outputPath),
		clients: map[string]*Client{
			Default: DefaultClient,
		},
		localPackages: make(map[string]string),
	}
}

// DefaultCache creates a new cache with the default output path. The behavior
// can be controlled with the presence of the FHIR_CACHE environment variable.
func DefaultCache() *Cache {
	path, ok := os.LookupEnv("FHIR_CACHE")
	if ok {
		return NewCache(filepath.FromSlash(path))
	}
	var err error
	path, err = os.UserHomeDir()
	if err == nil {
		return NewCache(filepath.Join(path, ".fhir"))
	}
	return NewCache(".fhir")
}

// AddListener adds a listener to the cache.
func (c *Cache) AddListener(listener CacheListener) {
	c.listeners = append(c.listeners, listener)
}

// AddClient adds a client to the cache.
func (c *Cache) AddClient(registry string, client *Client) {
	c.clients[registry] = client
}

// AddLocalPackage adds a local package to the cache that may be referenced
// later. The package registry is always considered "local".
func (c *Cache) AddLocalPackage(pkg, version, path string) {
	c.localPackages[c.localKey(pkg, version)] = path
}

// Contains returns true if the cache contains the specified package.
// Containment does not imply that the package is valid or usable -- just that
// the contents can be found on disk.
func (c *Cache) Contains(registry, pkg, version string) bool {
	pkgFile := filepath.Join(c.CacheDir(registry, pkg, version), "package.json")
	if _, err := os.Stat(pkgFile); err != nil {
		return false
	}
	return true
}

// Root returns the root directory of the cache.
func (c *Cache) Root() string {
	return c.outputPath
}

// Delete removes the specified package from the cache.
func (c *Cache) Delete(registry, pkg, version string) error {
	if registry == "" || pkg == "" || version == "" {
		return fmt.Errorf("fhir cache: registry, package, and version must be specified")
	}
	if registry == Local {
		delete(c.localPackages, c.localKey(pkg, version))
		return nil
	}

	cache := c.CacheDir(registry, pkg, version)
	if cache == "" {
		return fmt.Errorf("fhir cache: unknown registry %q", registry)
	}
	c.listeners.OnDelete(registry, pkg, version)
	return os.RemoveAll(cache)
}

// Fetch downloads the specified package from the registry.
func (c *Cache) Fetch(ctx context.Context, registry, pkg, version string) error {
	if c.Contains(registry, pkg, version) {
		c.listeners.OnCacheHit(registry, pkg, version)
		return nil
	}
	return c.ForceFetch(ctx, registry, pkg, version)
}

// ForceFetch forces a download of the specified package from the registry.
func (c *Cache) ForceFetch(ctx context.Context, registry, pkg, version string) error {
	if registry == Local {
		c.listeners.OnCacheHit(registry, pkg, version)
		return nil
	}

	client, ok := c.clients[registry]
	if !ok {
		return fmt.Errorf("fhir cache: unknown name %q", registry)
	}
	c.listeners.BeforeFetch(registry, pkg, version)
	content, size, err := client.Fetch(ctx, pkg, version)
	if err != nil {
		return err
	}
	defer content.Close()

	r := io.TeeReader(content, writerFunc(func(p []byte) {
		c.listeners.OnFetchWrite(registry, pkg, version, p)
	}))
	c.listeners.OnFetch(registry, pkg, version, size)
	defer c.listeners.AfterFetch(registry, pkg, version, err)

	unpackers := archive.Unpackers{
		archive.UnpackFunc(func(s string, i int64, _ io.Reader) error {
			c.listeners.OnUnpack(registry, pkg, version, s, i)
			return nil
		}),
		&archive.DiskUnpacker{
			Root: c.CacheDir(registry, pkg, version),
			Tee: func(name string, r io.Reader) io.Reader {
				return io.TeeReader(r, writerFunc(func(bytes []byte) {
					c.listeners.OnUnpackWrite(registry, pkg, version, name, bytes)
				}))
			},
		},
	}
	var opts []archive.Option

	// Only unpack JSON files that are not the .index.json
	opts = append(opts, archive.Filter(func(s string) bool {
		base := filepath.Base(s)
		ext := filepath.Ext(base)
		return ext == ".json" && (base != ".index.json")
	}))

	// Flatten all output
	opts = append(opts, archive.Transform(func(s string) string {
		return filepath.Base(s)
	}))

	tar := archive.New(r, opts...)
	return tar.Unpack(unpackers)
}

// CacheDir returns the directory where the specified package is cached.
// If the registry is unknown, or if any parameters are not set, an empty string
// is returned.
func (c *Cache) CacheDir(registry, pkg, version string) string {
	if registry == Local {
		return c.localPackages[c.localKey(pkg, version)]
	}
	client, ok := c.clients[registry]
	if !ok || pkg == "" || version == "" || registry == "" {
		return ""
	}
	return filepath.Join(c.outputPath, stripURL(client.url), pkg, version)
}

// Get returns the package from the cache.
func (c *Cache) Get(registry, pkg, version string) (*Package, error) {
	if registry == "" || pkg == "" || version == "" {
		return nil, fmt.Errorf("fhir cache: registry, package and version must be specified")
	}

	path := c.CacheDir(registry, pkg, version)
	if path == "" {
		return nil, fmt.Errorf("fhir cache: unknown registry %q", registry)
	}

	return NewPackage(path)
}

// GetOrFetch returns the package from the cache, or fetches it if it is not
// present.
func (c *Cache) GetOrFetch(ctx context.Context, registry, pkg, version string) (*Package, error) {
	if !c.Contains(registry, pkg, version) {
		if err := c.Fetch(ctx, registry, pkg, version); err != nil {
			return nil, err
		}
	}

	return c.Get(registry, pkg, version)
}

// localKey returns a key for a local package.
func (c *Cache) localKey(pkg, version string) string {
	return fmt.Sprintf("%s@%s", pkg, version)
}

type writerFunc func([]byte)

func (w writerFunc) Write(p []byte) (n int, err error) {
	w(p)
	return len(p), nil
}

// StripURL strips the scheme, ports, and parameters from a URL and returns a valid filepath.
func stripURL(u string) string {
	// Remove scheme
	parts := strings.Split(u, "://")
	if len(parts) == 2 {
		u = parts[1]
	}

	// Remove parameters
	u = strings.Split(u, "?")[0]

	// Remove port, if present
	parts = strings.Split(u, "/")
	if len(parts) > 1 {
		parts[0] = strings.Split(parts[0], ":")[0]
	}
	return filepath.Join(parts...)
}
