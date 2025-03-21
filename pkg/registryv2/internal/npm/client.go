package npm

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"net/http"
)

// HTTPClient is an interface that models the [http.Client] struct type.
//
// This enables dependency-injecting the client around for different purposes,
// which is leveraged internally in this package so that the npmtest
// package can provide non-networked clients for testing.
type HTTPClient interface {
	// Do sends an HTTP request and returns an HTTP response, following
	Do(req *http.Request) (*http.Response, error)
}

// ReadInterceptor is a function that is called when reading bytes from the
// NPM response body for the specified package. It passes both the current bytes
// read and the total bytes to be read.
//
// This can be used to compute the percentage of a download.
type ReadInterceptor func(pkg PackageRef, current, total int64, bytes []byte)

// Client is a registry client for accessing packages from the registry.
type Client struct {
	// URL is the URL of the registry to connect to. If unspecified, the default
	// registry will be https://packages.simplifier.net.
	URL string

	// Client is the HTTP client to use for making requests. If unspecified, the
	// default client will be http.Default
	Client HTTPClient

	// ReadInterceptors are a list of functions that are called when reading
	// bytes from the NPM response body for the specified package. This can be
	// used to compute the percentage of a download.
	ReadInterceptors []ReadInterceptor
}

// DefaultClient returns a new registry client with the default simplifier
// registry configured, using a non-authenticated client.
var DefaultClient = &Client{}

// Fetch will fetch the given package with the specified version from the
// connected registry.
func (c *Client) Fetch(ctx context.Context, name, version string) (content io.ReadCloser, err error) {
	url := fmt.Sprintf("%s/%s/%s", c.url(), name, version)
	return c.fetchURL(ctx, NewPackageRef(name, version), url, false)
}

func (c *Client) fetchURL(ctx context.Context, ref PackageRef, url string, dist bool) (content io.ReadCloser, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("npm: %w", err)
	}

	resp, err := c.client().Do(req)
	if err != nil {
		return nil, fmt.Errorf("npm: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("%w: %d - %s", ErrStatusCode, resp.StatusCode, resp.Status)
	}

	switch content := resp.Header.Get("Content-Type"); content {
	case "application/gzip", "application/tar+gzip", "application/x-gzip":
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			_ = resp.Body.Close()
			return nil, fmt.Errorf("npm: %w", err)
		}
		adapted := c.adapt(reader, ref, resp.ContentLength, resp.Body)

		return adapted, nil

	case "application/tar", "application/x-tar":
		adapted := c.adapt(resp.Body, ref, resp.ContentLength)

		return adapted, nil

	case "application/json": // indirect
		// Disallow manifests that point to other manifests. Manifests should
		// only point to tarballs.
		if dist {
			return nil, fmt.Errorf("%w: %s", ErrBadContentType, content)
		}
		var pkg struct {
			Dist struct {
				Tarball      string `json:"tarball"`
				UnpackedSize int64  `json:"unpackedSize"`
			} `json:"dist"`
		}
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&pkg); err != nil {
			return nil, fmt.Errorf("%w: %q", ErrBadManifest, err)
		}

		if pkg.Dist.Tarball == "" {
			return nil, fmt.Errorf("%w - missing dist.tarball", ErrBadManifest)
		}

		return c.fetchURL(ctx, ref, pkg.Dist.Tarball, true)
	}
	return nil, fmt.Errorf("%w: %s", ErrBadContentType, content)
}

func (c *Client) url() string {
	if c.URL == "" {
		return "https://packages.simplifier.net"
	}
	return c.URL
}

func (c *Client) client() HTTPClient {
	if c.Client == nil {
		return http.DefaultClient
	}
	return c.Client
}

func (c *Client) adapt(r io.ReadCloser, pkg PackageRef, total int64, closers ...io.Closer) io.ReadCloser {
	if len(c.ReadInterceptors) == 0 {
		return r
	}

	return &readInterceptor{
		pkg:     pkg,
		fns:     c.ReadInterceptors,
		total:   total,
		r:       r,
		current: 0,
		closers: closers,
	}
}

type readInterceptor struct {
	pkg     PackageRef
	fns     []ReadInterceptor
	current int64
	total   int64
	r       io.ReadCloser
	closers []io.Closer
}

func (r *readInterceptor) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	r.current += int64(n)
	for _, fn := range r.fns {
		fn(r.pkg, r.current, r.total, p[:n])
	}
	return
}

func (r *readInterceptor) Close() error {
	errs := make([]error, 0, len(r.closers)+1)

	errs = append(errs, r.r.Close())
	for _, closer := range r.closers {
		errs = append(errs, closer.Close())
	}
	return errors.Join(errs...)
}
