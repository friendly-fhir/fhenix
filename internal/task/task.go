package task

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"

	"golang.org/x/sync/errgroup"
)

// Task is a function that can be executed by a Queue.
type Task interface {
	// Execute the task, returning its result.
	Execute(ctx context.Context) error
}

// Func is a convenience function that implements the [Task] interface.
type Func func(ctx context.Context) error

// Execute executes the task.
func (j Func) Execute(ctx context.Context) error {
	return j(ctx)
}

var _ Task = (*Func)(nil)

// Runner is a task queue that can execute multiple tasks concurrently.
type Runner struct {
	// This mutex protects the tasks channel from being closed multiple times, or
	// written to after it has been closed.
	m sync.Mutex

	tasks   chan Task
	workers int
	wg      *sync.WaitGroup
}

type ctxKey struct{}

// RunnerFromContext returns the Runner from the context.
//
// The context provided must have been created using [WithRunner], otherwise
// this function will return nil.
func RunnerFromContext(ctx context.Context) *Runner {
	value := ctx.Value(ctxKey{})
	if value == nil {
		return nil
	}
	return value.(*Runner)
}

// WithRunner returns a new context with the Runner associated to it.
//
// A context returned from this function will return the [Runner] when calling
// [RunnerFromContext].
func WithRunner(ctx context.Context, q *Runner) context.Context {
	return context.WithValue(ctx, ctxKey{}, q)
}

// NewRunner creates a new Runner with the specified number of workers.
//
// If workers is 0 or negative, this will set the number of workers to the
// number of CPUs available.
func NewRunner(workers int) *Runner {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	return &Runner{
		tasks:   make(chan Task),
		workers: workers,
		wg:      &sync.WaitGroup{},
	}
}

// Add adds a task to the runner.
//
// Tasks must be added before calling [Run], otherwise they may not be executed,
// depending on the state of the runner at the time of the call.
//
// Calling this function from within a [Task] is perfectly safe, and supported,
// for providing generator-jobs.
func (q *Runner) Add(j Task) {
	wg := q.wg
	wg.Add(1)
	ch := q.tasksChannel()
	go func() {
		ch <- Func(func(ctx context.Context) error {
			defer wg.Done()
			return j.Execute(ctx)
		})
	}()
}

func (q *Runner) tasksChannel() chan<- Task {
	q.m.Lock()
	defer q.m.Unlock()
	return q.tasks
}

// errDone is a sentinel error used to signal that the runner has completed.
var errDone = errors.New("done")

// Run executes all tasks in the runner, returning the number of tasks executed
// and an error if any of the tasks failed.
//
// If any tasks are added into the runner as part of their operation, this
// function will wait until all task sources have been exhausted, and all tasks
// are completed (unless cancelled manually by the caller).
//
// The function will block until all tasks have been executed.
func (q *Runner) Run(ctx context.Context) (int, error) {
	group, ctx := errgroup.WithContext(ctx)

	tasks := q.tasks

	// When our wait-group finishes, we close the tasks channel to signal that
	// we're done -- then reset it for a future execution.
	go func() {
		q.wg.Wait()

		q.m.Lock()
		close(tasks)
		q.tasks = make(chan Task)
		q.m.Unlock()
	}()

	var n atomic.Int32

	for i := 0; i < q.workers; i++ {
		group.Go(func() error {
			for {
				if err := q.runOne(ctx, tasks, &n); err != nil {
					if err == errDone {
						return nil
					}
					return err
				}
			}
		})
	}
	err := group.Wait()
	value := n.Load()
	return int(value), err
}

func (q *Runner) runOne(ctx context.Context, tasks <-chan Task, n *atomic.Int32) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case j, ok := <-tasks:
		if !ok {
			return errDone
		}

		// Imbue the current runner into the context so the task can add more
		// jobs into the queue.
		ctx = WithRunner(ctx, q)
		if err := j.Execute(ctx); err != nil {
			n.Add(1)
			return err
		}
		n.Add(1)
	}
	return nil
}
