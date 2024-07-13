/*
Package fhirig enables a way to reference and cache entries in a FHIR
Implementation Guide that has been published to simplifier.net.
*/
package fhirig

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

// Fetcher is a mechansim for fetching content from a URL.
type Fetcher interface {
	// Fetch fetches the content from the URL, returning a reader.
	Fetch(url string) (io.ReadCloser, error)
}

// Package represents a versioned and named FHIR package from a registry.
type Package struct {
	name    string
	version string
}

// String returns a string representing the package version.
func (p *Package) String() string {
	if p == nil {
		return ""
	}
	return fmt.Sprintf("%s@%s", p.name, p.version)
}

// Equal returns true if the package is equal to the other package.
func (p *Package) Equal(other *Package) bool {
	return p.String() == other.String()
}

// Name returns the name of the package.
func (p *Package) Name() string {
	return p.name
}

// Version returns the version of the package.
func (p *Package) Version() string {
	return p.version
}

// NewPackage constructs a Package with the specified name and version.
func NewPackage(name, version string) *Package {
	return &Package{name: name, version: version}
}

// ParsePackage parses a package string into a Package.
func ParsePackage(s string) (*Package, error) {
	parts := strings.Split(s, "@")
	if len(parts) != 2 {
		return nil, errors.New("invalid package string")
	}
	return NewPackage(parts[0], parts[1]), nil
}

type Listener interface {
	OnFetchStart(pkg *Package)
	OnFetchEnd(pkg *Package, err error)
	OnCacheHit(pkg *Package)
}

type BaseListener struct{}

func (*BaseListener) OnFetchStart(pkg *Package)          {}
func (*BaseListener) OnFetchEnd(pkg *Package, err error) {}
func (*BaseListener) OnCacheHit(pkg *Package)            {}

type FetcherFunc func(url string) (io.ReadCloser, error)

func (f FetcherFunc) Fetch(url string) (io.ReadCloser, error) {
	return f(url)
}

// PackageCache is a cache for FHIR packages.
type PackageCache struct {
	// Registry is the URL of the FHIR registry (e.g. https://simplifier.net)
	Registry string

	// Root is the root directory of the cache.
	Root string

	// Fetcher is the fetching mechanism.
	Fetcher Fetcher

	// Listener is the listener for fetch events.
	Listener Listener
}

// Clear removes the package from the cache.
func (r *PackageCache) Clear(pkg *Package) error {
	if !r.Has(pkg) {
		return nil
	}
	return os.RemoveAll(r.path(pkg))
}

func (r *PackageCache) listener() Listener {
	if r.Listener != nil {
		return r.Listener
	}
	return &BaseListener{}
}

// Get returns the directory entries for the package and version.
func (r *PackageCache) Get(pkg *Package) ([]string, error) {
	root := r.path(pkg)
	dir, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var entries []string
	for _, entry := range dir {
		entries = append(entries, filepath.Join(root, entry.Name()))
	}
	return entries, nil
}

// Path returns the path to the specified package.
func (r *PackageCache) Path(pkg *Package) (string, error) {
	if r.Has(pkg) {
		return r.path(pkg), nil
	}
	return "", os.ErrNotExist
}

func (r *PackageCache) path(pkg *Package, parts ...string) string {
	args := make([]string, 0, 3+len(parts))
	args = append(args, r.Root, pkg.Name(), pkg.Version())
	args = append(args, parts...)
	return filepath.Join(args...)
}

// FetchAndGet fetches the versioned package if it is not already in the cache,
// and returns the directory entries for the package.
func (r *PackageCache) FetchAndGet(ctx context.Context, pkg *Package) ([]string, error) {
	if err := r.Fetch(ctx, pkg); err != nil {
		return nil, err
	}
	return r.Get(pkg)
}

// ForceFetch forces a fetch of the package and version.
func (r *PackageCache) ForceFetch(ctx context.Context, pkg *Package) error {
	return r.fetch(ctx, pkg, true)
}

// Fetch fetches the package and version.
func (r *PackageCache) Fetch(ctx context.Context, pkg *Package) error {
	return r.fetch(ctx, pkg, false)
}

