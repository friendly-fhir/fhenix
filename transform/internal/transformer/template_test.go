package transformer_test

import (
	"io/fs"
	"strings"
	"testing"

	"github.com/friendly-fhir/fhenix/transform/internal/template"
	"github.com/friendly-fhir/fhenix/transform/internal/transformer"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewTemplate(t *testing.T) {
	testCases := []struct {
		name      string
		engine    template.Engine
		templates map[string]string
		wantValue bool
		wantErr   error
	}{
		{
			name:      "empty templates",
			engine:    template.Text(),
			templates: nil,
			wantValue: false,
		}, {
			name:      "template does not exist",
			engine:    template.Text(),
			templates: map[string]string{"main": "testdata/does-not-exist.tmpl"},
			wantErr:   fs.ErrNotExist,
		}, {
			name:      "valid template",
			engine:    template.Text(),
			templates: map[string]string{"main": "testdata/hello.tmpl"},
			wantValue: true,
		}, {
			name:      "invalid template",
			engine:    template.Text(),
			templates: map[string]string{"main": "testdata/bad.tmpl"},
			wantErr:   cmpopts.AnyError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := transformer.NewTemplate(tc.engine, tc.templates)

			if got, want := err, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("NewTemplate() = error %v, want %v", got, want)
			}
			if got == nil && tc.wantValue {
				t.Errorf("NewTemplate() = nil, want non-nil")
			}
		})
	}
}

type FakeModelData struct {
	StructureDefinitions []struct{}
	ValueSets            []struct{}
	CodeSystems          []struct{}
}

func NewFakeModelData(structureDefs, valueSets, codeSystems int) *FakeModelData {
	return &FakeModelData{
		StructureDefinitions: make([]struct{}, structureDefs),
		ValueSets:            make([]struct{}, valueSets),
		CodeSystems:          make([]struct{}, codeSystems),
	}
}

func lines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func normalize(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

func TestTemplateExecute_SanityCheck(t *testing.T) {
	testCases := []struct {
		name      string
		engine    template.Engine
		templates map[string]string
		want      string

		structureDefinitions int
		valueSets            int
		codeSystems          int
	}{
		{
			name:   "simple template invokes once",
			engine: template.Text(),
			templates: map[string]string{
				"header": "testdata/header.tmpl",
				"footer": "testdata/footer.tmpl",
			},
			want: lines(
				"header",
				"footer",
			),
		}, {
			name:   "template with multiple structure definitions",
			engine: template.Text(),
			templates: map[string]string{
				"structure-definition": "testdata/body.tmpl",
				"header":               "testdata/header.tmpl",
				"footer":               "testdata/footer.tmpl",
			},
			want: lines(
				"header",
				"body",
				"body",
				"footer",
			),
			structureDefinitions: 2,
		}, {
			name:   "template with multiple value sets",
			engine: template.Text(),
			templates: map[string]string{
				"value-set": "testdata/body.tmpl",
				"header":    "testdata/header.tmpl",
				"footer":    "testdata/footer.tmpl",
			},
			want: lines(
				"header",
				"body",
				"body",
				"footer",
			),
			valueSets: 2,
		}, {
			name:   "template with multiple code systems",
			engine: template.Text(),
			templates: map[string]string{
				"code-system": "testdata/body.tmpl",
				"header":      "testdata/header.tmpl",
				"footer":      "testdata/footer.tmpl",
			},
			want: lines(
				"header",
				"body",
				"body",
				"footer",
			),
			codeSystems: 2,
		}, {
			name:   "replacing main called once",
			engine: template.Text(),
			templates: map[string]string{
				"main":   "testdata/body.tmpl",
				"header": "testdata/header.tmpl",
				"footer": "testdata/footer.tmpl",
			},
			want: lines(
				"header",
				"body",
				"footer",
			),
			structureDefinitions: 10,
			valueSets:            10,
			codeSystems:          10,
		}, {
			name:   "No header and footer",
			engine: template.Text(),
			templates: map[string]string{
				"main": "testdata/body.tmpl",
			},
			want: lines(
				"body",
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpl, err := transformer.NewTemplate(tc.engine, tc.templates)
			if err != nil {
				t.Fatalf("NewTemplate() = error %v, want nil", err)
			}

			data := NewFakeModelData(tc.structureDefinitions, tc.valueSets, tc.codeSystems)
			var sb strings.Builder
			err = tmpl.Execute(&sb, data)
			if err != nil {
				t.Fatalf("Template.Execute() = error %v", err)
			}

			if got, want := normalize(strings.TrimSpace(sb.String())), tc.want; !cmp.Equal(got, want) {
				t.Errorf("Template.Execute() = %q, want %q", got, want)
			}
		})
	}
}
