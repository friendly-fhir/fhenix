package transformer_test

import (
	"io/fs"
	"testing"

	"github.com/friendly-fhir/fhenix/transform/internal/transformer"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewFunc(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		input   any
		want    string
		wantErr error
	}{
		{
			name:    "invalid template",
			path:    "testdata/bad.tmpl",
			wantErr: cmpopts.AnyError,
		}, {
			name:  "valid template",
			path:  "testdata/uppercase.tmpl",
			input: "hello",
			want:  "HELLO",
		}, {
			name:    "file does not exist",
			path:    "testdata/does-not-exist.tmpl",
			wantErr: fs.ErrNotExist,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fn, err := transformer.NewFunc(tc.path, nil)
			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("NewFunc() = %v, want %v", got, want)
			}

			got := zeroFn(fn)(tc.input)
			if want := tc.want; !cmp.Equal(got, want) {
				t.Errorf("NewFunc() = %v, want %v", got, want)
			}
		})
	}
}

func zeroFn(fn func(...any) string) func(...any) string {
	if fn == nil {
		return func(...any) string { return "" }
	}
	return fn
}
