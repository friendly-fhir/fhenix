package templatefuncs_test

import (
	"testing"

	"github.com/friendly-fhir/fhenix/internal/templatefuncs"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestJSON_Encode(t *testing.T) {
	testCases := []struct {
		name    string
		input   any
		want    string
		wantErr error
	}{
		{
			name:  "empty",
			input: nil,
			want:  "null",
		}, {
			name: "basic object",
			input: map[string]any{
				"key": "value",
			},
			want: `{"key":"value"}`,
		}, {
			name: "bad object reports error",
			input: map[string]any{
				"key": func() {},
			},
			want:    templatefuncs.StringOnError,
			wantErr: cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotErr error
			m := &templatefuncs.JSONModule{
				Reporter: templatefuncs.ReporterFunc(func(err error) {
					gotErr = err
				}),
			}

			got := m.Encode(tc.input)

			if got, want := gotErr, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("JSONModule.Encode(...) error = %v, want %v", got, want)
			}
			if got != tc.want {
				t.Errorf("JSONModule.Encode(...) = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestJSON_EncodeIndent(t *testing.T) {
	testCases := []struct {
		name    string
		input   any
		want    string
		wantErr error
	}{
		{
			name:  "empty",
			input: nil,
			want:  "null",
		}, {
			name: "basic object",
			input: map[string]any{
				"key": "value",
			},
			want: `{
  "key": "value"
}`,
		}, {
			name: "bad object reports error",
			input: map[string]any{
				"key": func() {},
			},
			want:    templatefuncs.StringOnError,
			wantErr: cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotErr error
			m := &templatefuncs.JSONModule{
				Reporter: templatefuncs.ReporterFunc(func(err error) {
					gotErr = err
				}),
			}

			got := m.EncodeIndent(tc.input)

			if got, want := gotErr, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("JSONModule.EncodeIndent(...) error = %v, want %v", got, want)
			}
			if got != tc.want {
				t.Errorf("JSONModule.EncodeIndent(...) = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestJSON_Decode(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    map[string]any
		wantErr error
	}{
		{
			name:  "empty",
			input: "{}",
			want:  map[string]any{},
		}, {
			name:  "basic object",
			input: `{"key":"value"}`,
			want: map[string]any{
				"key": "value",
			},
		}, {
			name:    "bad object reports error",
			input:   `{"key":}`,
			wantErr: cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotErr error
			m := &templatefuncs.JSONModule{
				Reporter: templatefuncs.ReporterFunc(func(err error) {
					gotErr = err
				}),
			}

			got := m.Decode(tc.input)

			if got, want := gotErr, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("JSONModule.Decode(...) error = %v, want %v", got, want)
			}
			if diff := cmp.Diff(got, tc.want, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("JSONModule.Decode(...) mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
