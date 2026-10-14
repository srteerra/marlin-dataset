package main

import (
	"flag"
	"fmt"
)

// runValidate implements `validate <fish-dir>`.
func runValidate(args []string) int {
	flags := flag.NewFlagSet("validate", flag.ContinueOnError)
	schemaPath := addSchemaFlag(flags)
	datasetDir, ok := parseDatasetArgs(flags, args)
	if !ok {
		return exitUsage
	}
	result, ok := validateDataset(datasetDir, *schemaPath)
	if !ok {
		return exitFailure
	}
	fmt.Printf("ok: %d valid species\n", len(result.Records))
	return exitSuccess
}
