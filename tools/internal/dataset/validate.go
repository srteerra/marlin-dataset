package dataset

import "fmt"

// Result holds the species that passed the schema and every issue found.
type Result struct {
	Records []Record
	Issues  []Issue
}

// Validate loads every species under datasetDir and checks it against the schema and dataset rules.
func Validate(datasetDir, schemaPath string) (Result, error) {
	schema, err := compileSchema(schemaPath)
	if err != nil {
		return Result{}, fmt.Errorf("compiling schema %s: %w", schemaPath, err)
	}
	files, err := listFiles(datasetDir)
	if err != nil {
		return Result{}, fmt.Errorf("reading %s: %w", datasetDir, err)
	}
	var result Result
	for _, path := range files {
		record, issues, err := loadRecord(schema, datasetDir, path)
		if err != nil {
			return Result{}, err
		}
		result.Issues = append(result.Issues, issues...)
		if record != nil {
			result.Records = append(result.Records, *record)
		}
	}
	result.Issues = append(result.Issues, duplicateIssues(result.Records)...)
	sortIssues(result.Issues)
	return result, nil
}
