package task_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/friendly-fhir/fhenix/internal/task"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRunner(t *testing.T) {
	runner := task.NewRunner(10)

	runner.Add(task.Func(func(ctx context.Context) error {
		runner := task.RunnerFromContext(ctx)
		runner.Add(task.Func(func(ctx context.Context) error {
			return nil
		}))
		runner.Add(task.Func(func(ctx context.Context) error {
			return nil
		}))
		return nil
	}))
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	count, err := runner.Run(ctx)

	if got, want := err, (error)(nil); got != want {
		t.Errorf("Runner.Run() = %v, want %v", got, want)
	}
	if got, want := count, 3; got != want {
		t.Errorf("Runner.Run() = %d, want %d", got, want)
	}
}

func TestRunner_CancelledContext(t *testing.T) {
	runner := task.NewRunner(10)

	runner.Add(task.Func(func(ctx context.Context) error {
		runner := task.RunnerFromContext(ctx)
		runner.Add(task.Func(func(ctx context.Context) error {
			return nil
		}))
		runner.Add(task.Func(func(ctx context.Context) error {
			return nil
		}))
		return nil
	}))
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()

	_, err := runner.Run(ctx)

	if got, want := err, context.Canceled; got != want {
		t.Errorf("Runner.Run() = %v, want %v", got, want)
	}
}

func TestRunner_ErroredJob(t *testing.T) {
	testErr := errors.New("test error")
	runner := task.NewRunner(10)

	runner.Add(task.Func(func(ctx context.Context) error {
		runner := task.RunnerFromContext(ctx)
		runner.Add(task.Func(func(ctx context.Context) error {
			return testErr
		}))
		runner.Add(task.Func(func(ctx context.Context) error {
			return nil
		}))
		return nil
	}))
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	_, err := runner.Run(ctx)

	if got, want := err, testErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
		t.Errorf("Runner.Run() = %v, want %v", got, want)
	}
}
