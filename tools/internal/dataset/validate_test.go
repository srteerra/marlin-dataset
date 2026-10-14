package dataset

import (
	"path/filepath"
	"strings"
	"testing"
)

var schemaPath = filepath.Join("..", "..", "..", "schema", "fish.schema.json")

// validateFixture validates a folder under testdata and fails on unexpected errors.
func validateFixture(t *testing.T, folder ...string) Result {
	t.Helper()
	parts := append([]string{"testdata"}, folder...)
	result, err := Validate(filepath.Join(parts...), schemaPath)
	if err != nil {
		t.Fatalf("Validate returned an error: %v", err)
	}
	return result
}

// hasIssue reports whether an issue matches the field (when given) and contains the reason.
func hasIssue(issues []Issue, field, reason string) bool {
	for _, issue := range issues {
		fieldMatches := field == "" || issue.Field == field
		if fieldMatches && strings.Contains(issue.Reason, reason) {
			return true
		}
	}
	return false
}

func TestValidateAcceptsValidDataset(t *testing.T) {
	result := validateFixture(t, "valid")
	if len(result.Issues) > 0 {
		t.Fatalf("expected no issues, got %v", result.Issues)
	}
	if len(result.Records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(result.Records))
	}
}

func TestValidateAcceptsRepositoryDataset(t *testing.T) {
	result, err := Validate(filepath.Join("..", "..", "..", "fish"), schemaPath)
	if err != nil {
		t.Fatalf("Validate returned an error: %v", err)
	}
	if len(result.Issues) > 0 {
		t.Fatalf("expected the repository dataset to be valid, got %v", result.Issues)
	}
}

func TestValidateRejectsInvalidFiles(t *testing.T) {
	cases := []struct {
		folder     string
		wantField  string
		wantReason string
	}{
		{folder: "extra-field", wantReason: "color"},
		{folder: "wrong-schema-reference", wantField: "$schema"},
		{folder: "image-not-pixabay", wantField: "images.0.url"},
		{folder: "missing-required", wantReason: "status"},
		{folder: "negative-size", wantField: "adult_size_cm"},
		{folder: "bad-json", wantReason: reasonInvalidJSON},
		{folder: "wrong-extension", wantReason: reasonWrongExtension},
		{folder: "wrong-folder", wantReason: reasonWrongLocation},
		{folder: "wrong-filename", wantReason: reasonWrongLocation},
		{folder: "id-mismatch", wantField: "id", wantReason: reasonIDMismatch},
		{folder: "min-above-max", wantField: "temperature_c", wantReason: reasonMinAboveMax},
		{folder: "ph-out-of-bounds", wantField: "ph", wantReason: reasonOutOfBounds},
		{folder: "temperature-out-of-bounds", wantField: "temperature_c", wantReason: reasonOutOfBounds},
		{folder: "reviewed-without-sources", wantField: "sources", wantReason: reasonReviewedWithoutSources},
		{folder: "duplicate-id", wantField: "id", wantReason: reasonDuplicateID},
		{folder: "unknown-reference", wantField: "incompatible_with.0.id", wantReason: reasonUnknownReference},
	}
	for _, testCase := range cases {
		t.Run(testCase.folder, func(t *testing.T) {
			result := validateFixture(t, "invalid", testCase.folder)
			if !hasIssue(result.Issues, testCase.wantField, testCase.wantReason) {
				t.Fatalf("expected issue (field %q, reason %q), got %v",
					testCase.wantField, testCase.wantReason, result.Issues)
			}
		})
	}
}

func TestValidateReportsAllIssuesTogether(t *testing.T) {
	result := validateFixture(t, "invalid", "multiple-errors")
	if !hasIssue(result.Issues, "temperature_c", reasonMinAboveMax) {
		t.Errorf("missing temperature_c issue in %v", result.Issues)
	}
	if !hasIssue(result.Issues, "ph", reasonOutOfBounds) {
		t.Errorf("missing ph issue in %v", result.Issues)
	}
}
