package snek

import (
	"sync"

	"github.com/spf13/cobra"
)

type commandProperty[T any] struct {
	defaultValue T
	lookup       *sync.Map
}

func newCommandProperty[T any](defaultValue T) commandProperty[T] {
	return commandProperty[T]{
		defaultValue: defaultValue,
		lookup:       &sync.Map{},
	}
}

func (cp *commandProperty[T]) Get(cmd *cobra.Command) T {
	val, ok := cp.lookup.Load(cmd)
	if ok {
		return val.(T)
	}
	return cp.defaultValue
}

func (cp *commandProperty[T]) Set(cmd *cobra.Command, value T) {
	cp.lookup.Store(cmd, value)
}

func (cp *commandProperty[T]) SetDefault(value T) {
	cp.defaultValue = value
}

var (
	commandFlags = newCommandProperty([]*FlagSet{})
)
