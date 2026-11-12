package main

import (
	"flag"
	"fmt"
)

// runValidate implements `validate <dataset-dir>`, e.g. ../fish or ../plants.
func runValidate(args []string) int {
	flags := flag.NewFlagSet("validate", flag.ContinueOnError)
	schemaPath := addSchemaFlag(flags)
	datasetDir, ok := parseDatasetArgs(flags, args)
	if !ok {
		return exitUsage
	}
	collection, ok := resolveCollection(datasetDir)
	if !ok {
		return exitUsage
	}
	result, ok := validateDataset(datasetDir, collection, *schemaPath)
	if !ok {
		return exitFailure
	}
	fmt.Printf("ok: %d valid %s species\n", len(result.Records), collection.Name)
	return exitSuccess
}
