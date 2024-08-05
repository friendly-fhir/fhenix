package test_test

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/friendly-fhir/fhenix/pkg/registry"
)

var opts = godog.Options{
	Output:      colors.Colored(os.Stdout),
	Concurrency: runtime.NumCPU(),
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	format := "progress"
	verbose := false
	for _, arg := range os.Args[1:] {
		if arg == "-test.v=true" { // go test transforms -v option
			format = "pretty"
			verbose = true
			break
		}
	}

	opts := godog.Options{
		Format:   format,
		Paths:    []string{"testdata/features"},
		TestingT: t,
		Strict:   true,
	}

	binary, init := InitializeTestSuite(t, verbose)
	status := godog.TestSuite{
		Name:                 "Fhenix",
		TestSuiteInitializer: init,
		ScenarioInitializer:  InitializeScenario(t, verbose, binary),
		Options:              &opts,
	}.Run()

	if status != 0 {
		t.Fatal("Non-zero status")
	}
}

func repoRootDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		t.Fatalf("repoRootDir(): unable to get caller information")
	}
	return filepath.Clean(filepath.Dir(filepath.Dir(file)))
}

func run(t *testing.T, verbose bool, cmd string, args ...string) (string, error) {
	t.Helper()
	if verbose {
		t.Logf("calling %s %s", cmd, strings.Join(args, " "))
	}
	var combined strings.Builder
	command := exec.Command(cmd, args...)
	if verbose {
		command.Stdout = io.MultiWriter(os.Stdout, &combined)
		command.Stderr = io.MultiWriter(os.Stderr, &combined)
	} else {
		command.Stdout = &combined
		command.Stderr = &combined
	}

	return combined.String(), command.Run()
}

func InitializeTestSuite(t *testing.T, verbose bool) (binary string, init func(ctx *godog.TestSuiteContext)) {
	t.Helper()
	dir := t.TempDir()
	binDir := filepath.Join(dir, "bin")
	binary = filepath.Join(binDir, "fhenix")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	init = func(ctx *godog.TestSuiteContext) {
		ctx.BeforeSuite(func() {
			_ = os.MkdirAll(binDir, 0755)
			if verbose {
				t.Logf("Building fhenix...")
			}
			_, err := run(t, verbose, "go", "build", "-o", binary, repoRootDir(t))
			if err != nil {
				t.Fatalf("failed to build fhenix: %v", err)
			}
		})
		ctx.AfterSuite(func() {
			if verbose {
				t.Logf("removing %s...", binary)
			}
			_ = os.RemoveAll(binary)
		})
	}
	return
}

func InitializeScenario(t *testing.T, verbose bool, binary string) func(ctx *godog.ScenarioContext) {
	dir := t.TempDir()
	return func(ctx *godog.ScenarioContext) {
		ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
			t.Setenv("FHIR_CACHE", dir)
			if verbose {
				t.Logf("FHIR_CACHE=%s", dir)
			}
			return ctx, nil
		})
		ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
			if verbose {
				t.Logf("removing %s...", dir)
			}
			_ = os.RemoveAll(dir)
			return ctx, err
		})
		cache := registry.NewCache(dir)

		tc := &TestCase{
			Timeout: 3 * time.Minute,
			cache:   cache,
			Verbose: verbose,
			T:       t,
		}

		ctx.Given(`^a timeout of (\d+[a-zA-Z]+)$`, tc.SetTimeout)
		ctx.Given(`^the cache is empty$`, tc.PurgeCache)
		ctx.Given(`^the cache contains packages:$`, tc.PrimeCache)
		ctx.Given(`^the command '(.+)'$`, tc.SetCommand(binary))
		ctx.Given(`^the registry is '(.+)'$`, tc.SetRegistry)

		ctx.When(`^the command is executed$`, tc.ExecuteCommand)

		ctx.Then(`^the FHIR cache contains packages:$`, tc.HasPackages)
		ctx.Then(`^the FHIR cache does not contain packages:$`, tc.DoesNotHavePackages)
		ctx.Then(`^the program exits with status-code (\d+)$`, tc.StatusCodeIs)
		ctx.Then(`^the program exits with non-(\d+) status-code$`, tc.StatusCodeIsNot)
		ctx.Then(`^stdout contains '(.+)'$`, tc.StdoutContains)
		ctx.Then(`^stderr contains '(.+)'$`, tc.StderrContains)
		ctx.Then(`^stdout does not contain '(.+)'$`, tc.StdoutDoesNotContain)
		ctx.Then(`^stderr does not contain '(.+)'$`, tc.StderrDoesNotContain)
		ctx.Then(`^stdout is '(.+)'$`, tc.StdoutIs)
		ctx.Then(`^stderr is '(.+)'$`, tc.StderrIs)
		ctx.Then(`^stdout is not '(.+)'$`, tc.StdoutIsNot)
		ctx.Then(`^stderr is not '(.+)'$`, tc.StderrIsNot)
	}
}

type TestCase struct {
	// Testing
	T       *testing.T
	Verbose bool

	// Setup
	Command  []string
	Timeout  time.Duration
	Registry string

	// Command Output
	ExitCode int
	Stdout   string
	Stderr   string

	// Local Cache
	cache *registry.Cache
}

