package dataset

import "fmt"

// duplicateIssues reports every record whose id was already used by an earlier file.
func duplicateIssues(records []Record) []Issue {
	firstPathByID := make(map[string]string, len(records))
	var issues []Issue
	for _, record := range records {
		id := record.Species.ID
		firstPath, exists := firstPathByID[id]
		if !exists {
			firstPathByID[id] = record.Path
			continue
		}
		reason := fmt.Sprintf("%s %s", reasonDuplicateID, firstPath)
		issues = append(issues, Issue{File: record.Path, Field: "id", Reason: reason})
	}
	return issues
}

// referenceIssues reports incompatible_with entries that point to unknown species.
func referenceIssues(records []Record) []Issue {
	knownIDs := make(map[string]bool, len(records))
	for _, record := range records {
		knownIDs[record.Species.ID] = true
	}
	var issues []Issue
	for _, record := range records {
		for index, incompatibility := range record.Species.IncompatibleWith {
			if knownIDs[incompatibility.ID] {
				continue
			}
			field := fmt.Sprintf("incompatible_with.%d.id", index)
			reason := fmt.Sprintf("%s %q", reasonUnknownReference, incompatibility.ID)
			issues = append(issues, Issue{File: record.Path, Field: field, Reason: reason})
		}
	}
	return issues
}
