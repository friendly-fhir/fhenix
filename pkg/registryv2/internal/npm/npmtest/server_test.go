package npmtest_test

import (
	"compress/gzip"
	"context"
	"net/http"
	"testing"

	"github.com/friendly-fhir/fhenix/pkg/registryv2/internal/npm"
	"github.com/friendly-fhir/fhenix/pkg/registryv2/internal/npm/npmtest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestServer_Client_InvalidAddress(t *testing.T) {
	server := npmtest.NewServer()
	client := server.Client()

	_, err := client.Fetch(context.Background(), "foo", "1.0.0")

	if got, want := err, npm.ErrStatusCode; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
		t.Fatalf("unexpected error: got %v, want %v", got, want)
	}
}

func TestServer_RoundTrip(t *testing.T) {
	server := npmtest.NewServer()
	client := server.Client()
	server.SetTarball("tarball", "1.0.0", npmtest.SimplePackageFS())
	server.SetGzipTarball("gzip-tarball", "1.0.0", npmtest.SimplePackageFS())
	server.SetManifestTarball("indirect-tarball", "1.0.0", npmtest.SimplePackageFS())
	server.SetManifestGzipTarball("indirect-gzip-tarball", "1.0.0", npmtest.SimplePackageFS())
	server.SetResponseStatus("error-forbidden", "1.0.0", http.StatusForbidden)
	server.SetBadManifest("bad-manifest", "1.0.0")
	server.SetEmptyManifest("empty-manifest", "1.0.0")
	server.SetBadContentType("bad-content-type", "1.0.0")
	server.SetMalformedTarball("malformed-tarball", "1.0.0")
	server.SetDistLoop("dist-loop", "1.0.0")

	testCases := []struct {
		name    string
		want    map[string]string
		wantErr error
	}{
		{
			name: "tarball",
			want: npmtest.FSContents(t, npmtest.SimplePackageFS()),
		}, {
			name: "gzip-tarball",
			want: npmtest.FSContents(t, npmtest.SimplePackageFS()),
		}, {
			name: "indirect-tarball",
			want: npmtest.FSContents(t, npmtest.SimplePackageFS()),
		}, {
			name: "indirect-gzip-tarball",
			want: npmtest.FSContents(t, npmtest.SimplePackageFS()),
		}, {
			name:    "error-forbidden",
			wantErr: npm.ErrStatusCode,
		}, {
			name:    "bad-manifest",
			wantErr: npm.ErrBadManifest,
		}, {
			name:    "empty-manifest",
			wantErr: npm.ErrBadManifest,
		}, {
			name:    "bad-content-type",
			wantErr: npm.ErrBadContentType,
		}, {
			name:    "malformed-tarball",
			wantErr: gzip.ErrHeader,
		}, {
			name:    "dist-loop",
			wantErr: npm.ErrBadContentType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rc, err := client.Fetch(context.Background(), tc.name, "1.0.0")
			got := npmtest.TarContents(t, rc)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("Client.Fetch: want %q, got: %q", want, got)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Client.Fetch: (-want +got):\n%s", diff)
			}
		})
	}
}
