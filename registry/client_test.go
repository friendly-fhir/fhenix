package registry_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	_ "embed"

	"github.com/friendly-fhir/fhenix/registry"
	"github.com/friendly-fhir/fhenix/registry/registrytest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	//go:embed testdata/good-archive.tar.gz
	goodArchive []byte

	//go:embed testdata/invalid-archive.tar.gz
	invalidArchive []byte

	//go:embed testdata/malformed-archive.tar.gz
	malformedArchive []byte
)

func TestClient_Fetch(t *testing.T) {
	server := registrytest.NewFakeServer()

	server.SetIndirectTarball("good.indirect", "4.0.1", bytes.NewReader(goodArchive))
	server.SetIndirectTarball("malformed.indirect", "4.0.1", bytes.NewReader(malformedArchive))
	server.SetIndirectTarball("invalid.indirect", "4.0.1", bytes.NewReader(invalidArchive))

	server.SetTarball("good.direct", "5.0.1", bytes.NewReader(goodArchive))
	server.SetTarball("malformed.direct", "5.0.1", bytes.NewReader(malformedArchive))
	server.SetTarball("invalid.direct", "5.0.1", bytes.NewReader(invalidArchive))

	server.SetStatusCode("bad.not-found", "4.0.4", 404)
	server.SetContent("bad.content-type", "4.0.4", "application/slam-poetry", bytes.NewReader(nil))
	server.SetContent("bad.json-content", "4.0.4", "application/json", bytes.NewReader([]byte(`{"pkg":`)))
	server.SetContent("bad.no-tarball", "4.0.4", "application/json", bytes.NewReader([]byte(`{}`)))
	server.SetError("bad.error", "4.0.4", fmt.Errorf("server error"))

	ctx := context.Background()
	client, err := registry.NewClient(ctx, registry.URL(server.URL()))
	if err != nil {
		t.Fatal("NewClient() failed unexpected:", err)
	}
	_ = client

	testCases := []struct {
		name      string
		pkg       string
		version   string
		wantBytes int64
		wantErr   error
	}{
		{
			name:      "good indirect",
			pkg:       "good.indirect",
			version:   "4.0.1",
			wantBytes: int64(len(goodArchive)),
		}, {
			name:      "good direct",
			pkg:       "good.direct",
			version:   "5.0.1",
			wantBytes: int64(len(goodArchive)),
		}, {
			name:      "malformed indirect",
			pkg:       "malformed.indirect",
			version:   "4.0.1",
			wantBytes: int64(len(malformedArchive)),
		}, {
			name:      "malformed direct",
			pkg:       "malformed.direct",
			version:   "5.0.1",
			wantBytes: int64(len(malformedArchive)),
		}, {
			name:      "invalid indirect",
			pkg:       "invalid.indirect",
			version:   "4.0.1",
			wantBytes: int64(len(invalidArchive)),
		}, {
			name:      "invalid direct",
			pkg:       "invalid.direct",
			version:   "5.0.1",
			wantBytes: int64(len(invalidArchive)),
		}, {
			name:    "non-200 status code",
			pkg:     "bad.not-found",
			version: "4.0.4",
			wantErr: registry.ErrStatusCode,
		}, {
			name:    "bad content type",
			pkg:     "bad.content-type",
			version: "4.0.4",
			wantErr: registry.ErrBadContentType,
		}, {
			name:    "bad json content",
			pkg:     "bad.json-content",
			version: "4.0.4",
			wantErr: registry.ErrBadContent,
		}, {
			name:    "no tarball field",
			pkg:     "bad.no-tarball",
			version: "4.0.4",
			wantErr: registry.ErrBadContent,
		}, {
			name:    "server error",
			pkg:     "bad.error",
			version: "4.0.4",
			wantErr: cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, count, err := client.Fetch(ctx, tc.pkg, tc.version)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("Fetch() err = %v, want err %v", got, want)
			}
			if got, want := count, tc.wantBytes; got != want {
				t.Errorf("Fetch() count = %v, want count %v", got, want)
			}
			if r != nil {
				defer r.Close()
			}
		})
	}
}
