package set_test

import (
	stdcmp "cmp"
	"testing"

	"github.com/friendly-fhir/fhenix/internal/set"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestOf(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		elements []int
		want     set.Set[int]
	}{
		{
			name:     "nil input",
			elements: nil,
			want:     nil,
		}, {
			name:     "empty input",
			elements: []int{},
			want:     nil,
		}, {
			name:     "single element",
			elements: []int{1},
			want:     set.Set[int]{1: {}},
		}, {
			name:     "duplicate elements",
			elements: []int{1, 1, 1},
			want:     set.Set[int]{1: {}},
		}, {
			name:     "multiple elements",
			elements: []int{1, 2, 3, 2, 1},
			want:     set.Set[int]{1: {}, 2: {}, 3: {}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := set.Of(tc.elements...)

			if !got.Equal(tc.want) {
				t.Errorf("Of(...) = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFromSlice(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    []int
		want set.Set[int]
	}{
		{
			name: "nil input",
			s:    nil,
			want: nil,
		}, {
			name: "empty input",
			s:    []int{},
			want: nil,
		}, {
			name: "single element",
			s:    []int{1},
			want: set.Set[int]{1: {}},
		}, {
			name: "duplicate elements",
			s:    []int{1, 1, 1},
			want: set.Set[int]{1: {}},
		}, {
			name: "multiple elements",
			s:    []int{1, 2, 3, 2, 1},
			want: set.Set[int]{1: {}, 2: {}, 3: {}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := set.FromSlice(tc.s)

			if !got.Equal(tc.want) {
				t.Errorf("FromSlice(...) = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFromKeys(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		m    map[int]string
		want set.Set[int]
	}{
		{
			name: "nil input",
			m:    nil,
			want: nil,
		}, {
			name: "empty input",
			m:    map[int]string{},
			want: nil,
		}, {
			name: "single element",
			m:    map[int]string{1: ""},
			want: set.Set[int]{1: {}},
		}, {
			name: "multiple elements",
			m:    map[int]string{1: "", 2: "", 3: ""},
			want: set.Set[int]{1: {}, 2: {}, 3: {}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := set.FromKeys(tc.m)

			if !got.Equal(tc.want) {
				t.Errorf("FromKeys(...) = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFromValues(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		m    map[string]int
		want set.Set[int]
	}{
		{
			name: "nil input",
			m:    nil,
			want: nil,
		}, {
			name: "empty input",
			m:    map[string]int{},
			want: nil,
		}, {
			name: "single element",
			m:    map[string]int{"a": 1},
			want: set.Set[int]{1: {}},
		}, {
			name: "multiple elements",
			m:    map[string]int{"a": 1, "b": 2, "c": 3},
			want: set.Set[int]{1: {}, 2: {}, 3: {}},
		}, {
			name: "duplicate elements",
			m:    map[string]int{"a": 1, "b": 1, "c": 1},
			want: set.Set[int]{1: {}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := set.FromValues(tc.m)

			if !got.Equal(tc.want) {
				t.Errorf("FromValues(...) = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIntersection(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		sets []set.Set[int]
		want set.Set[int]
	}{
		{
			name: "nil input",
			sets: nil,
			want: nil,
		}, {
			name: "empty input",
			sets: []set.Set[int]{},
			want: nil,
		}, {
			name: "single set",
			sets: []set.Set[int]{
				{1: {}, 2: {}, 3: {}},
			},
			want: set.Set[int]{1: {}, 2: {}, 3: {}},
		}, {
			name: "two sets with overlap",
			sets: []set.Set[int]{
				{1: {}, 2: {}, 3: {}},
				{2: {}, 3: {}, 4: {}},
			},
			want: set.Set[int]{2: {}, 3: {}},
		}, {
			name: "three sets with overlap",
			sets: []set.Set[int]{
				{1: {}, 2: {}, 3: {}},
				{2: {}, 3: {}, 4: {}},
				{3: {}, 4: {}, 5: {}},
			},
			want: set.Set[int]{3: {}},
		}, {
			name: "three sets with no overlap",
			sets: []set.Set[int]{
				{1: {}, 2: {}, 3: {}},
				{4: {}, 5: {}, 6: {}},
				{7: {}, 8: {}, 9: {}},
			},
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := set.Intersection(tc.sets...)

			if !got.Equal(tc.want) {
				t.Errorf("Intersection(...) = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDifference(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		lhs  set.Set[int]
		rhs  set.Set[int]
		want set.Set[int]
	}{
		{
			name: "nil input",
			lhs:  nil,
			rhs:  nil,
			want: nil,
		}, {
			name: "empty input",
			lhs:  set.Set[int]{},
			rhs:  set.Set[int]{},
			want: nil,
		}, {
			name: "single set",
			lhs:  set.Set[int]{1: {}, 2: {}, 3: {}},
			rhs:  nil,
			want: set.Set[int]{1: {}, 2: {}, 3: {}},
		}, {
			name: "two sets with overlap",
			lhs:  set.Set[int]{1: {}, 2: {}, 3: {}},
			rhs:  set.Set[int]{2: {}, 3: {}, 4: {}},
			want: set.Set[int]{1: {}, 4: {}},
		}, {
			name: "two sets with no overlap",
			lhs:  set.Set[int]{1: {}, 2: {}, 3: {}},
			rhs:  set.Set[int]{4: {}, 5: {}, 6: {}},
			want: set.Set[int]{1: {}, 2: {}, 3: {}, 4: {}, 5: {}, 6: {}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := set.Difference(tc.lhs, tc.rhs)

			if !got.Equal(tc.want) {
				t.Errorf("Difference(...) = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestUnion(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		sets []set.Set[int]
		want set.Set[int]
	}{
		{
			name: "nil input",
			sets: nil,
			want: nil,
		}, {
			name: "empty input",
			sets: []set.Set[int]{},
			want: nil,
		}, {
			name: "single set",
			sets: []set.Set[int]{
				{1: {}, 2: {}, 3: {}},
			},
			want: set.Set[int]{1: {}, 2: {}, 3: {}},
		}, {
			name: "two sets with overlap",
			sets: []set.Set[int]{
				{1: {}, 2: {}, 3: {}},
				{2: {}, 3: {}, 4: {}},
			},
			want: set.Set[int]{1: {}, 2: {}, 3: {}, 4: {}},
		}, {
			name: "three sets with overlap",
			sets: []set.Set[int]{
				{1: {}, 2: {}, 3: {}},
				{2: {}, 3: {}, 4: {}},
				{3: {}, 4: {}, 5: {}},
			},
			want: set.Set[int]{1: {}, 2: {}, 3: {}, 4: {}, 5: {}},
		}, {
			name: "three sets with no overlap",
			sets: []set.Set[int]{
				{1: {}, 2: {}, 3: {}},
				{4: {}, 5: {}, 6: {}},
				{7: {}, 8: {}, 9: {}},
			},
			want: set.Set[int]{1: {}, 2: {}, 3: {}, 4: {}, 5: {}, 6: {}, 7: {}, 8: {}, 9: {}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := set.Union(tc.sets...)

			if !got.Equal(tc.want) {
				t.Errorf("Union(...) = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSet_Clone(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    set.Set[int]
	}{
		{
			name: "nil input",
			s:    nil,
		}, {
			name: "empty input",
			s:    set.Set[int]{},
		}, {
			name: "single element",
			s:    set.Set[int]{1: {}},
		}, {
			name: "multiple elements",
			s:    set.Set[int]{1: {}, 2: {}, 3: {}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.s.Clone()

			if !got.Equal(tc.s) {
				t.Errorf("Set.Clone() = %v, want %v", got, tc.s)
			}
		})
	}
}

func TestSet_Len(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    set.Set[int]
		want int
	}{
		{
			name: "nil input",
			s:    nil,
			want: 0,
		}, {
			name: "empty input",
			s:    set.Set[int]{},
			want: 0,
		}, {
			name: "single element",
			s:    set.Set[int]{1: {}},
			want: 1,
		}, {
			name: "multiple elements",
			s:    set.Set[int]{1: {}, 2: {}, 3: {}},
			want: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.s.Len()

			if got != tc.want {
				t.Errorf("Set.Len() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSet_IsEmpty(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    set.Set[int]
		want bool
	}{
		{
			name: "nil input",
			s:    nil,
			want: true,
		}, {
			name: "empty input",
			s:    set.Set[int]{},
			want: true,
		}, {
			name: "single element",
			s:    set.Set[int]{1: {}},
			want: false,
		}, {
			name: "multiple elements",
			s:    set.Set[int]{1: {}, 2: {}, 3: {}},
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.s.IsEmpty()

			if got != tc.want {
				t.Errorf("Set.IsEmpty() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSet_Contains(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		s     set.Set[int]
		value int
		want  bool
	}{
		{
			name:  "set does not contains",
			s:     nil,
			value: 1,
			want:  false,
		}, {
			name:  "set does contain",
			s:     set.Set[int]{1: {}},
			value: 1,
			want:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.s.Contains(tc.value)

			if got != tc.want {
				t.Errorf("Set.Contains(...) = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSet_Equal(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		s     set.Set[int]
		other set.Set[int]
		want  bool
	}{
		{
			name:  "nil inputs equate equal",
			s:     nil,
			other: nil,
			want:  true,
		}, {
			name:  "empty inputs equate equal with nil",
			s:     set.Set[int]{},
			other: nil,
			want:  true,
		}, {
			name:  "sets are equivalent",
			s:     set.Set[int]{1: {}, 2: {}, 3: {}},
			other: set.Set[int]{1: {}, 2: {}, 3: {}},
			want:  true,
		}, {
			name:  "sets have no overlap",
			s:     set.Set[int]{1: {}, 2: {}, 3: {}},
			other: set.Set[int]{4: {}, 5: {}, 6: {}},
			want:  false,
		}, {
			name:  "sets have some overlap",
			s:     set.Set[int]{1: {}, 2: {}, 3: {}},
			other: set.Set[int]{2: {}, 3: {}, 4: {}},
			want:  false,
		}, {
			name:  "sets have different lengths",
			s:     set.Set[int]{1: {}, 2: {}, 3: {}},
			other: set.Set[int]{1: {}, 2: {}},
			want:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.s.Equal(tc.other)

			if got != tc.want {
				t.Errorf("Set.Equal(...) = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSet_Add(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		s     set.Set[int]
		value int
		want  set.Set[int]
	}{
		{
			name:  "nil set creates new set",
			s:     nil,
			value: 1,
			want:  set.Set[int]{1: {}},
		}, {
			name:  "empty set",
			s:     set.Set[int]{},
			value: 1,
			want:  set.Set[int]{1: {}},
		}, {
			name:  "existing set",
			s:     set.Set[int]{1: {}},
			value: 2,
			want:  set.Set[int]{1: {}, 2: {}},
		}, {
			name:  "duplicate value",
			s:     set.Set[int]{1: {}},
			value: 1,
			want:  set.Set[int]{1: {}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.s.Add(tc.value)

			if !tc.s.Equal(tc.want) {
				t.Errorf("Set.Add(...) = %v, want %v", tc.s, tc.want)
			}
		})
	}
}

func TestSet_Remove(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		s     set.Set[int]
		value int
		want  set.Set[int]
	}{
		{
			name:  "nil set",
			s:     nil,
			value: 1,
			want:  nil,
		}, {
			name:  "empty set",
			s:     set.Set[int]{},
			value: 1,
			want:  set.Set[int]{},
		}, {
			name:  "existing set",
			s:     set.Set[int]{1: {}},
			value: 1,
			want:  set.Set[int]{},
		}, {
			name:  "non-existing value",
			s:     set.Set[int]{1: {}},
			value: 2,
			want:  set.Set[int]{1: {}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.s.Remove(tc.value)

			if !tc.s.Equal(tc.want) {
				t.Errorf("Set.Remove(...) = %v, want %v", tc.s, tc.want)
			}
		})
	}
}

func TestSet_Slice(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    set.Set[int]
		want []int
	}{
		{
			name: "nil input",
			s:    nil,
			want: nil,
		}, {
			name: "empty input",
			s:    set.Set[int]{},
			want: nil,
		}, {
			name: "single element",
			s:    set.Set[int]{1: {}},
			want: []int{1},
		}, {
			name: "multiple elements",
			s:    set.Set[int]{1: {}, 2: {}, 3: {}},
			want: []int{1, 2, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.s.Slice()

			if !cmp.Equal(got, tc.want, cmpopts.SortSlices(stdcmp.Less[int]), cmpopts.EquateEmpty()) {
				t.Errorf("Set.Slice() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSet_SortedSlice(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    set.Set[int]
		want []int
	}{
		{
			name: "nil input",
			s:    nil,
			want: nil,
		}, {
			name: "empty input",
			s:    set.Set[int]{},
			want: nil,
		}, {
			name: "single element",
			s:    set.Set[int]{1: {}},
			want: []int{1},
		}, {
			name: "multiple elements",
			s:    set.Set[int]{3: {}, 1: {}, 2: {}},
			want: []int{1, 2, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.s.SortedSlice()

			if !cmp.Equal(got, tc.want, cmpopts.EquateEmpty()) {
				t.Errorf("Set.SortedSlice() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSet_SortedSliceFunc(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    set.Set[int]
		want []int
	}{
		{
			name: "nil input",
			s:    nil,
			want: nil,
		}, {
			name: "empty input",
			s:    set.Set[int]{},
			want: nil,
		}, {
			name: "single element",
			s:    set.Set[int]{1: {}},
			want: []int{1},
		}, {
			name: "multiple elements",
			s:    set.Set[int]{3: {}, 1: {}, 2: {}},
			want: []int{1, 2, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.s.SortedSliceFunc(stdcmp.Compare[int])

			if !cmp.Equal(got, tc.want, cmpopts.EquateEmpty()) {
				t.Errorf("Set.SortedSlice() = %v, want %v", got, tc.want)
			}
		})
	}
}
