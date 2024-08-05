package registry

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"golang.org/x/sync/errgroup"
)

// Downloader is a utility for downloading packages and their dependencies
// into the [Cache].
type Downloader struct {
	cache *Cache

	workers int
	force   bool

	requests []*request

	m sync.Mutex
}

// NewDownloader creates a new [Downloader].
func NewDownloader(dest *Cache) *Downloader {
	return &Downloader{
		cache:   dest,
		workers: runtime.NumCPU(),
		force:   false,
	}
}

// Workers sets the number of worker threads. If it is less than 1, the number
// of worker threads will be set to the number of CPUs on the system.
func (d *Downloader) Workers(n int) *Downloader {
	d.m.Lock()
	defer d.m.Unlock()
	d.workers = n
	if d.workers < 1 {
		d.workers = runtime.NumCPU()
	}
	return d
}

// Force forces the downloader to fetch the package from the registry, even if
// it is already present in the cache.
func (d *Downloader) Force(b bool) *Downloader {
	d.m.Lock()
	defer d.m.Unlock()
	d.force = b
	return d
}

// Add an explicit package to download from the registry. If includeDependencies
// is set to true, the downloader will also download the dependencies for the
// package.
func (d *Downloader) Add(registry, pkg, version string, includeDependencies bool) {
	d.m.Lock()
	defer d.m.Unlock()
	d.requests = append(d.requests, &request{
		registry:            registry,
		pkg:                 pkg,
		version:             version,
		includeDependencies: includeDependencies,
	})
}

type queue struct {
	group    *errgroup.Group
	wg       *sync.WaitGroup
	requests chan *request
}

func (q *queue) Add(req *request) {
	q.wg.Add(1)
	q.group.Go(func() error {
		q.requests <- req
		return nil
	})
}

func (q *queue) More() <-chan *request {
	return q.requests
}

func (q *queue) Close() {
	close(q.requests)
}

// Start the download process which will load the cache with the downloaded
// packages and their dependencies.
//
// This function will block until all worker threads have completed downloading
// the packages and their dependencies. If an error occurs during the download
// process, the function will return the first encountered error.
func (d *Downloader) Start(ctx context.Context) error {
	// Copy the state of the downloader as it is seen by the time we start.
	d.m.Lock()
	initialRequests := make([]*request, len(d.requests))
	copy(initialRequests, d.requests)
	force := d.force
	workers := d.workers
	d.m.Unlock()

	requests := make(chan *request)
	wg := &sync.WaitGroup{}

	visited := &sync.Map{}
	group, ctx := errgroup.WithContext(ctx)
	queue := &queue{
		group:    group,
		requests: requests,
		wg:       wg,
	}

	// pre-load the queue with the initial requests
	for _, req := range initialRequests {
		queue.Add(req)
	}

	// start the worker threads
	for i := 0; i < workers; i++ {
		group.Go(d.worker(ctx, wg, queue, force, visited))
	}
	// Wait for all jobs to be completed.
	go func() {
		wg.Wait()
		queue.Close()
	}()

	return group.Wait()
}

func (d *Downloader) worker(ctx context.Context, wg *sync.WaitGroup, queue *queue, force bool, visited *sync.Map) func() error {
	return func() error {
		for req := range queue.More() {
			pkg, err := d.download(ctx, req, force, visited)
			if err != nil {
				wg.Done()
				return err
			}
			// a nil package means it was already visited it before in a different
			// request.
			if pkg == nil {
				wg.Done()
				continue
			}

			if req.includeDependencies {
				for name, version := range pkg.Dependencies() {
					queue.Add(&request{
						registry:            req.registry,
						pkg:                 name,
						version:             version,
						includeDependencies: req.includeDependencies,
					})
				}
			}
			wg.Done()
		}
		return nil
	}
}

func (d *Downloader) download(ctx context.Context, req *request, force bool, visited *sync.Map) (*Package, error) {
	if _, ok := visited.LoadOrStore(req.key(), struct{}{}); ok {
		return nil, nil
	}

	var pkg *Package
	var err error
	if force {
		if err := d.cache.ForceFetch(ctx, req.registry, req.pkg, req.version); err != nil {
			return nil, err
		}
		pkg, err = d.cache.Get(req.registry, req.pkg, req.version)
		if err != nil {
			return nil, err
		}
	} else {
		pkg, err = d.cache.GetOrFetch(ctx, req.registry, req.pkg, req.version)
		if err != nil {
			return nil, err
		}
	}
	return pkg, nil
}

type request struct {
	registry string
	pkg      string
	version  string

	includeDependencies bool
}

func (r *request) key() string {
	return fmt.Sprintf("%s::%s@%s", r.registry, r.pkg, r.version)
}
