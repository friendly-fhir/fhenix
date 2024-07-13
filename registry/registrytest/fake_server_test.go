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
	sut := registrytest.NewFakeServer()
	defer sut.Close()
	testErr := errors.New("test error")
	content := []byte("good package")
	sut.SetTarball("good.package", "1.0.0", bytes.NewReader(content))
	sut.SetIndirectTarball("good.indirect.package", "1.0.0", bytes.NewReader(content))
	sut.SetError("bad.package", "1.0.0", testErr)
	sut.SetStatusCode("bad.not-found", "1.0.0", 404)
	sut.SetContent("bad.content-type", "1.0.0", "application/slam-poetry", bytes.NewReader(nil))

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
			name:            "good package",
			method:          http.MethodGet,
			path:            "/good.package/1.0.0",
			wantContentType: "application/gzip",
			wantStatusCode:  http.StatusOK,
			wantContent:     content,
		}, {
			name:            "good indirect package",
			method:          http.MethodGet,
			path:            "/good.indirect.package/1.0.0",
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
			path:           "/good.package/1.0.0",
			wantStatusCode: http.StatusMethodNotAllowed,
		}, {
			name:           "bad method indirect",
			method:         http.MethodPost,
			path:           "/good.indirect.package/1.0.0",
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
