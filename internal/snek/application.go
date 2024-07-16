package snek

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"text/template"
	"unicode"

	"github.com/friendly-fhir/fhenix/internal/snek/dedent"
	"github.com/spf13/cobra"
)

var (
	//go:embed templates/panic.tmpl
	defaultPanicTemplate string

	//go:embed templates/help.tmpl
	defaultHelpTemplate string

	//go:embed templates/version.tmpl
	defaultVersionTemplate string

	//go:embed templates/usage.tmpl
	defaultUsageTemplate string
)

// Application is a wrapper around a cobra command that represents an
// application.
type Application struct {
	name string

	command *cobra.Command

	exitPanic int

	panicTemplate *template.Template
}

// NewApplication creates a new application from the given root command.
// This function will panic if any of the commands in the root do not define
// a Command.Info function that returns non-nil information.
func NewApplication(name string, root Command, opts ...Option) *Application {
	cfg := &config{
		ExitPanic: 2,
		ExitError: 1,

		ApplicationName: name,

		PanicTemplate:   defaultPanicTemplate,
		UsageTemplate:   defaultUsageTemplate,
		HelpTemplate:    defaultHelpTemplate,
		VersionTemplate: defaultVersionTemplate,
	}
	for _, opt := range opts {
		opt.set(cfg)
	}

	cobra.AddTemplateFuncs(funcs)

	return &Application{
		name:      name,
		command:   toCobraCommand(cfg, root),
		exitPanic: 2,
	}
}

// SetPanicTemplate sets the template that will be used to render the panic
// output.
//
// This function will panic if the template is invalid.
func (a *Application) SetPanicTemplate(tmpl string) {
	t, err := template.New("panic").Funcs(funcs).Parse(tmpl)
	if err != nil {
		panic(err)
	}
	a.panicTemplate = t
}

// SetPanicExitCode sets the exit code that the application will use when
// a panic occurs. The default is 2.
func (a *Application) SetPanicExitCode(code int) {
	a.exitPanic = code
}

// Execute runs the application with the given context.
func (a *Application) Execute(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			a.handlePanic(a.command.ErrOrStderr(), r)
			os.Exit(a.exitPanic)
		}
	}()

	a.command.ExecuteContext(ctx)
}

func (a *Application) handlePanic(w io.Writer, r any) {
	stack := string(debug.Stack())

	p := struct {
		Error      string
		StackTrace []string
	}{
		Error:      fmt.Sprintf("%v", r),
		StackTrace: strings.Split(stack, "\n"),
	}

	err := a.panicTemplate.Execute(w, p)
	if err != nil {
		panic(err)
	}
	os.Exit(a.exitPanic)
}

func toCobraCommand(cfg *config, command Command) *cobra.Command {
	info := command.Info()
	if info == nil {
		panic("Command.Info() must return non-nil information")
	}

	result := &cobra.Command{
		Use:     info.Use,
		Short:   info.Summary,
		Long:    info.Description,
		Aliases: info.Aliases,
		Example: dedent.Strings(info.Examples...),
		Version: info.Version,

		Hidden: info.Hidden,

		DisableAutoGenTag: true,

		Annotations: map[string]string(info.Annotations),

		ValidArgsFunction: toCompletionFunc(command),

		Args: command.PositionalArgs().positionArg(),

		RunE: toRunFunc(cfg, command),
	}

	// Install the command flags and completion functions
	flagsets := command.Flags()
	for _, fs := range flagsets {
		result.Flags().AddFlagSet(fs.FlagSet())

		for name, completion := range fs.CompletionFuncs() {
			result.RegisterFlagCompletionFunc(name, toCompletionFunc(completion))
		}
	}
	commandFlags.Set(result, flagsets)

	// Install sub-commands and apply the appropriate groups
	for group, cmds := range command.Commands() {
		id := toID(group)
		result.AddGroup(&cobra.Group{
			ID:    id,
			Title: group,
		})

		for _, cmd := range *cmds {
			cmd := toCobraCommand(cfg, cmd)

			cmd.GroupID = id
			result.AddCommand(cmd)
		}
	}

	result.SetErrPrefix(errPrefix(cfg.ApplicationName))
	result.SetHelpTemplate(cfg.HelpTemplate)
	result.SetUsageTemplate(cfg.UsageTemplate)
	result.SetVersionTemplate(cfg.VersionTemplate)

	return result
}

func toCompletionFunc(completer Completer) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		got := completer.Complete(args, toComplete)
		if got == nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return got.completion()
	}
}

func toID(name string) string {
	replaceSpecial := func(r rune) rune {
		if !unicode.IsLetter(r) {
			return '_'
		}
		return r
	}
	return strings.Map(replaceSpecial, name)
}

func toRunFunc(cfg *config, command Command) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		// set application context information
		ctx = withAppName(ctx, cfg.ApplicationName)
		ctx = withStdout(ctx, cmd.OutOrStdout())
		ctx = withStderr(ctx, cmd.ErrOrStderr())

		err := command.Run(ctx, args)
		if err != nil {
			if err == errNotImplemented {
				return cmd.Usage()
			} else if IsUsageError(err) {
				return err
			}

			Errorf(ctx, "%v", err)
			os.Exit(cfg.ExitError)
		}
		return nil
	}
}
