package transformer

import (
	"os"
	"strings"
	texttemplate "text/template"

	"github.com/friendly-fhir/fhenix/internal/templatefuncs"
)

func NewFunc(path string, reporter templatefuncs.Reporter) (func(...any) string, error) {
	fntmpl := texttemplate.New("").Funcs(templatefuncs.NewFuncs(reporter))

	var err error
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if fntmpl, err = fntmpl.Parse(string(bytes)); err != nil {
		return nil, err
	}
	return func(data ...any) string {
		var in any
		if len(data) == 1 {
			in = data[0]
		}
		var sb strings.Builder
		_ = fntmpl.Execute(&sb, in)
		result := sb.String()
		return strings.TrimSpace(result)
	}, nil
}

func FuncsFromConfig(funcs map[string]string, reporter templatefuncs.Reporter) (map[string]any, error) {
	result := make(map[string]any, len(funcs))
	for name, path := range funcs {
		fn, err := NewFunc(path, reporter)
		if err != nil {
			return nil, err
		}
		result[name] = fn
	}
	return result, nil
}
