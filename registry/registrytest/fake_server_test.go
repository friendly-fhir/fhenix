package registrytest_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/friendly-fhir/fhenix/registry/registrytest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/iancoleman/strcase"
)

func TestFakeServer(t *testing.T) {
	t.Parallel()

	sut := registrytest.NewFakeServer()
	defer sut.Close()
	testErr := errors.New("test error")
	content := []byte("good package")
	sut.SetTarball("good.tar", "1.0.0", content)
	sut.SetTarballFS("good.fs.tar", "1.0.0", TestPackage)
	sut.SetGzipTarball("good.tar.gzip", "1.0.0", content)

	sut.SetIndirectTarball("good.indirect.tar", "1.0.0", content)
	sut.SetIndirectTarballFS("good.fs.indirect.tar", "1.0.0", TestPackage)
	sut.SetIndirectGzipTarball("good.indirect.tar.gzip", "1.0.0", content)

	sut.SetError("bad.package", "1.0.0", testErr)
	sut.SetStatusCode("bad.not-found", "1.0.0", 404)
	sut.SetContent("bad.content-type", "1.0.0", "application/slam-poetry", nil)

	testCases := []struct {
		name            string
		method          string
		path            string
		wantStatusCode  int
		wantContentType string
		wantContent     []byte
		wantErr         error
	}{
		{
			name:            "good tar package",
			method:          http.MethodGet,
			path:            "/good.tar/1.0.0",
			wantContentType: "application/tar",
			wantStatusCode:  http.StatusOK,
			wantContent:     content,
		}, {
			name:            "good tar gzip package",
			method:          http.MethodGet,
			path:            "/good.tar.gzip/1.0.0",
			wantContentType: "application/tar+gzip",
			wantStatusCode:  http.StatusOK,
			wantContent:     content,
		}, {
			name:            "good tar fs package",
			method:          http.MethodGet,
			path:            "/good.fs.tar/1.0.0",
			wantContentType: "application/tar",
			wantStatusCode:  http.StatusOK,
			wantContent:     registrytest.TarballBytes(TestPackage),
		}, {
			name:            "good indirect tar package",
			method:          http.MethodGet,
			path:            "/good.indirect.tar/1.0.0",
			wantContentType: "application/json",
			wantStatusCode:  http.StatusOK,
		}, {
			name:            "good indirect tar gzip package",
			method:          http.MethodGet,
			path:            "/good.indirect.tar.gzip/1.0.0",
			wantContentType: "application/json",
			wantStatusCode:  http.StatusOK,
		}, {
			name:           "bad package",
			method:         http.MethodGet,
			path:           "/bad.package/1.0.0",
			wantStatusCode: http.StatusInternalServerError,
		}, {
			name:           "not found",
			method:         http.MethodGet,
			path:           "/bad.not-found/1.0.0",
			wantStatusCode: http.StatusNotFound,
		}, {
			name:           "bad method",
			method:         http.MethodPost,
			path:           "/good.tar/1.0.0",
			wantStatusCode: http.StatusMethodNotAllowed,
		}, {
			name:           "bad method indirect",
			method:         http.MethodPost,
			path:           "/good.indirect.tar/1.0.0",
			wantStatusCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := sut.URL() + tc.path
			req, err := http.NewRequest(tc.method, url, nil)
			if err != nil {
				t.Fatalf("http.NewRequest() failed: %v", err)
			}

			resp, err := http.DefaultClient.Do(req)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("http.%s(%q) = error %v, want %v", strcase.ToCamel(tc.method), url, err, tc.wantErr)
			}
			if got, want := resp.StatusCode, tc.wantStatusCode; got != want {
				t.Errorf("http.%s(%q) = %d, want %d", strcase.ToCamel(tc.method), url, got, want)
			}
			if got, want := resp.Header.Get("Content-Type"), tc.wantContentType; want != "" && got != want {
				t.Errorf("http.%s(%q) = %q, want %q", strcase.ToCamel(tc.method), url, got, want)
			}
			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("io.ReadAll() failed: %v", err)
			}

			if want := tc.wantContent; len(want) != 0 && !bytes.Equal(got, want) {
				t.Errorf("http.%s(%q) = %v, want %v", strcase.ToCamel(tc.method), url, got, want)
			}
		})
	}
}
