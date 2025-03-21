package npmtest_test

import (
	"testing"

	"github.com/friendly-fhir/fhenix/pkg/registryv2/internal/npm/npmtest"
	"github.com/google/go-cmp/cmp"
)

func TestSimplePackageTarball_HasSameContentsAsFS(t *testing.T) {
	want := npmtest.FSContents(t, npmtest.SimplePackageFS())
	r := npmtest.SimplePackageTarball()
	defer r.Close()

	got := npmtest.TarContents(t, r)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected contents: %v", diff)
	}
}
