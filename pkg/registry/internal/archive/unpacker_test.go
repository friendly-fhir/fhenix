package archive_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/friendly-fhir/fhenix/pkg/registry/internal/archive"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestDiskUnpacker_Unpack(t *testing.T) {
	testCases := []struct {
		name     string
		tee      func(name string, r io.Reader) io.Reader
		content  []byte
		wantErr  error
		wantFile bool
	}{
		{
			name:     "good content",
			content:  []byte("hello, world"),
			wantFile: true,
		},
	}

	path := filepath.Join("foo", "bar", "baz.txt")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			unpacker := &archive.DiskUnpacker{
				Root: t.TempDir(),
				Tee:  tc.tee,
			}
			out := filepath.Join(unpacker.Root, path)

			err := unpacker.Unpack(path, int64(len(tc.content)), bytes.NewReader(tc.content))

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("Unpack() error = %v, want %v", got, want)
			}
			if got, want := fileExists(out), tc.wantFile; got != want {
				t.Fatalf("Unpack() file exists = %v, want %v", got, want)
			}
		})
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
