package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/srteerra/marlin-dataset/tools/internal/dataset"
)

// addSchemaFlag registers the optional --schema override shared by every subcommand.
func addSchemaFlag(flags *flag.FlagSet) *string {
	return flags.String("schema", "", "path to the schema (default: the collection's schema, e.g. <dataset-dir>/../schema/plant.schema.json)")
}

// parseDatasetArgs parses flags and returns the single dataset folder argument.
func parseDatasetArgs(flags *flag.FlagSet, args []string) (string, bool) {
	if err := flags.Parse(args); err != nil {
		return "", false
	}
	if flags.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "pass exactly one dataset folder, e.g. ../fish or ../plants")
		return "", false
	}
	return flags.Arg(0), true
}

// resolveSchemaPath returns the override, or the collection's schema next to the dataset folder.
func resolveSchemaPath(datasetDir string, collection dataset.Collection, override string) string {
	if override != "" {
		return override
	}
	return datasetSibling(datasetDir, collection.SchemaFile)
}

// datasetSibling resolves a path relative to the parent of the dataset folder.
func datasetSibling(datasetDir, relativePath string) string {
	return filepath.Join(filepath.Dir(filepath.Clean(datasetDir)), relativePath)
}

// resolveCollection names the collection from the dataset folder, printing the error when unknown.
func resolveCollection(datasetDir string) (dataset.Collection, bool) {
	collection, err := dataset.CollectionFor(datasetDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return dataset.Collection{}, false
	}
	return collection, true
}

// validateDataset runs the validation, prints every issue and reports success.
func validateDataset(datasetDir string, collection dataset.Collection, schemaOverride string) (dataset.Result, bool) {
	result, err := dataset.Validate(datasetDir, resolveSchemaPath(datasetDir, collection, schemaOverride))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return result, false
	}
	for _, issue := range result.Issues {
		fmt.Fprintln(os.Stderr, issue)
	}
	if len(result.Issues) > 0 {
		fmt.Fprintf(os.Stderr, "\n%d errors found\n", len(result.Issues))
		return result, false
	}
	return result, true
}
