package dataset

import (
	"path/filepath"
	"strings"
	"testing"
)

var plantSchemaPath = filepath.Join("..", "..", "..", "schema", "plant.schema.json")

// validatePlantFixture validates a folder under testdata/plants and fails on unexpected errors.
func validatePlantFixture(t *testing.T, folder ...string) Result {
	t.Helper()
	parts := append([]string{"testdata", "plants"}, folder...)
	result, err := Validate(filepath.Join(parts...), plantSchemaPath)
	if err != nil {
		t.Fatalf("Validate returned an error: %v", err)
	}
	return result
}

func TestValidateAcceptsValidPlants(t *testing.T) {
	result := validatePlantFixture(t, "valid")
	if len(result.Issues) > 0 {
		t.Fatalf("expected no issues, got %v", result.Issues)
	}
	if len(result.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(result.Records))
	}
}

func TestValidateRejectsUnknownPlantZone(t *testing.T) {
	result := validatePlantFixture(t, "invalid", "unknown-zone")
	if !hasIssue(result.Issues, "zones.0", "") {
		t.Fatalf("expected an issue on zones.0, got %v", result.Issues)
	}
}

func TestValidateAcceptsRepositoryPlants(t *testing.T) {
	result, err := Validate(filepath.Join("..", "..", "..", "plants"), plantSchemaPath)
	if err != nil {
		t.Fatalf("Validate returned an error: %v", err)
	}
	if len(result.Issues) > 0 {
		t.Fatalf("expected the repository plants to be valid, got %v", result.Issues)
	}
}

func TestPlantBundleUsesPlantsKey(t *testing.T) {
	result := validatePlantFixture(t, "valid")
	bundle, err := BuildBundle(result.Records, collections[CollectionPlants], "v2026.11.0", fixtureGeneratedAt)
	if err != nil {
		t.Fatalf("BuildBundle returned an error: %v", err)
	}
	var output strings.Builder
	if err := bundle.Encode(&output, ""); err != nil {
		t.Fatalf("Encode returned an error: %v", err)
	}
	if !strings.Contains(output.String(), `"count":1,"plants":[`) {
		t.Errorf("expected the records under the plants key, got %s", output.String())
	}
}

func TestCollectionFor(t *testing.T) {
	for _, testCase := range []struct{ dir, want string }{{"../fish", "fish"}, {"../plants/", "plants"}} {
		collection, err := CollectionFor(testCase.dir)
		if err != nil || collection.Name != testCase.want {
			t.Errorf("CollectionFor(%q) = %q, %v; want %q", testCase.dir, collection.Name, err, testCase.want)
		}
	}
	if _, err := CollectionFor("../corals"); err == nil {
		t.Error("CollectionFor accepted an unknown folder")
	}
}
