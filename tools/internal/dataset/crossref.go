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
