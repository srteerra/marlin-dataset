package dataset

import (
	"strings"
	"testing"
	"time"
)

var fixtureGeneratedAt = time.Date(2026, 10, 14, 12, 0, 0, 0, time.UTC)

func TestBuildBundle(t *testing.T) {
	result := validateFixture(t, "valid")
	bundle, err := BuildBundle(result.Records, "v2026.10.0", fixtureGeneratedAt)
	if err != nil {
		t.Fatalf("BuildBundle returned an error: %v", err)
	}
	if bundle.Count != 2 || len(bundle.Fish) != 2 {
		t.Fatalf("expected 2 fish, got count %d and %d items", bundle.Count, len(bundle.Fish))
	}
	for index, wantID := range []string{`"testus-alpha"`, `"testus-beta"`} {
		fish := bundle.Fish[index]
		if got := string(fish["id"]); got != wantID {
			t.Errorf("fish %d: id = %s, want %s", index, got, wantID)
		}
		if _, found := fish[schemaKey]; found {
			t.Errorf("fish %d: %s should be removed", index, schemaKey)
		}
	}
}

func TestBuildBundleRejectsInvalidVersion(t *testing.T) {
	for _, version := range []string{"", "2026.10.0", "v2026.1.0", "v2026.10"} {
		if _, err := BuildBundle(nil, version, fixtureGeneratedAt); err == nil {
			t.Errorf("BuildBundle accepted version %q", version)
		}
	}
}

func TestBundleEncodeKeepsValuesAsWritten(t *testing.T) {
	result := validateFixture(t, "valid")
	bundle, err := BuildBundle(result.Records, "v2026.10.0", fixtureGeneratedAt)
	if err != nil {
		t.Fatalf("BuildBundle returned an error: %v", err)
	}
	var output strings.Builder
	if err := bundle.Encode(&output, ""); err != nil {
		t.Fatalf("Encode returned an error: %v", err)
	}
	if !strings.Contains(output.String(), `"ph":{"min":6.0,"max":7.5}`) {
		t.Errorf("expected pH written as in the source file, got %s", output.String())
	}
}
