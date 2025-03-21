package npmtest

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"

	"github.com/friendly-fhir/fhenix/pkg/registryv2/internal/npm"
)

// Server is a test server for an npm registry.
// This models a real registry server by allowing [fs.FS] to be served as
// tarballs.
type Server struct {
	server *httptest.Server
	mux    *http.ServeMux
	client *npm.Client
}

// NewServer creates a new test server for the npm registry.
func NewServer() *Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	return &Server{
		server: server,
		mux:    mux,
		client: &npm.Client{
			URL:    server.URL,
			Client: server.Client(),
		},
	}
}

// Client returns the [npm.Client] for this [Server].
func (s *Server) Client() *npm.Client {
	return s.client
}

// SetMalformedTarball sets the package and version to return a badly formatted
// gzipped tarball.
func (s *Server) SetMalformedTarball(pkg, version string) {
	s.mux.HandleFunc(fmt.Sprintf("/%v/%v", pkg, version), func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/tar+gzip")
		_, _ = w.Write([]byte("bad gzip tar content"))
	})
}

// SetBadContentType sets the package and version to return a bad content type.
func (s *Server) SetBadContentType(pkg, version string) {
	s.mux.HandleFunc(fmt.Sprintf("/%v/%v", pkg, version), func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/invalid-content-type")
	})
}

// SetDistLoop specifies a package and version that will return a JSON
// payload which has a dist.tarball that redirects to the same package and
// version.
func (s *Server) SetDistLoop(pkg, version string) {
	pattern := fmt.Sprintf("/%v/%v", pkg, version)
	s.setDist(pattern, s.server.URL+pattern, 0)
}

// SetTarball specifies an [fs.FS] that will be served as a tarball for a
// package and version. This is served directly from this URL without any
// JSON indirection.
func (s *Server) SetTarball(pkg, version string, fs fs.FS) {
	s.setTarball(fmt.Sprintf("/%v/%v", pkg, version), fs)
}

// SetGzipTarball specifies an [fs.FS] that will be served as a gzipped
// tarball for a package and version. This is served directly from this URL
// without any JSON indirection.
func (s *Server) SetGzipTarball(pkg, version string, fs fs.FS) {
	s.setGzipTarball(fmt.Sprintf("/%v/%v", pkg, version), fs)
}

// SetManifestTarball specifies an [fs.FS] that will be served as a tarball for
// a package and version. This is served indirectly through a JSON response
// that specifies the dist.tarball URL which actually contains the data.
func (s *Server) SetManifestTarball(pkg, version string, fs fs.FS) {
	tarballPattern := fmt.Sprintf("/%v/%v/tarball", pkg, version)
	s.setTarball(tarballPattern, fs)
	address := s.server.URL + tarballPattern
	size := s.sizeOfFS(fs)

	s.setDist(fmt.Sprintf("/%v/%v", pkg, version), address, size)
}

// SetManifestGzipTarball specifies an [fs.FS] that will be served as a gzipped
// tarball for a package and version. This is served indirectly through a JSON
// response that specifies the dist.tarball URL which actually contains the data.
func (s *Server) SetManifestGzipTarball(pkg, version string, fs fs.FS) error {
	tarballPattern := fmt.Sprintf("/%v/%v/tarball", pkg, version)
	s.setGzipTarball(tarballPattern, fs)
	address := s.server.URL + tarballPattern
	size := s.sizeOfFS(fs)

	s.setDist(fmt.Sprintf("/%v/%v", pkg, version), address, size)
	return nil
}

// SetBadManifest sets the package and version to return a badly formatted JSON
// response.
func (s *Server) SetBadManifest(pkg, version string) {
	s.mux.HandleFunc(fmt.Sprintf("/%v/%v", pkg, version), func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("bad json"))
	})
}

// SetResponseStatus sets the HTTP status code that will be returned for a
// package and version.
func (s *Server) SetResponseStatus(pkg, version string, status int) {
	s.mux.HandleFunc(fmt.Sprintf("/%v/%v", pkg, version), func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, http.StatusText(status), status)
	})
}

// SetEmptyManifest sets the package and version to return an empty manifest.
func (s *Server) SetEmptyManifest(pkg, version string) {
	s.setDist(fmt.Sprintf("/%v/%v", pkg, version), "", 0)
}

func (s *Server) setDist(pattern string, tarball string, unpackedSize int64) {
	var pkg struct {
		Dist struct {
			Tarball      string `json:"tarball"`
			UnpackedSize int64  `json:"unpackedSize"`
		} `json:"dist"`
	}

	pkg.Dist.Tarball = tarball
	pkg.Dist.UnpackedSize = unpackedSize

	s.mux.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(pkg); err != nil {
			// Unreachable?
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func (s *Server) sizeOfFS(f fs.FS) int64 {
	var size int64
	_ = fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.Type().IsRegular() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		size += info.Size()
		return nil
	})
	return size
}

func (s *Server) setTarball(pattern string, fs fs.FS) {
	s.mux.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/tar")

		pr, pw := io.Pipe()
		go func() {
			writer := tar.NewWriter(pw)
			if err := writer.AddFS(fs); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			_ = writer.Close()
			_ = pw.Close()
		}()

		_, _ = io.Copy(w, pr)
		_ = pr.Close()
	})
}

func (s *Server) setGzipTarball(pattern string, fs fs.FS) {
	s.mux.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/tar+gzip")

		pr, pw := io.Pipe()
		go func() {
			gzipWriter := gzip.NewWriter(pw)
			tarWriter := tar.NewWriter(gzipWriter)
			if err := tarWriter.AddFS(fs); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			_ = tarWriter.Close()
			_ = gzipWriter.Close()
			_ = pw.Close()
		}()

		_, _ = io.Copy(w, pr)
		_ = pr.Close()
	})
}
