package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/srteerra/marlin-dataset/tools/internal/dataset"
)

// defaultSchemaFile is resolved relative to the parent of the dataset folder.
var defaultSchemaFile = filepath.Join("schema", "fish.schema.json")

// addSchemaFlag registers the optional --schema override shared by every subcommand.
func addSchemaFlag(flags *flag.FlagSet) *string {
	return flags.String("schema", "", "path to the schema (default <fish-dir>/../schema/fish.schema.json)")
}

// parseDatasetArgs parses flags and returns the single dataset folder argument.
func parseDatasetArgs(flags *flag.FlagSet, args []string) (string, bool) {
	if err := flags.Parse(args); err != nil {
		return "", false
	}
	if flags.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "pass exactly one dataset folder, e.g. ../fish")
		return "", false
	}
	return flags.Arg(0), true
}

// resolveSchemaPath returns the override, or the schema next to the dataset folder.
func resolveSchemaPath(datasetDir, override string) string {
	if override != "" {
		return override
	}
	return filepath.Join(filepath.Dir(filepath.Clean(datasetDir)), defaultSchemaFile)
}

// validateDataset runs the validation, prints every issue and reports success.
func validateDataset(datasetDir, schemaOverride string) (dataset.Result, bool) {
	result, err := dataset.Validate(datasetDir, resolveSchemaPath(datasetDir, schemaOverride))
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
