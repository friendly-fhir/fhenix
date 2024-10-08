package snek

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"regexp"
	"runtime/debug"
	"strings"
	"text/template"
	"time"
	"unicode"

	"atomicgo.dev/cursor"
	"github.com/friendly-fhir/fhenix/internal/ansi"
	"github.com/friendly-fhir/fhenix/internal/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
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

	appInfo *AppInfo

	command *cobra.Command

	panicTemplate *template.Template

	buildDate time.Time
}

type AppInfo struct {
	// Name is the name of the application.
	Name string

	// Website is the URL for the contact person or organization responsible for
	// the application.
	Website string

	// RepositoryURL is the URL for the repository for the application.
	RepositoryURL string

	// ReportIssueURL is the URL for the issue tracker for the application.
	ReportIssueURL string

	// DocsURL is the URL for the documentation for the application, which may be
	// different from the base Website URL.
	DocsURL string

	// KeyTerms are words that are important in the command, and will be
	// highlighted in the help output. (Optional)
	KeyTerms []string

	// Variables are words that are variable inputs in the command, and will be
	// highlighted in the help output. (Optional)
	Variables []string

	// UserData is an optional (unspecified) piece of data that can be provided
	// to the application. Default is nil.
	UserData any

	// Version is the version of the application. Default is "snapshot" if
	// not specified
	Version string

	// Date is the date the application was built at.
	Date time.Time
}

// NewApplication creates a new application from the given root command.
// This function will panic if any of the commands in the root do not define
// a Command.Info function that returns non-nil information.
func NewApplication(root Command, appInfo *AppInfo) *Application {
	info := root.Info()
	if info == nil {
		panic("Command.Info() must return non-nil information")
	}
	if appInfo == nil {
		panic("meta must not be nil")
	}
	if appInfo.Name == "" {
		appInfo.Name = strings.Split(info.Use, " ")[0]
	}

	cfg := &config{
		ApplicationName: appInfo.Name,

		UsageTemplate:   usageTemplateString,
		HelpTemplate:    helpTemplateString,
		VersionTemplate: versionTemplateString,

		KeyTerms:  appInfo.KeyTerms,
		Variables: appInfo.Variables,
	}

	cobra.AddTemplateFuncs(funcs)
	cobra.AddTemplateFuncs(template.FuncMap{
		"AppName": get(appInfo.Name),
		"AppInfo": get(appInfo),
	})
	cmd := toCobraCommand(cfg, root)
	cmd.Version = appInfo.Version
	if cmd.Version == "" {
		cmd.Version = "snapshot"
	}
	if appInfo.Date.IsZero() {
		appInfo.Date = time.Now()
	}

	return &Application{
		name:          appInfo.Name,
		command:       cmd,
		appInfo:       appInfo,
		panicTemplate: panicTemplate,
		buildDate:     appInfo.Date,
	}
}

func get[T any](v T) func() T {
	return func() T {
		return v
	}
}

// Execute runs the application with the given context.
func (a *Application) Execute(ctx context.Context) (code *StatusCode) {
	defer func() {
		if r := recover(); r != nil {
			Panicf(ctx, "%v", r)
			a.handlePanic(a.command.ErrOrStderr(), r)

			code = &StatusCode{
				Result: PanicError(fmt.Sprintf("%v", r)),
				Code:   ExitPanic,
			}
		}
	}()
	defer cursor.Show()

	err := a.command.ExecuteContext(ctx)
	status := ExitSuccess
	if err != nil {
		status = ExitError
	}

	return &StatusCode{
		Result: err,
		Code:   status,
	}
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
		Meta:       a.appInfo,
	}
	tmpl := a.panicTemplate
	if tmpl == nil {
		tmpl = panicTemplate
	}

	err := tmpl.Execute(w, p)
	if err != nil {
		panic(err)
	}
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
	t := template.New("panic").Funcs(funcs).Funcs(template.FuncMap{
		"AppName": get(a.name),
		"AppInfo": get(a.appInfo),
	})

	t, err := t.Parse(tmpl)
	if err != nil {
		panic(err)
	}
	a.panicTemplate = t
}

