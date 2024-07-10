/*
Package filter provides a mechanism for filtering files to be included or
excluded from processing, and leverages the [config.TransformFilter] for its
base definition.
*/
package filter

import (
	"html/template"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/friendly-fhir/fhenix/config"
	"github.com/friendly-fhir/fhenix/internal/templatefuncs"
	"github.com/friendly-fhir/fhenix/model"
)

// Filter represents a filter that can be applied to a set of definitions.
type Filter struct {
	config *config.TransformFilter
}

// New creates a new filter from the given configuration.
func newFilter(cfg *config.TransformFilter) *Filter {
	return &Filter{
		config: cfg,
	}
}

// Matches returns true if the given value matches the filter.
func (f *Filter) Matches(v any) bool {
	switch v := v.(type) {
	case *model.Type:
		return f.MatchesType(v)
	}
	return false
}

var zero config.TransformFilter

// MatchesType returns true if the given type matches the filter.
func (f *Filter) MatchesType(t *model.Type) bool {
	if t == nil || f.config == nil {
		return false
	}
	if tp := f.config.Type; tp != "" && tp != "StructureDefinition" {
		return false
	}
	if name := f.config.Name; name != "" && !f.match(name, t.Name) {
		return false
	}
	if source := f.config.Source; source != "" && !f.match(source, filepath.Base(t.Source.File)) {
		return false
	}
	if pkg := f.config.Package; pkg != "" && pkg != t.Source.Package.Name() {
		return false
	}
	if url := f.config.URL; url != "" && url != t.URL {
		return false
	}
	if condition := f.config.Condition; condition != "" && !f.evaluateTemplate(condition, t) {
		return false
	}
	return *f.config != zero
}

func (f *Filter) match(regex, needle string) bool {
	got, err := regexp.MatchString(strings.TrimSpace(regex), needle)
	return err == nil && got
}

func (f *Filter) evaluateTemplate(condition string, v any) bool {
	tmpl := template.New("").Funcs(templatefuncs.DefaultFuncs)
	_, err := tmpl.Parse(strings.TrimSpace(condition))
	if err != nil {
		return false
	}

	var sb strings.Builder
	err = tmpl.Execute(&sb, v)
	if err != nil {
		return false
	}

	b, err := strconv.ParseBool(sb.String())
	return b || err != nil
}

// Filters represents a set of filters.
type Filters []*Filter

// NewFilters creates a new set of filters from the given configurations.
func New(cfgs ...*config.TransformFilter) Filters {
	var filters Filters
	for _, cfg := range cfgs {
		filters = append(filters, newFilter(cfg))
	}
	return filters
}

// Matches returns true if the given value matches any of the filters.
func (f Filters) Matches(v any) bool {
	for _, filter := range f {
		if filter.Matches(v) {
			return true
		}
	}
	return false
}

// MatchesType returns true if the given type matches any of the filters.
func (f Filters) MatchesType(t *model.Type) bool {
	for _, filter := range f {
		if filter.MatchesType(t) {
			return true
		}
	}
	return false
}
