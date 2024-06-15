package template_test

import (
	"testing"

	"github.com/friendly-fhir/fhenix/internal/template"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gopkg.in/yaml.v3"
)

func nodeFrom(t *testing.T, input string) *yaml.Node {
	t.Helper()

	var node yaml.Node
	if err := yaml.Unmarshal([]byte(input), &node); err != nil {
		t.Fatalf("yaml.Unmarshal: %v", err)
	}

	return &node
}

func TestTemplate_UnmarshalYAML(t *testing.T) {
	testCases := []struct {
		name    string
		input   *yaml.Node
		wantErr error
	}{
		{
			name:  "valid",
			input: nodeFrom(t, `"{{- .Name -}}"`),
		}, {
			name:    "invalid template format",
			input:   nodeFrom(t, `" {{ .Name | | }} "`),
			wantErr: cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got template.Template
			err := tc.input.Decode(&got)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("UnmarshalYAML(%s) = %v, want = %v", tc.input.Value, got, want)
			}
		})
	}
}

func TestTemplate_ExecuteBool(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		template *template.Template
		want     bool
		wantErr  error
	}{
		{
			name:     "bad execution",
			input:    "",
			template: template.MustParse("template", `{{ .Name }}`),
			wantErr:  cmpopts.AnyError,
		},
		{
			name:     "boolean true",
			input:    true,
			template: template.MustParse("template", `{{ . }}`),
			want:     true,
		}, {
			name:     "boolean false",
			input:    false,
			template: template.MustParse("template", `{{ . }}`),
			want:     false,
		}, {
			name:     "integer 0",
			input:    0,
			template: template.MustParse("template", `{{ . }}`),
			want:     false,
		}, {
			name:     "integer 42",
			input:    42,
			template: template.MustParse("template", `{{ . }}`),
			want:     true,
		}, {
			name:     "string empty",
			input:    "",
			template: template.MustParse("template", `{{ . }}`),
			want:     false,
		}, {
			name:     "non-empty non-json content",
			input:    "hello world",
			template: template.MustParse("template", `{{ . }}`),
			want:     true,
		}, {
			name:     "empty quoted string",
			input:    `""`,
			template: template.MustParse("template", `{{ . }}`),
			want:     false,
		}, {
			name:     "empty list node",
			input:    `[]`,
			template: template.MustParse("template", `{{ . }}`),
			want:     false,
		}, {
			name:     "non-empty list node",
			input:    `["hello", "world"]`,
			template: template.MustParse("template", `{{ . }}`),
			want:     true,
		}, {
			name:     "empty object",
			input:    `{}`,
			template: template.MustParse("template", `{{ . }}`),
			want:     false,
		}, {
			name:     "non-empty object",
			input:    `{"hello": "world"}`,
			template: template.MustParse("template", `{{ . }}`),
			want:     true,
		}, {
			name:     "quoted string",
			input:    `"hello world"`,
			template: template.MustParse("template", `{{ . }}`),
			want:     true,
		}, {
			name:     "null node",
			input:    `null`,
			template: template.MustParse("template", `{{ . }}`),
			want:     false,
		}, {
			name:     "falsey non-json",
			input:    "FALSE",
			template: template.MustParse("template", `{{ . }}`),
			want:     false,
		}, {
			name:     "truthy non-json",
			input:    "TRUE",
			template: template.MustParse("template", `{{ . }}`),
			want:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.template.ExecuteBool(tc.input)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("ExecuteBool(%v) = %v, want = %v", tc.input, got, want)
			}
			if got, want := got, tc.want; got != want {
				t.Errorf("ExecuteBool(%v) = %v, want = %v", tc.input, got, want)
			}
		})
	}
}

func TestTemplate_Funcs(t *testing.T) {
	testCases := []struct {
		name    string
		expr    string
		want    string
		wantErr error
	}{
		{
			name: "lowercase transforms to lowercase",
			expr: `{{ lowercase "Hello World!" }}`,
			want: "hello world!",
		}, {
			name: "lowercase transforms to lowercase pipeline",
			expr: `{{ "Hello World!" | lowercase }}`,
			want: "hello world!",
		}, {
			name: "uppercase transforms to uppercase",
			expr: `{{ uppercase "Hello World!" }}`,
			want: "HELLO WORLD!",
		}, {
			name: "uppercase transforms to uppercase pipeline",
			expr: `{{ "Hello World!" | uppercase }}`,
			want: "HELLO WORLD!",
		}, {
			name: "title transforms to title case",
			expr: `{{ titlecase "hello world!" }}`,
			want: "Hello World!",
		}, {
			name: "title transforms to title case pipeline",
			expr: `{{ "hello world!" | titlecase }}`,
			want: "Hello World!",
		}, {
			name: "pascalcase transforms to pascal case",
			expr: `{{ pascalcase "hello-world!" }}`,
			want: "HelloWorld",
		}, {
			name: "pascalcase transforms to pascal case pipeline",
			expr: `{{ "hello-world!" | pascalcase }}`,
			want: "HelloWorld",
		}, {
			name: "camelcase transforms to pascal case",
			expr: `{{ camelcase "hello-world!" }}`,
			want: "helloWorld",
		}, {
			name: "camelcase transforms to pascal case pipeline",
			expr: `{{ "hello-world!" | camelcase }}`,
			want: "helloWorld",
		}, {
			name:    "too many inputs",
			expr:    `{{ "hello" | titlecase "Hello World!" }}`,
			wantErr: cmpopts.AnyError,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpl := template.MustParse("template", tc.expr)

			got, err := tmpl.ExecuteString(nil)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("ExecuteString(nil) = %v, want = %v", got, want)
			}
			if got, want := got, tc.want; got != want {
				t.Errorf("ExecuteString(nil) = %v, want = %v", got, want)
			}
		})
	}
}
