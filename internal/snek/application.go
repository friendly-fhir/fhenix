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

	"github.com/friendly-fhir/fhenix/internal/dedent"
	"github.com/spf13/cobra"
)

var (
	//go:embed templates/panic.tmpl
	panicTemplateString string

	//go:embed templates/help.tmpl
	helpTemplateString string

	//go:embed templates/version.tmpl
	versionTemplateString string

	//go:embed templates/usage.tmpl
	usageTemplateString string

	panicTemplate = template.Must(template.New("panic").Funcs(funcs).Parse(panicTemplateString))
)

const (
	ExitSuccess = 0
	ExitError   = 1
	ExitPanic   = 2
)

// Application is a wrapper around a cobra command that represents an
// application.
type Application struct {
	name string

	appinfo *AppInfo

	command *cobra.Command

	panicTemplate *template.Template
}

type AppInfo struct {
	// Name is the name of the application.
	Name string

	// Website is the URL for the contact person or organization responsible for
	// the application.
	Website string

	// GitHubRepository is the URL for the GitHub repository for the application.
	GitHubRepository string

	// IssueURL is the URL for the issue tracker for the application.
	IssueURL string
}

// NewApplication creates a new application from the given root command.
// This function will panic if any of the commands in the root do not define
// a Command.Info function that returns non-nil information.
func NewApplication(root Command, appinfo *AppInfo) *Application {
	info := root.Info()
	if info == nil {
		panic("Command.Info() must return non-nil information")
	}
	if appinfo == nil {
		panic("meta must not be nil")
	}
	if appinfo.Name == "" {
		appinfo.Name = strings.Split(info.Use, " ")[0]
	}

	cfg := &config{
		ApplicationName: appinfo.Name,

		UsageTemplate:   usageTemplateString,
		HelpTemplate:    helpTemplateString,
		VersionTemplate: versionTemplateString,
	}

	cobra.AddTemplateFuncs(funcs)

	return &Application{
		name:          appinfo.Name,
		command:       toCobraCommand(cfg, root),
		appinfo:       appinfo,
		panicTemplate: panicTemplate,
	}
}

// Execute runs the application with the given context.
func (a *Application) Execute(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			Panicf(ctx, "%v", r)
			a.handlePanic(a.command.ErrOrStderr(), r)
			os.Exit(ExitPanic)
		}
	}()

	return a.command.ExecuteContext(ctx)
}

func (a *Application) handlePanic(w io.Writer, r any) {
	stack := string(debug.Stack())

	p := struct {
		Error      string
		StackTrace []string
		Meta       *AppInfo
	}{
		Error:      fmt.Sprintf("%v", r),
		StackTrace: strings.Split(stack, "\n"),
		Meta:       a.appinfo,
	}
	tmpl := a.panicTemplate
	if tmpl == nil {
		tmpl = panicTemplate
	}

	err := tmpl.Execute(w, p)
	if err != nil {
		panic(err)
	}
	os.Exit(ExitPanic)
}

// SetOut sets the out writer for the application.
func (a *Application) SetOut(w io.Writer) {
	a.visitAllCommands(func(cmd *cobra.Command) {
		cmd.SetOut(w)
	})
}

// SetErr sets the error writer for the application.
func (a *Application) SetErr(w io.Writer) {
	a.visitAllCommands(func(cmd *cobra.Command) {
		cmd.SetErr(w)
	})
}

// SetHelpTemplate sets the template that will be used to render the help output.
func (a *Application) SetHelpTemplate(tmpl string) {
	a.visitAllCommands(func(cmd *cobra.Command) {
		cmd.SetHelpTemplate(tmpl)
	})
}

// SetUsageTemplate sets the template that will be used to render the usage
// output.
func (a *Application) SetUsageTemplate(tmpl string) {
	a.visitAllCommands(func(cmd *cobra.Command) {
		cmd.SetUsageTemplate(tmpl)
	})
}

// SetVersionTemplate sets the template that will be used to render the version
// output.
func (a *Application) SetVersionTemplate(tmpl string) {
	a.visitAllCommands(func(cmd *cobra.Command) {
		cmd.SetVersionTemplate(tmpl)
	})
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

func (a *Application) visitAllCommands(fn func(cmd *cobra.Command)) {
	a.visitCommand(a.command, fn)
}

func (a *Application) visitCommand(cmd *cobra.Command, fn func(cmd *cobra.Command)) {
	fn(cmd)
	for _, cmd := range cmd.Commands() {
		a.visitCommand(cmd, fn)
	}
}

type config struct {
	ApplicationName string

	UsageTemplate   string
	HelpTemplate    string
	VersionTemplate string
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
			// This can only error if a flag is either not registered, or already
			// registered -- but this is only populated if the flag exists, and only
			// set once, here, while building.
			_ = result.RegisterFlagCompletionFunc(name, toCompletionFunc(completion))
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
			os.Exit(ExitError)
		}
		return nil
	}
}