func (r *PackageCache) Dependencies(pkg *Package) ([]*Package, error) {
	path := filepath.Join(r.Root, pkg.Name(), pkg.Version(), "package.json")
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	packageJSON := struct {
		Dependencies map[string]string `json:"dependencies"`
	}{}
	err = json.NewDecoder(file).Decode(&packageJSON)
	if err != nil {
		return nil, err
	}

	var deps []*Package
	for dep, ver := range packageJSON.Dependencies {
		deps = append(deps, NewPackage(dep, ver))
	}
	return deps, nil
}

// Has returns true if the package and version is in the cache.
func (r *PackageCache) Has(pkg *Package) bool {
	path := filepath.Join(r.Root, pkg.Name(), pkg.Version(), "package.tar.gz")
	_, err := os.Stat(path)
	return err == nil
}

func (r *PackageCache) fetcher() Fetcher {
	if r.Fetcher == nil {
		return defaultFetcher(0)
	}
	return r.Fetcher
}

func (r *PackageCache) fetch(ctx context.Context, pkg *Package, force bool) error {
	listener := r.listener()
	if !force && r.Has(pkg) {
		listener.OnCacheHit(pkg)
		return nil
	}

	url := fmt.Sprintf("%s/%s/%s", r.registry(), pkg.Name(), pkg.Version())
	listener.OnFetchStart(pkg)
	pkgFile, err := r.fetcher().Fetch(url)
	if err != nil {
		listener.OnFetchEnd(pkg, err)
		return err
	}
	defer pkgFile.Close()
	root := r.path(pkg)
	if err := os.MkdirAll(root, 0755); err != nil {
		listener.OnFetchEnd(pkg, err)
		return err
	}
	path := filepath.Join(root, "package.tar.gz")
	file, err := os.Create(path)
	if err != nil {
		listener.OnFetchEnd(pkg, err)
		return err
	}
	defer file.Close()
	success := false
	// If fetching didn't succeed, remove anything that was downloaded.
	defer func() {
		if !success {
			_ = os.RemoveAll(path)
		}
	}()

	reader := io.TeeReader(pkgFile, file)

	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		listener.OnFetchEnd(pkg, err)
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			listener.OnFetchEnd(pkg, err)
			return err
		}

		if header.Typeflag != tar.TypeReg {
			continue
		}

		name, ok := strings.CutPrefix(header.Name, "package/")
		if !ok {
			continue
		}

		// Skip files in subdirectories
		if strings.Contains(name, "/") {
			continue
		}

		// Skip .index.json
		if name == ".index.json" {
			continue
		}

		// Skip non StructureDefinition, CodeSystem, ValueSet, or ConceptMap files
		if name != "package.json" && !hasPrefix(name, "StructureDefinition-", "CodeSystem-", "ValueSet-", "ConceptMap-") {
			continue
		}

		// Create the destination file
		path := r.path(pkg, name)
		file, err := os.Create(path)
		if err != nil {
			listener.OnFetchEnd(pkg, err)
			return err
		}
		pkgReader := io.TeeReader(tarReader, file)

		// If the file is a package.json file, download its dependencies, and still
		// copy the file contents
		if name == "package.json" {
			err = r.getDependencies(ctx, pkgReader, force)
			if err != nil {
				return err
			}
		} else if _, err := io.Copy(io.Discard, pkgReader); err != nil {
			return err
		}
	}
	success = true
	listener.OnFetchEnd(pkg, nil)
	return nil
}

func (r *PackageCache) getDependencies(ctx context.Context, reader io.Reader, force bool) error {
	packageJSON := struct {
		Dependencies map[string]string `json:"dependencies"`
	}{}
	err := json.NewDecoder(reader).Decode(&packageJSON)
	if err != nil {
		return err
	}
	group, ctx := errgroup.WithContext(ctx)
	for dep, ver := range packageJSON.Dependencies {
		dep, ver := dep, ver
		group.Go(func() error {
			return r.fetch(ctx, NewPackage(dep, ver), force)
		})
	}
	if err := group.Wait(); err != nil {
		return err
	}
	return nil
}

func (r *PackageCache) registry() string {
	if r.Registry == "" {
		return "https://packages.simplifier.net"
	}
	return r.Registry
}

func hasPrefix(name string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

type defaultFetcher int

func (defaultFetcher) Fetch(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
