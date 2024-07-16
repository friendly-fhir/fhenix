package snek

// Option is a configuration option for a Snek [Application].
type Option interface {
	set(*config)
}

type option func(*config)

func (o option) set(a *config) {
	o(a)
}

var _ Option = (*option)(nil)

type config struct {
	ExitPanic int
	ExitError int

	ApplicationName string

	UsageTemplate   string
	HelpTemplate    string
	VersionTemplate string
	PanicTemplate   string
}

func WithUsageTemplate(tmpl string) Option {
	return option(func(c *config) {
		c.UsageTemplate = tmpl
	})
}

func WithHelpTemplate(tmpl string) Option {
	return option(func(c *config) {
		c.HelpTemplate = tmpl
	})
}

func WithVersionTemplate(tmpl string) Option {
	return option(func(c *config) {
		c.VersionTemplate = tmpl
	})
}

func WithPanicTemplate(tmpl string) Option {
	return option(func(c *config) {
		c.PanicTemplate = tmpl
	})
}

func WithPanicExitCode(code int) Option {
	return option(func(c *config) {
		c.ExitPanic = code
	})
}

func WithErrorExitCode(code int) Option {
	return option(func(c *config) {
		c.ExitError = code
	})
}
