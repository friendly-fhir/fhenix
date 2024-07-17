package snek

import "github.com/spf13/cobra"

// Completion is an interface that can be implemented by commands to provide
// custom completions for arguments.
type Completion interface {
	completion() ([]string, cobra.ShellCompDirective)
}

type completion func() ([]string, cobra.ShellCompDirective)

func (c completion) completion() ([]string, cobra.ShellCompDirective) {
	return c()
}

var _ Completion = (*completion)(nil)

// CompletionError indicates that an error occurred while trying to
// provide completions, and the shell should not provide any completions.
var CompletionError = completion(func() ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveError
})

// CompleteWithoutSpace indicates that the shell should not add a space
// after the completion even if there is only a single completion provided.
func CompleteWithoutSpace(completions ...string) Completion {
	return completion(func() ([]string, cobra.ShellCompDirective) {
		return completions, cobra.ShellCompDirectiveNoSpace
	})
}

// CompleteFileExts indicates that the provided completions should be used as
// file extension filters.
//
// For flags, using Command.MarkFlagFilename() and Command.MarkPersistentFlagFilename()
// is a shortcut to using this directive explicitly.  The BashCompFilenameExt
// annotation can also be used to obtain the same behavior for flags.
func CompleteFileExts(extensions ...string) Completion {
	return completion(func() ([]string, cobra.ShellCompDirective) {
		return extensions, cobra.ShellCompDirectiveFilterFileExt
	})
}

var NoCompletion = CompleteNoFiles()

// CompleteNoFiles indicates that the shell should not provide
// file completion even when no completion is provided.
func CompleteNoFiles(completions ...string) Completion {
	return completion(func() ([]string, cobra.ShellCompDirective) {
		return completions, cobra.ShellCompDirectiveNoFileComp
	})
}

// CompleteDirs indicates that only directory names should
// be provided in file completion.  To request directory names within another
// directory, the returned completions should specify the directory within
// which to search.  The BashCompSubdirsInDir annotation can be used to
// obtain the same behavior but only for flags.
func CompleteDirs(completions ...string) Completion {
	return completion(func() ([]string, cobra.ShellCompDirective) {
		return completions, cobra.ShellCompDirectiveFilterDirs
	})
}

// Completer is an interface that can be implemented by commands to provide
// custom completions for arguments.
type Completer interface {
	// Complete returns a list of possible completions for the current set of args
	// and argument being completed.
	Complete(args []string, toComplete string) Completion
}

// CompleterFunc is a function that can be used to implement the Completer
// interface.
type CompleterFunc func(args []string, toComplete string) Completion

// Complete calls the underlying function to provide completions.
func (f CompleterFunc) Complete(args []string, toComplete string) Completion {
	return completion(func() ([]string, cobra.ShellCompDirective) {
		return f(args, toComplete).completion()
	})
}

var _ Completer = (*CompleterFunc)(nil)
