package definition_test

import (
	"encoding/json"
	"io/fs"
	"os"
	"testing"

	"github.com/friendly-fhir/fhenix/pkg/model/conformance/definition"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func mustParseJSON[T any](t *testing.T, bytes []byte) *T {
	t.Helper()

	var out T
	if err := json.Unmarshal(bytes, &out); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	return &out
}

func mustReadJSON[T any](t *testing.T, path string) *T {
	t.Helper()

	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	return mustParseJSON[T](t, bytes)
}

func TestReadFile(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		want    any
		wantErr error
	}{
		{
			name: "structure definition",
			path: "testdata/structure-definition.json",
			want: mustReadJSON[definition.StructureDefinition](t, "testdata/structure-definition.json"),
		}, {
			name: "value set",
			path: "testdata/value-set.json",
			want: mustReadJSON[definition.ValueSets](t, "testdata/value-set.json"),
		}, {
			name: "code system",
			path: "testdata/code-system.json",
			want: mustReadJSON[definition.CodeSystem](t, "testdata/code-system.json"),
		}, {
			name: "concept map",
			path: "testdata/concept-map.json",
			want: mustReadJSON[definition.ConceptMap](t, "testdata/concept-map.json"),
		}, {
			name:    "invalid path",
			path:    "testdata/invalid.json",
			wantErr: fs.ErrNotExist,
		}, {
			name:    "not a resource",
			path:    "testdata/not-resource.json",
			wantErr: cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := definition.FromFile(tc.path)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("FromFile(%q) = error %v, want %v", tc.path, got, want)
			}
			if got, want := got, tc.want; !cmp.Equal(got, want) {
				t.Fatalf("FromFile(%q) = %v, want %v", tc.path, got, want)
			}
		})
	}
}
