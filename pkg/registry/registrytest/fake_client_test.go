package registrytest_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/friendly-fhir/fhenix/pkg/registry"
	"github.com/friendly-fhir/fhenix/pkg/registry/registrytest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestFakeClient_Fetch(t *testing.T) {
	client := registrytest.NewFakeClient()
	testErr := errors.New("test error")
	content := []byte("good package")
	client.SetTarball("good.tar", "1.0.0", content)
	client.SetGzipTarball("good.tar.gzip", "1.0.0", registrytest.GzipTarballBytes(TestPackage))
	client.SetTarballFS("good.fs.tar", "1.0.0", TestPackage)
	client.SetError("bad.package", "1.0.0", testErr)
	client.Set("bad.not-found", "1.0.0", http.StatusNotFound, nil)

	testCases := []struct {
		name      string
		pkg       string
		version   string
		wantBytes int64
		wantErr   error
	}{
		{
			name:      "good tar",
			pkg:       "good.tar",
			version:   "1.0.0",
			wantBytes: int64(len(content)),
		}, {
			name:      "good gzip tar",
			pkg:       "good.tar.gzip",
			version:   "1.0.0",
			wantBytes: int64(len(registrytest.GzipTarballBytes(TestPackage))),
		}, {
			name:      "good tar fs",
			pkg:       "good.fs.tar",
			version:   "1.0.0",
			wantBytes: int64(len(registrytest.TarballBytes(TestPackage))),
		}, {
			name:    "bad package",
			pkg:     "bad.package",
			version: "1.0.0",
			wantErr: testErr,
		}, {
			name:    "not found",
			pkg:     "bad.not-found",
			version: "1.0.0",
			wantErr: registry.ErrStatusCode,
		}, {
			name:    "package not configured returns 404",
			pkg:     "not.configured",
			version: "1.0.0",
			wantErr: registry.ErrStatusCode,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, length, err := client.Fetch(context.Background(), tc.pkg, tc.version)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("FakeClient.Fetch(%q,%q) = error %v, want %v", tc.pkg, tc.version, err, tc.wantErr)
			}
			if got, want := length, tc.wantBytes; got != want {
				t.Errorf("FakeClient.Fetch(%q,%q) = %d bytes, want %d", tc.pkg, tc.version, got, want)
			}
		})
	}
}
