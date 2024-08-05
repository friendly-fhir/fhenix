package registrytest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"sync"

	"github.com/friendly-fhir/fhenix/pkg/registry"
	"github.com/friendly-fhir/fhenix/pkg/registry/internal/auth"
)

// FakeClient is a fake registry client which can be used for local-only
// testing, without needing to spawn up an http server.
//
// This makes it easier to write tests that require a client without needing
// to worry about network connectivity.
//
// FakeClient is thread-safe/goroutine-safe.
type FakeClient struct {
	*registry.Client
	client *authClient
}

// NewFakeClient creates a new fake client for testing.
func NewFakeClient() *FakeClient {
	result := &FakeClient{
		client: &authClient{
			entries: &sync.Map{},
		},
	}
	client, err := registry.NewClient(
		context.Background(),
		registry.URL(""),
		registry.Auth(auth.AuthenticationFunc(func(ctx context.Context) (auth.Client, error) {
			return result.client, nil
		})),
	)
	if err != nil {
		// If this happens, this library is at fault for not constructing a valid
		// client.
		panic(err)
	}
	result.Client = client
	return result
}

// SetGzipTarball sets the gzip response for the given package and version.
func (fc *FakeClient) SetGzipTarball(name, version string, content []byte) {
	entry := &contentEntry{
		responseCode: http.StatusOK,
		content:      content,
		contentType:  "application/tar+gzip",
		length:       int64(len(content)),
	}
	fc.client.entries.Store(fmt.Sprintf("/%s/%s", name, version), entry)
}

// SetTarball sets the tar response for the given package and version.
func (fc *FakeClient) SetTarball(name, version string, content []byte) {
	entry := &contentEntry{
		responseCode: http.StatusOK,
		content:      content,
		contentType:  "application/tar",
		length:       int64(len(content)),
	}
	fc.client.entries.Store(fmt.Sprintf("/%s/%s", name, version), entry)
}

// SetTarballFS sets the tarball respons for the given package and version, done as
// a filesystem
func (fc *FakeClient) SetTarballFS(name, version string, fs fs.FS) {
	fc.SetTarball(name, version, TarballBytes(fs))
}

// SetError sets the error response for the given package and version.
func (fc *FakeClient) SetError(name, version string, err error) {
	entry := &contentEntry{
		responseCode: http.StatusInternalServerError,
		err:          err,
	}
	fc.client.entries.Store(fmt.Sprintf("/%s/%s", name, version), entry)
}

// Set sets the exact response/status-code for the given package and version.
func (fc *FakeClient) Set(name, version string, statusCode int, content []byte) {
	entry := &contentEntry{
		responseCode: statusCode,
		content:      content,
		length:       int64(len(content)),
	}
	fc.client.entries.Store(fmt.Sprintf("/%s/%s", name, version), entry)
}

// authClient is a fake [auth.Client] that retains a cache of responses on-hand
// rather than doing any real communication. This makes fake responses that act
// like they came from a server, without requirin the real connectivity.
type authClient struct {
	entries *sync.Map
}

type contentEntry struct {
	responseCode int
	content      []byte
	contentType  string
	length       int64
	err          error
}

func (ac *authClient) Do(req *http.Request) (*http.Response, error) {
	entry, ok := ac.entries.Load(req.URL.String())
	if !ok {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       nil,
			Header:     make(http.Header),
		}, nil
	}

	content := entry.(*contentEntry)
	if content.err != nil {
		return nil, content.err
	}
	headers := http.Header{}
	headers.Add("Content-Type", content.contentType)
	return &http.Response{
		StatusCode:    content.responseCode,
		Status:        http.StatusText(content.responseCode),
		Body:          io.NopCloser(bytes.NewReader(content.content)),
		Header:        headers,
		ContentLength: content.length,
		ProtoMajor:    1,
		ProtoMinor:    0,
	}, nil
}

var _ auth.Client = (*authClient)(nil)
