package set

import (
	"cmp"
	"slices"
)

// Set is an implementation of the missing "set" datastructure for Go.
type Set[T cmp.Ordered] map[T]struct{}

// New creates a new, but empty, set.
func New[T cmp.Ordered]() Set[T] {
	return make(Set[T])
}

// Of creates a new set with the given values.
func Of[T cmp.Ordered](values ...T) Set[T] {
	if len(values) == 0 {
		return nil
	}
	set := New[T]()
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

// FromSlice creates a new set from the given slice.
func FromSlice[T cmp.Ordered](values []T) Set[T] {
	return Of(values...)
}

func FromKeys[T cmp.Ordered, V any](m map[T]V) Set[T] {
	if len(m) == 0 {
		return nil
	}
	set := make(Set[T])
	for key := range m {
		set[key] = struct{}{}
	}
	return set
}

// FromValues creates a new set from the values of the given map.
func FromValues[T, V cmp.Ordered](m map[T]V) Set[V] {
	if len(m) == 0 {
		return nil
	}
	set := make(Set[V])
	for _, value := range m {
		set[value] = struct{}{}
	}
	return set
}

// Intersection returns the intersection of the given sets.
func Intersection[T cmp.Ordered](sets ...Set[T]) Set[T] {
	if len(sets) == 0 {
		return nil
	}
	intersection := make(Set[T])
	for value := range sets[0] {
		found := true
		for _, set := range sets[1:] {
			if !set.Contains(value) {
				found = false
				break
			}
		}
		if found {
			intersection[value] = struct{}{}
		}
	}
	return intersection
}

// Difference returns the difference of the given sets.
func Difference[T cmp.Ordered](lhs, rhs Set[T]) Set[T] {
	difference := make(Set[T])
	for value := range lhs {
		if !rhs.Contains(value) {
			difference[value] = struct{}{}
		}
	}
	for value := range rhs {
		if !lhs.Contains(value) {
			difference[value] = struct{}{}
		}
	}
	return difference
}

// Union returns the union of the given sets.
func Union[T cmp.Ordered](sets ...Set[T]) Set[T] {
	union := make(Set[T])
	for _, set := range sets {
		for value := range set {
			union[value] = struct{}{}
		}
	}
	return union
}

// Clone returns a shallow copy of the set.
func (s Set[T]) Clone() Set[T] {
	if s == nil {
		return nil
	}
	other := make(Set[T], len(s))
	for value := range s {
		other[value] = struct{}{}
	}
	return other
}

// Len returns the number of elements in the set.
func (s Set[T]) Len() int {
	return len(s)
}

// IsEmpty returns true if the set is empty.
func (s Set[T]) IsEmpty() bool {
	return len(s) == 0
}

// Add adds a value to the set.
func (s *Set[T]) Add(value T) {
	if *s == nil {
		*s = make(Set[T])
	}
	(*s)[value] = struct{}{}
}

// Remove removes a value from the set.
func (s Set[T]) Remove(value T) {
	delete(s, value)
}

// Contains returns true if the set contains the given value.
func (s Set[T]) Contains(value T) bool {
	if s == nil {
		return false
	}
	_, ok := s[value]
	return ok
}

// Equal returns true if the set is equal to the other set.
func (s Set[T]) Equal(other Set[T]) bool {
	if len(s) != len(other) {
		return false
	}
	for value := range s {
		if !other.Contains(value) {
			return false
		}
	}
	return true
}

// Slice returns the values of the set as a slice, without any sorting.
func (s Set[T]) Slice() []T {
	values := make([]T, 0, len(s))
	for value := range s {
		values = append(values, value)
	}
	return values
}

// SortedSlice returns the values of the set as a slice, sorted by the default
// comparison function.
func (s Set[T]) SortedSlice() []T {
	return s.SortedSliceFunc(cmp.Compare[T])
}

// SortedSliceFunc returns the values of the set as a slice, sorted using the
// given comparison function.
func (s Set[T]) SortedSliceFunc(cmp func(T, T) int) []T {
	slice := s.Slice()
	slices.SortFunc(slice, cmp)
	return slice
}
