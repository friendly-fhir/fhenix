package model_test

import (
	"testing"

	"github.com/friendly-fhir/fhenix/internal/model"
)

func TestCardinality_IsRequired(t *testing.T) {
	testCases := []struct {
		name        string
		cardinality model.Cardinality
		want        bool
	}{
		{
			name:        "List",
			cardinality: model.Cardinality{Min: 0, Max: 2},
			want:        false,
		}, {
			name:        "UnboundedList",
			cardinality: model.Cardinality{Min: 0, Max: model.Unbound},
			want:        false,
		}, {
			name:        "Scalar",
			cardinality: model.Cardinality{Min: 0, Max: 1},
			want:        false,
		}, {
			name:        "Required",
			cardinality: model.Cardinality{Min: 1, Max: 1},
			want:        true,
		}, {
			name:        "Disabled",
			cardinality: model.Cardinality{Min: 0, Max: 0},
			want:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut := tc.cardinality

			got := sut.IsRequired()

			if got != tc.want {
				t.Errorf("Cardinality.IsOptional() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCardinality_IsDisabled(t *testing.T) {
	testCases := []struct {
		name        string
		cardinality model.Cardinality
		want        bool
	}{
		{
			name:        "List",
			cardinality: model.Cardinality{Min: 0, Max: 2},
			want:        false,
		}, {
			name:        "UnboundedList",
			cardinality: model.Cardinality{Min: 0, Max: model.Unbound},
			want:        false,
		}, {
			name:        "Scalar",
			cardinality: model.Cardinality{Min: 0, Max: 1},
			want:        false,
		}, {
			name:        "Required",
			cardinality: model.Cardinality{Min: 1, Max: 1},
			want:        false,
		}, {
			name:        "Disabled",
			cardinality: model.Cardinality{Min: 0, Max: 0},
			want:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut := tc.cardinality

			got := sut.IsDisabled()

			if got != tc.want {
				t.Errorf("Cardinality.IsOptional() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCardinality_IsScalar(t *testing.T) {
	testCases := []struct {
		name        string
		cardinality model.Cardinality
		want        bool
	}{
		{
			name:        "List",
			cardinality: model.Cardinality{Min: 0, Max: 2},
			want:        false,
		}, {
			name:        "UnboundedList",
			cardinality: model.Cardinality{Min: 0, Max: model.Unbound},
			want:        false,
		}, {
			name:        "Optional",
			cardinality: model.Cardinality{Min: 0, Max: 1},
			want:        true,
		}, {
			name:        "Required",
			cardinality: model.Cardinality{Min: 1, Max: 1},
			want:        true,
		}, {
			name:        "Disabled",
			cardinality: model.Cardinality{Min: 0, Max: 0},
			want:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut := tc.cardinality

			got := sut.IsScalar()

			if got != tc.want {
				t.Errorf("Cardinality.IsOptional() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCardinality_IsOptional(t *testing.T) {
	testCases := []struct {
		name        string
		cardinality model.Cardinality
		want        bool
	}{
		{
			name:        "List",
			cardinality: model.Cardinality{Min: 0, Max: 2},
			want:        false,
		}, {
			name:        "UnboundedList",
			cardinality: model.Cardinality{Min: 0, Max: model.Unbound},
			want:        false,
		}, {
			name:        "Optional",
			cardinality: model.Cardinality{Min: 0, Max: 1},
			want:        true,
		}, {
			name:        "Required",
			cardinality: model.Cardinality{Min: 1, Max: 1},
			want:        false,
		}, {
			name:        "Disabled",
			cardinality: model.Cardinality{Min: 0, Max: 0},
			want:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut := tc.cardinality

			got := sut.IsOptional()

			if got != tc.want {
				t.Errorf("Cardinality.IsOptional() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCardinality_IsList(t *testing.T) {
	testCases := []struct {
		name        string
		cardinality model.Cardinality
		want        bool
	}{
		{
			name:        "List",
			cardinality: model.Cardinality{Min: 0, Max: 2},
			want:        true,
		}, {
			name:        "UnboundedList",
			cardinality: model.Cardinality{Min: 0, Max: model.Unbound},
			want:        true,
		}, {
			name:        "Scalar",
			cardinality: model.Cardinality{Min: 0, Max: 1},
			want:        false,
		}, {
			name:        "Required",
			cardinality: model.Cardinality{Min: 1, Max: 1},
			want:        false,
		}, {
			name:        "Disabled",
			cardinality: model.Cardinality{Min: 0, Max: 0},
			want:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut := tc.cardinality

			got := sut.IsList()

			if got != tc.want {
				t.Errorf("Cardinality.IsList() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCardinality_IsUnboundedList(t *testing.T) {
	testCases := []struct {
		name        string
		cardinality model.Cardinality
		want        bool
	}{
		{
			name:        "Unbound",
			cardinality: model.Cardinality{Min: 0, Max: model.Unbound},
			want:        true,
		}, {
			name:        "Bounded",
			cardinality: model.Cardinality{Min: 0, Max: 1},
			want:        false,
		}, {
			name:        "Required",
			cardinality: model.Cardinality{Min: 1, Max: 1},
			want:        false,
		}, {
			name:        "Disabled",
			cardinality: model.Cardinality{Min: 0, Max: 0},
			want:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut := tc.cardinality

			got := sut.IsUnboundedList()

			if got != tc.want {
				t.Errorf("Cardinality.IsUnboundedList() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCardinality_String(t *testing.T) {
	testCases := []struct {
		name        string
		cardinality model.Cardinality
		want        string
	}{
		{
			name:        "Unbound",
			cardinality: model.Cardinality{Min: 0, Max: model.Unbound},
			want:        "0..*",
		}, {
			name:        "Bounded",
			cardinality: model.Cardinality{Min: 0, Max: 1},
			want:        "0..1",
		}, {
			name:        "Required",
			cardinality: model.Cardinality{Min: 1, Max: 1},
			want:        "1..1",
		}, {
			name:        "Disabled",
			cardinality: model.Cardinality{Min: 0, Max: 0},
			want:        "0..0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut := tc.cardinality

			got := sut.String()

			if got != tc.want {
				t.Errorf("Cardinality.String() = %q, want %q", got, tc.want)
			}
		})
	}
}
