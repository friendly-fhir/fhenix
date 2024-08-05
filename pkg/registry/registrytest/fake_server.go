package registrytest

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"path/filepath"
)

// FakeServer is a server that can be spawned up for testing registry clients.
type FakeServer struct {
	mux    *http.ServeMux
	server *httptest.Server
}

// URL returns the URL of the fake server, which can be used in the client for
// testing.
func (f *FakeServer) URL() string {
	return f.server.URL
}

// NewFakeServer creates a new fake server for testing.
//
// The server is not started by default, and must be started by calling the
// Start function.
func NewFakeServer() *FakeServer {
	mux := http.NewServeMux()
	return &FakeServer{
		mux:    mux,
		server: httptest.NewServer(mux),
	}
}

// Close shuts down the fake server.
func (f *FakeServer) Close() {
	f.server.Close()
}

// SetIndirectGzipTarball adds a package to the fake server that redirects to
// another URL for the tarball content.
func (f *FakeServer) SetIndirectGzipTarball(name, version string, content []byte) {
	f.setIndirect(name, version, "application/tar+gzip", content)
}

// SetIndirectTarball adds a package to the fake server that redirects to
func (f *FakeServer) SetIndirectTarball(name, version string, content []byte) {
	f.setIndirect(name, version, "application/tar", content)
}

// SetIndirectTarballFS sets the indirect tarball response for the given package
// and version, done as a filesystem.
func (f *FakeServer) SetIndirectTarballFS(name, version string, fs fs.FS) {
	f.SetIndirectTarball(name, version, TarballBytes(fs))
}

func (f *FakeServer) setIndirect(name, version, contentType string, content []byte) {
	redirectPattern := fmt.Sprintf("/%s/-/%s-%s.tar.gz", name, filepath.Base(name), version)
	redirect := fmt.Sprintf("%s%s", f.server.URL, redirectPattern)

	reader := bytes.NewReader(content)
	br := bytes.NewBuffer(nil)
	sha := sha1.New()

	tee := io.TeeReader(reader, sha)
	if _, err := io.Copy(br, tee); err != nil {
		panic(err)
	}
	shasum := sha.Sum(nil)

	f.mux.HandleFunc(fmt.Sprintf("/%s/%s", name, version), func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		var manifest struct {
			Name    string `json:"name"`
			Version string `json:"version"`
			Dist    struct {
				Shasum  string `json:"shasum"`
				TarBall string `json:"tarball"`
			} `json:"dist"`
		}
		manifest.Name = name
		manifest.Version = version
		manifest.Dist.TarBall = redirect
		manifest.Dist.Shasum = fmt.Sprintf("%x", shasum)

		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(manifest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	f.setContent(redirectPattern, contentType, br.Bytes())
}

// SetTarball adds a package to the fake server.
func (f *FakeServer) SetTarball(name, version string, bytes []byte) {
	f.setContent(fmt.Sprintf("/%s/%s", name, version), "application/tar", bytes)
}

// SetGzipTarball adds a gzipped tarball to the fake server.
func (f *FakeServer) SetGzipTarball(name, version string, bytes []byte) {
	f.setContent(fmt.Sprintf("/%s/%s", name, version), "application/tar+gzip", bytes)
}

// SetTarballFS sets the tarball response for the given package and version,
// done as a filesystem.
func (f *FakeServer) SetTarballFS(name, version string, fs fs.FS) {
	f.SetTarball(name, version, TarballBytes(fs))
}

// SetContent sets the content response that will be returned by the server
// for the given package.
func (f *FakeServer) SetContent(name, version, contentType string, bytes []byte) {
	f.setContent(fmt.Sprintf("/%s/%s", name, version), contentType, bytes)
}

func (f *FakeServer) setContent(pattern, contentType string, content []byte) {
	f.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		w.Header().Add("Content-Type", contentType)

		reader := bytes.NewReader(content)
		if _, err := io.Copy(w, reader); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

// SetStatusCode makes the fake server return an error for the given package.
func (f *FakeServer) SetStatusCode(name, version string, code int) {
	f.mux.HandleFunc(fmt.Sprintf("/%s/%s", name, version), func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(code), code)
	})
}

// SetError makes the fake server return an error for the given package.
func (f *FakeServer) SetError(name, version string, err error) {
	f.mux.HandleFunc(fmt.Sprintf("/%s/%s", name, version), func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	})
}
