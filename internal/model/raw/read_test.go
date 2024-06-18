package raw_test

import (
	"testing"

	"github.com/friendly-fhir/fhenix/internal/model/raw"
)

func TestReadStructureDefinition(t *testing.T) {
	file := "testdata/structure-definition.json"
	got, err := raw.ReadStructureDefinition(file)

	if err != nil {
		t.Errorf("ReadStructureDefinition(%s) = %v, want = %v", file, err, nil)
	}
	if got == nil {
		t.Fatalf("ReadStructureDefinition(%s) = %v, want non-nil", file, got)
	}
}

func TestReadValueSet(t *testing.T) {
	file := "testdata/value-set.json"
	got, err := raw.ReadValueSet(file)

	if err != nil {
		t.Errorf("ReadValueSet(%s) = %v, want = %v", file, err, nil)
	}
	if got == nil {
		t.Fatalf("ReadValueSet(%s) = %v, want non-nil", file, got)
	}
}

func TestReadCodeSystem(t *testing.T) {
	file := "testdata/code-system.json"
	got, err := raw.ReadCodeSystem(file)

	if err != nil {
		t.Errorf("ReadCodeSystem(%s) = %v, want = %v", file, err, nil)
	}
	if got == nil {
		t.Fatalf("ReadCodeSystem(%s) = %v, want non-nil", file, got)
	}
}
