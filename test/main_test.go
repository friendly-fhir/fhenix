package test_test

import (
	"context"
	"errors"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/friendly-fhir/fhenix/registry"
)

var opts = godog.Options{
	Output:      colors.Colored(os.Stdout),
	Concurrency: runtime.NumCPU(),
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

func TestIntegration(t *testing.T) {
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

	status := godog.TestSuite{
		Name:                 "Fhenix",
		TestSuiteInitializer: InitializeTestSuite(t),
		ScenarioInitializer:  InitializeScenario(t, verbose),
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

func run(t *testing.T, cmd string, args ...string) {
	t.Helper()
	t.Logf("Running %s %s", cmd, strings.Join(args, " "))
	output, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}
	t.Logf(string(output))
}

func InitializeTestSuite(t *testing.T) func(ctx *godog.TestSuiteContext) {
	dir := t.TempDir()
	return func(ctx *godog.TestSuiteContext) {
		ctx.BeforeSuite(func() {
			t.Setenv("GOPATH", dir)
			t.Setenv("PATH",
				strings.Join([]string{os.Getenv("PATH"), filepath.Join(dir, "bin")},
					string(filepath.ListSeparator)),
			)
			t.Logf("Installing fhenix application...")
			run(t, "go", "install", repoRootDir(t))
		})
		ctx.AfterSuite(func() {
			t.Setenv("GOPATH", dir)
			t.Logf("Cleaning up...")
			run(t, "go", "clean", "-cache", "-modcache", "-testcache")
			_ = os.RemoveAll(dir)
		})
	}
}

func InitializeScenario(t *testing.T, verbose bool) func(ctx *godog.ScenarioContext) {
	dir := t.TempDir()
	return func(ctx *godog.ScenarioContext) {
		ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
			t.Setenv("FHIR_CACHE", dir)
			return ctx, nil
		})
		ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
			return ctx, errors.Join(err, os.RemoveAll(dir))
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
		ctx.Given(`^the command '(.+)'$`, tc.SetCommand)
		ctx.Given(`^the registry is '(.+)'$`, tc.SetRegistry)

		ctx.When(`^the command is executed$`, tc.ExecuteCommand)

		ctx.Then(`^the FHIR cache contains packages:$`, tc.HasPackages)
	}
}

type TestCase struct {
	Timeout  time.Duration
	Command  []string
	Result   error
	Registry string
	T        *testing.T
	Verbose  bool

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

func (tc *TestCase) SetCommand(command string) error {
	tc.Command = strings.Split(command, " ")
	return nil
}

func (tc *TestCase) SetRegistry(registry string) error {
	tc.Registry = registry
	return nil
}

func (tc *TestCase) ExecuteCommand() error {
	if len(tc.Command) == 0 {
		return godog.ErrPending
	}
	command := tc.Command
	if tc.Timeout != 0 {
		command = append(command, "--timeout", tc.Timeout.String())
	}
	if tc.Registry != "" {
		command = append(command, "--registry", tc.Registry)
	}
	command = append(command, "--fhir-cache", tc.cache.Root())
	if tc.Verbose {
		tc.T.Logf("Running command: %v", strings.Join(command, " "))
	}
	out, err := exec.Command(command[0], command[1:]...).CombinedOutput()
	if tc.Verbose {
		tc.T.Logf(string(out))
	}
	return err
}

func (tc *TestCase) HasPackages(packages *godog.Table) error {
	for _, row := range packages.Rows[1:] {
		name := row.Cells[0].Value
		version := row.Cells[1].Value

		if _, err := tc.cache.Get("default", name, version); err != nil {
			return err
		}
	}
	return nil
}
