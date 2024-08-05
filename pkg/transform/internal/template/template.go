/*
Package template is an abstraction over the builtin 'template' package by
allowing either HTML or Text to be selected as the basic mechanism.
*/
package template

import (
	"fmt"
	htmltemplate "html/template"
	"io"
	texttemplate "text/template"
)

// Template is the interface used by Go's text and html packages.
type Template interface {
	New(string) Template

	Parse(string) (Template, error)

	Execute(w io.Writer, v any) error

	Funcs(funcs map[string]any) Template
}

// FromString returns a new TemplateEngine based on the given string.
func FromString(str string) (Engine, error) {
	switch str {
	case "text":
		return Text(), nil
	case "html":
		return HTML(), nil
	}
	return nil, fmt.Errorf("unknown template mode: %s", str)
}

// Engine is the interface used by Go's text and html packages.
type Engine interface {
	New(string) Template
}

// Text returns a new text template engine.
func Text() Engine {
	return textTemplateEngine{}
}

// HTML returns a new html template engine.
func HTML() Engine {
	return htmlTemplateEngine{}
}

// Default returns the default template engine.
func Default() Engine {
	return Text()
}

type textTemplate struct {
	*texttemplate.Template
}

func (tt *textTemplate) New(name string) Template {
	return &textTemplate{tt.Template.New(name)}
}

func (tt *textTemplate) Parse(content string) (Template, error) {
	tmpl, err := tt.Template.Parse(content)
	if err != nil {
		return nil, err
	}
	return &textTemplate{tmpl}, nil
}

func (tt *textTemplate) Execute(w io.Writer, v any) error {
	return tt.Template.Execute(w, v)
}

func (tt *textTemplate) Funcs(funcs map[string]any) Template {
	return &textTemplate{tt.Template.Funcs(funcs)}
}

type htmlTemplate struct {
	*htmltemplate.Template
}

func (ht *htmlTemplate) New(name string) Template {
	return &htmlTemplate{ht.Template.New(name)}
}

func (ht *htmlTemplate) Parse(content string) (Template, error) {
	tmpl, err := ht.Template.Parse(content)
	if err != nil {
		return nil, err
	}
	return &htmlTemplate{tmpl}, nil
}

func (ht *htmlTemplate) Execute(w io.Writer, v any) error {
	return ht.Template.Execute(w, v)
}

func (ht *htmlTemplate) Funcs(funcs map[string]any) Template {
	return &htmlTemplate{ht.Template.Funcs(funcs)}
}

var _ Template = (*htmlTemplate)(nil)

type textTemplateEngine struct{}

func (textTemplateEngine) New(name string) Template {
	return &textTemplate{texttemplate.New(name)}
}

var _ Engine = (*textTemplateEngine)(nil)

type htmlTemplateEngine struct{}

func (htmlTemplateEngine) New(name string) Template {
	return &htmlTemplate{htmltemplate.New(name)}
}

var _ Engine = (*htmlTemplateEngine)(nil)