// GenManTree generates man pages for the application into the given directory.
func (a *Application) GenManTree(directory string) error {
	header := &doc.GenManHeader{
		Title: a.appInfo.Name,
		Date:  &a.buildDate,
	}
	return doc.GenManTree(a.command, header, directory)
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

	KeyTerms  []string
	Variables []string

	ShowCursor bool
}

type set[T comparable] map[T]struct{}

func (s *set[T]) Add(v T) {
	if *s == nil {
		*s = make(map[T]struct{})
	}
	(*s)[v] = struct{}{}
}

func setOf[T comparable](vs ...T) set[T] {
	s := make(map[T]struct{})
	for _, v := range vs {
		s[v] = struct{}{}
	}
	return s
}

func (s *set[T]) Insert(vs ...T) set[T] {
	if *s == nil {
		*s = make(map[T]struct{})
	}
	for _, v := range vs {
		(*s)[v] = struct{}{}
	}
	return *s
}

func (s *set[T]) Contains(v T) bool {
	if *s == nil {
		return false
	}
	_, ok := (*s)[v]
	return ok
}

func highlight(s string, content set[string], color ansi.Display) string {
	for k := range content {
		regex := regexp.MustCompile(fmt.Sprintf("(?i)%s", k))
		s = string(regex.ReplaceAllFunc([]byte(s), func(b []byte) []byte {
			return []byte(color.Format(string(b)))
		}))
	}
	return s
}

var urlRegex = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)

func highlightURLs(s string) string {
	return string(urlRegex.ReplaceAllFunc([]byte(s), func(b []byte) []byte {
		return []byte(FormatLink.Format(string(b)))
	}))
}

func highlightAll(s string, keyterms set[string], variables set[string]) string {
	s = highlight(s, keyterms, FormatKeyword)
	s = highlight(s, variables, FormatArg)
	s = highlightURLs(s)
	return s
}

func toCobraCommand(cfg *config, command Command) *cobra.Command {
	info := command.Info()
	if info == nil {
		panic("Command.Info() must return non-nil information")
	}

	keyTerms := setOf(cfg.KeyTerms...)
	keyTerms.Insert(info.KeyTerms...)

	variables := setOf(cfg.Variables...)
	variables.Insert(info.Variables...)

	result := &cobra.Command{
		Use:     info.Use,
		Short:   highlightAll(info.Summary, keyTerms, variables),
		Long:    highlightAll(info.Description, keyTerms, variables),
		Aliases: info.Aliases,
		Example: highlightURLs(dedent.Strings(info.Examples...)),
		Version: info.Version,

		Hidden: info.Hidden,

		DisableAutoGenTag: true,
		SilenceUsage:      true,
		SilenceErrors:     true,

		Annotations: map[string]string(info.Annotations),

		ValidArgsFunction: toCompletionFunc(command),

		Args: command.PositionalArgs().positionArg(),

		RunE: toRunFunc(cfg, info, command),
	}
	result.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		Errorf(c.Context(), "%v", err)
		_ = c.Usage()
		return err
	})

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

		for _, transform := range fs.transforms {
			transform(result)
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

func toRunFunc(cfg *config, info *CommandInfo, command Command) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		// set application context information
		ctx = withAppName(ctx, cfg.ApplicationName)
		ctx = withStdout(ctx, cmd.OutOrStdout())
		ctx = withStderr(ctx, cmd.ErrOrStderr())

		hideCursor := !info.ShowCursor && IsTerminal(cmd.OutOrStdout())
		if hideCursor {
			cursor.Hide()
			defer cursor.Show()
		}

		err := command.Run(ctx, args)
		if err != nil {
			if err == errNotImplemented {
				return cmd.Usage()
			}
			Errorf(ctx, "%v", err)
			if IsUsageError(err) {
				_ = cmd.Usage()
			}
			return err
		}
		return nil
	}
}
