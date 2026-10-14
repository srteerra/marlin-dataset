package dataset

import (
	"fmt"
	"math"
	"path/filepath"
)

// Physical limits for range checks.
const (
	minTemperatureC = 0
	maxTemperatureC = 40
	minPH           = 0
	maxPH           = 14
	minHardnessDGH  = 0
)

// rangeRule binds a range field to its allowed bounds.
type rangeRule struct {
	field string
	value *Range
	lower float64
	upper float64
}

// recordIssues runs every per-file rule that the schema cannot express.
func recordIssues(datasetDir string, record Record) []Issue {
	var issues []Issue
	issues = append(issues, locationIssues(datasetDir, record)...)
	issues = append(issues, rangeIssues(record)...)
	issues = append(issues, statusIssues(record)...)
	return issues
}

// locationIssues checks that the id matches the name and the file sits in fish/<genus>/<id>.json.
func locationIssues(datasetDir string, record Record) []Issue {
	var issues []Issue
	species := record.Species
	expectedID := Slugify(species.ScientificName)
	if species.ID != expectedID {
		reason := fmt.Sprintf("%s %q", reasonIDMismatch, expectedID)
		issues = append(issues, Issue{File: record.Path, Field: "id", Reason: reason})
	}
	expectedPath := filepath.Join(Genus(species.ID), species.ID+jsonExtension)
	relativePath, err := filepath.Rel(datasetDir, record.Path)
	if err != nil || relativePath != expectedPath {
		reason := fmt.Sprintf("%s %s", reasonWrongLocation, expectedPath)
		issues = append(issues, Issue{File: record.Path, Reason: reason})
	}
	return issues
}

// rangeIssues checks min <= max and physical limits for every range field.
func rangeIssues(record Record) []Issue {
	species := record.Species
	rules := []rangeRule{
		{field: "temperature_c", value: species.TemperatureC, lower: minTemperatureC, upper: maxTemperatureC},
		{field: "ph", value: species.PH, lower: minPH, upper: maxPH},
		{field: "hardness_dgh", value: species.HardnessDGH, lower: minHardnessDGH, upper: math.Inf(1)},
	}
	var issues []Issue
	for _, rule := range rules {
		issues = append(issues, checkRange(record.Path, rule)...)
	}
	return issues
}

// checkRange validates a single range field; null ranges are allowed.
func checkRange(file string, rule rangeRule) []Issue {
	if rule.value == nil {
		return nil
	}
	var issues []Issue
	if rule.value.Min > rule.value.Max {
		reason := fmt.Sprintf("%s (%g > %g)", reasonMinAboveMax, rule.value.Min, rule.value.Max)
		issues = append(issues, Issue{File: file, Field: rule.field, Reason: reason})
	}
	if !rule.contains(rule.value.Min) || !rule.contains(rule.value.Max) {
		reason := fmt.Sprintf("%s, %s", reasonOutOfBounds, describeBounds(rule.lower, rule.upper))
		issues = append(issues, Issue{File: file, Field: rule.field, Reason: reason})
	}
	return issues
}

// contains reports whether value lies within the rule's physical limits.
func (rule rangeRule) contains(value float64) bool {
	return value >= rule.lower && value <= rule.upper
}

// describeBounds renders the allowed interval for error messages.
func describeBounds(lower, upper float64) string {
	if math.IsInf(upper, 1) {
		return fmt.Sprintf("must be >= %g", lower)
	}
	return fmt.Sprintf("must be between %g and %g", lower, upper)
}

// statusIssues enforces that reviewed species cite at least one source.
func statusIssues(record Record) []Issue {
	species := record.Species
	if species.Status != StatusReviewed || len(species.Sources) > 0 {
		return nil
	}
	return []Issue{{File: record.Path, Field: "sources", Reason: reasonReviewedWithoutSources}}
}
