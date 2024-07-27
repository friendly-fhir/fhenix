package registry

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"net/http"

	"github.com/friendly-fhir/fhenix/registry/internal/auth"
)

// Client is a registry client for accessing packages from the registry.
type Client struct {
	url    string
	client auth.Client
}

type config struct {
	url  string
	auth Authentication
}

type Option interface {
	set(*config)
}

type option func(*config)

func (o option) set(c *config) {
	o(c)
}

// URL returns an [Option] that sets the URL of the registry client.
// If unspecified, the default registry will be https://packages.simplifier.net
func URL(url string) Option {
	return option(func(c *config) {
		c.url = url
	})
}

// Auth returns an [Option] that sets the authentication method for the client.
func Auth(auth Authentication) Option {
	return option(func(c *config) {
		c.auth = auth
	})
}

// NewClient creates a new registry client with the specified options.
// If no options are provided, the client will be created with the default
// registry for Simplifier.net, using no authentication.
func NewClient(ctx context.Context, opts ...Option) (*Client, error) {
	cfg := config{
		url:  "https://packages.simplifier.net",
		auth: nil,
	}
	for _, opt := range opts {
		opt.set(&cfg)
	}
	if cfg.auth == nil {
		cfg.auth = NoAuthentication()
	}
	client, err := auth.LoadClient(ctx, cfg.auth)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
		url:    cfg.url,
	}, nil
}

// DefaultClient returns a new registry client with the default simplifier
// registry configured, using a non-authenticated client.
var DefaultClient = &Client{
	client: http.DefaultClient,
	url:    "https://packages.simplifier.net",
}

var (
	ErrStatusCode     = fmt.Errorf("unexpected status code")
	ErrNoTarball      = fmt.Errorf("missing tarball URL")
	ErrBadContentType = fmt.Errorf("unexpected content-type")
	ErrBadContent     = fmt.Errorf("bad content")
)

type multiCloser struct {
	reader  io.Reader
	closers []io.Closer
}

func (m *multiCloser) Read(p []byte) (n int, err error) {
	return m.reader.Read(p)
}

func (m *multiCloser) Close() error {
	var errs []error
	for _, r := range m.closers {
		errs = append(errs, r.Close())
	}
	return errors.Join(errs...)
}

func multiReadCloser(reader io.ReadCloser, closers ...io.Closer) io.ReadCloser {
	return &multiCloser{
		reader:  reader,
		closers: append(closers, reader),
	}
}

// Fetch will fetch the given package with the specified version from the
// connected registry.
func (c *Client) Fetch(ctx context.Context, name, version string) (content io.ReadCloser, bytes int64, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s/%s", c.url, name, version), nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("%w: %d - %s", ErrStatusCode, resp.StatusCode, resp.Status)
	}
	switch content := resp.Header.Get("Content-Type"); content {
	case "application/gzip", "application/tar+gzip":
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			resp.Body.Close()
			return nil, 0, err
		}
		return multiReadCloser(reader, resp.Body), resp.ContentLength, nil
	case "application/tar":
		return resp.Body, resp.ContentLength, nil
	case "application/json":
		var pkg struct {
			Dist struct {
				Shasum       string `json:"shasum"`
				Tarball      string `json:"tarball"`
				UnpackedSize int64  `json:"unpackedSize"`
			} `json:"dist"`
		}
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&pkg); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrBadContent, err)
		}

		if pkg.Dist.Tarball == "" {
			return nil, 0, fmt.Errorf("%w: missing tarball URL", ErrBadContent)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, pkg.Dist.Tarball, nil)
		if err != nil {
			return nil, 0, err
		}
		resp, err = c.client.Do(req)
		if err != nil {
			return nil, 0, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, 0, fmt.Errorf("%w: %d - %s", ErrStatusCode, resp.StatusCode, resp.Status)
		}
	default:
		return nil, 0, fmt.Errorf("%w: %s", ErrBadContentType, content)
	}

	return resp.Body, resp.ContentLength, nil
}
