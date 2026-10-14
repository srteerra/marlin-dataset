package dataset

import (
	"fmt"
	"sort"
)

// Reasons shared by the rules and the tests, so messages stay consistent.
const (
	reasonWrongExtension         = "invalid extension, use .json"
	reasonInvalidJSON            = "invalid JSON"
	reasonWrongLocation          = "wrong location, expected"
	reasonIDMismatch             = "does not match scientific_name, expected"
	reasonMinAboveMax            = "min is greater than max"
	reasonOutOfBounds            = "out of bounds"
	reasonReviewedWithoutSources = "status reviewed requires at least one source"
	reasonDuplicateID            = "duplicate id, already used in"
	reasonUnknownReference       = "unknown species"
)

// Issue describes a single validation problem in a dataset file.
type Issue struct {
	File   string
	Field  string
	Reason string
}

// String formats the issue as "file: field: reason".
func (issue Issue) String() string {
	if issue.Field == "" {
		return fmt.Sprintf("%s: %s", issue.File, issue.Reason)
	}
	return fmt.Sprintf("%s: %s: %s", issue.File, issue.Field, issue.Reason)
}

// sortIssues orders issues by file and field so output is stable.
func sortIssues(issues []Issue) {
	sort.SliceStable(issues, func(left, right int) bool {
		if issues[left].File != issues[right].File {
			return issues[left].File < issues[right].File
		}
		return issues[left].Field < issues[right].Field
	})
}
