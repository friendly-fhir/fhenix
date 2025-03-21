package npm_test

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/friendly-fhir/fhenix/pkg/registryv2/internal/npm"
	"github.com/friendly-fhir/fhenix/pkg/registryv2/internal/npm/npmtest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestClient_Fetch(t *testing.T) {
	// Arrange
	server := npmtest.NewServer()
	client := server.Client()
	server.SetTarball("tarball", "1.0.0", npmtest.SimplePackageFS())
	server.SetGzipTarball("gzip-tarball", "1.0.0", npmtest.SimplePackageFS())
	server.SetManifestTarball("indirect-tarball", "1.0.0", npmtest.SimplePackageFS())
	server.SetManifestGzipTarball("indirect-gzip-tarball", "1.0.0", npmtest.SimplePackageFS())
	server.SetResponseStatus("error-forbidden", "403", http.StatusForbidden)
	server.SetMalformedTarball("bad-gzip-tar", "1.0.0")
	server.SetBadManifest("bad-manifest", "1.0.0")
	server.SetEmptyManifest("empty-manifest", "1.0.0")
	server.SetBadContentType("bad-content-type", "1.0.0")
	server.SetDistLoop("dist-loop", "1.0.0")

	testCases := []struct {
		name    string
		ctx     *context.Context
		pkg     string
		version string
		want    map[string]string
		wantErr error
	}{
		{
			name:    "direct tarball",
			pkg:     "tarball",
			version: "1.0.0",
			want:    npmtest.FSContents(t, npmtest.SimplePackageFS()),
		}, {
			name:    "direct gzip tarball",
			pkg:     "gzip-tarball",
			version: "1.0.0",
			want:    npmtest.FSContents(t, npmtest.SimplePackageFS()),
		}, {
			name:    "indirect tarball",
			pkg:     "indirect-tarball",
			version: "1.0.0",
			want:    npmtest.FSContents(t, npmtest.SimplePackageFS()),
		}, {
			name:    "indirect gzip tarball",
			pkg:     "indirect-gzip-tarball",
			version: "1.0.0",
			want:    npmtest.FSContents(t, npmtest.SimplePackageFS()),
		}, {
			name:    "server returns bad status code",
			pkg:     "error-forbidden",
			version: "1.0.0",
			wantErr: npm.ErrStatusCode,
		}, {
			name:    "package version does not exist",
			pkg:     "tarball",
			version: "2.0.0",
			wantErr: npm.ErrStatusCode,
		}, {
			name:    "context expired",
			ctx:     cancelledContext(),
			pkg:     "tarball",
			version: "1.0.0",
			wantErr: context.Canceled,
		}, {
			name:    "bad gzip tarball content",
			pkg:     "bad-gzip-tar",
			version: "1.0.0",
			wantErr: cmpopts.AnyError,
		}, {
			name:    "bad manifest",
			pkg:     "bad-manifest",
			version: "1.0.0",
			wantErr: npm.ErrBadManifest,
		}, {
			name:    "empty manifest",
			pkg:     "empty-manifest",
			version: "1.0.0",
			wantErr: npm.ErrBadManifest,
		}, {
			name:    "bad content type",
			pkg:     "bad-content-type",
			version: "1.0.0",
			wantErr: npm.ErrBadContentType,
		}, {
			name:    "nil context",
			ctx:     nilContext(),
			pkg:     "tarball",
			version: "1.0.0",
			wantErr: cmpopts.AnyError,
		}, {
			name:    "dist loop",
			pkg:     "dist-loop",
			version: "1.0.0",
			wantErr: npm.ErrBadContentType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			rc, err := client.Fetch(nilToContext(tc.ctx), tc.pkg, tc.version)

			// Assert
			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("Client.Fetch: want %q, got: %q", want, got)
			}
			if err == nil { // branch to avoid nil pointer dereference
				got := npmtest.TarContents(t, rc)
				_ = rc.Close()

				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("Client.Fetch: (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestClient_ReadInterceptors_CalledWithPackage(t *testing.T) {
	// Arrange
	server := npmtest.NewServer()
	client := server.Client()
	server.SetTarball("tarball", "1.0.0", npmtest.SimplePackageFS())
	server.SetGzipTarball("gzip-tarball", "1.0.0", npmtest.SimplePackageFS())

	testCases := []struct {
		name    string
		pkg     string
		version string
	}{
		{
			name:    "tarball",
			pkg:     "tarball",
			version: "1.0.0",
		}, {
			name:    "gzip-tarball",
			pkg:     "gzip-tarball",
			version: "1.0.0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			wantRef := npm.NewPackageRef(tc.pkg, tc.version)
			var gotRef npm.PackageRef
			client.ReadInterceptors = append(client.ReadInterceptors, func(ref npm.PackageRef, current, total int64, bytes []byte) {
				gotRef = ref
			})
			r, err := client.Fetch(context.Background(), tc.pkg, tc.version)
			if err != nil {
				t.Fatalf("Client.Fetch: %v", err)
			}
			defer r.Close()

			// Act
			_, err = io.ReadAll(r)
			if err != nil {
				t.Fatalf("io.ReadAll: %v", err)
			}

			// Assert
			if got, want := gotRef, wantRef; !cmp.Equal(got, want) {
				t.Errorf("ReadInterceptor: want %v, got: %v", want, got)
			}
		})
	}
}

func nilToContext(ctx *context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return *ctx
}

func cancelledContext() *context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return &ctx
}

func nilContext() *context.Context {
	var ctx context.Context
	return &ctx
}