func (ts *TestCase) SetTimeout(timeout string) error {
	d, err := time.ParseDuration(timeout)
	if err != nil {
		return err
	}
	ts.Timeout = d
	return nil
}

func (tc *TestCase) PurgeCache() error {
	path := tc.cache.Root()
	return os.RemoveAll(path)
}

func (tc *TestCase) PrimeCache(table *godog.Table) error {
	return godog.ErrPending
}

func (tc *TestCase) SetCommand(binary string) func(string) error {
	return func(command string) error {
		parts := strings.Fields(command)
		if parts[0] == "fhenix" {
			parts[0] = binary
		}
		tc.Command = parts
		return nil
	}
}

func (tc *TestCase) SetRegistry(registry string) error {
	tc.Registry = registry
	return nil
}

func (tc *TestCase) ExecuteCommand() error {
	if len(tc.Command) == 0 {
		return godog.ErrPending
	}
	cmd := tc.Command
	if tc.Timeout != 0 {
		cmd = append(cmd, "--timeout", tc.Timeout.String())
	}
	if tc.Registry != "" {
		cmd = append(cmd, "--registry", tc.Registry)
	}
	cmd = append(cmd, "--fhir-cache", tc.cache.Root())
	return tc.runCommand(cmd[0], cmd[1:]...)
}

func (tc *TestCase) runCommand(cmd string, args ...string) error {
	if tc.Verbose {
		tc.T.Logf("calling %s %s", cmd, strings.Join(args, " "))
	}

	var stdout strings.Builder
	var stderr strings.Builder
	command := exec.Command(cmd, args...)
	if tc.Verbose {
		command.Stdout = io.MultiWriter(os.Stdout, &stdout)
		command.Stderr = io.MultiWriter(os.Stderr, &stderr)
	} else {
		command.Stdout = &stdout
		command.Stderr = &stderr
	}
	err := command.Run()

	tc.Stdout = stdout.String()
	tc.Stderr = stderr.String()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			tc.ExitCode = exiterr.ExitCode()
		} else {
			return err
		}
	}
	return nil
}

func (tc *TestCase) HasPackages(packages *godog.Table) error {
	var errs []error
	for _, row := range packages.Rows[1:] {
		name := row.Cells[0].Value
		version := row.Cells[1].Value

		if !tc.cache.Contains("default", name, version) {
			errs = append(errs, fmt.Errorf("expected package %s@%s not found in cache", name, version))
		}
	}
	return errors.Join(errs...)
}

func (tc *TestCase) DoesNotHavePackages(packages *godog.Table) error {
	var errs []error
	for _, row := range packages.Rows[1:] {
		name := row.Cells[0].Value
		version := row.Cells[1].Value

		if tc.cache.Contains("default", name, version) {
			errs = append(errs, fmt.Errorf("unexpected package %s@%s found in cache", name, version))
		}
	}
	return errors.Join(errs...)
}

func (tc *TestCase) StatusCodeIs(code int) error {
	if tc.ExitCode != code {
		return fmt.Errorf("expected exit code %d, got %d", code, tc.ExitCode)
	}
	return nil
}

func (tc *TestCase) StatusCodeIsNot(code int) error {
	if tc.ExitCode == code {
		return fmt.Errorf("expected exit code not to be %d", code)
	}
	return nil
}

func (tc *TestCase) StdoutContains(expected string) error {
	if !strings.Contains(tc.Stdout, expected) {
		return fmt.Errorf("expected stdout to contain %q, got %q", expected, tc.Stdout)
	}
	return nil
}

func (tc *TestCase) StderrContains(expected string) error {
	if !strings.Contains(tc.Stderr, expected) {
		return fmt.Errorf("expected stderr to contain %q, got %q", expected, tc.Stderr)
	}
	return nil
}

func (tc *TestCase) StdoutDoesNotContain(expected string) error {
	if strings.Contains(tc.Stdout, expected) {
		return fmt.Errorf("expected stdout not to contain %q, got %q", expected, tc.Stdout)
	}
	return nil
}

func (tc *TestCase) StderrDoesNotContain(expected string) error {
	if strings.Contains(tc.Stderr, expected) {
		return fmt.Errorf("expected stderr not to contain %q, got %q", expected, tc.Stderr)
	}
	return nil
}

func (tc *TestCase) StdoutIs(expected string) error {
	if tc.Stdout != expected {
		return fmt.Errorf("expected stdout to be %q, got %q", expected, tc.Stdout)
	}
	return nil
}

func (tc *TestCase) StderrIs(expected string) error {
	if tc.Stderr != expected {
		return fmt.Errorf("expected stderr to be %q, got %q", expected, tc.Stderr)
	}
	return nil
}

func (tc *TestCase) StdoutIsNot(expected string) error {
	if tc.Stdout == expected {
		return fmt.Errorf("expected stdout not to be %q, got %q", expected, tc.Stdout)
	}
	return nil
}

func (tc *TestCase) StderrIsNot(expected string) error {
	if tc.Stderr != expected {
		return fmt.Errorf("expected stderr not to be %q, got %q", expected, tc.Stderr)
	}
	return nil
}
