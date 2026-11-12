package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/srteerra/marlin-dataset/tools/internal/dataset"
)

// Release assets are <collection>.json and <collection>.min.json; the minified one is what consumers download.
const (
	defaultOutputDir     = "dist"
	bundleExtension      = ".json"
	minBundleExtension   = ".min.json"
	bundleIndent         = "  "
	minifiedBundleIndent = ""
)

// bundleFile pairs an output file name with its indentation.
type bundleFile struct {
	name   string
	indent string
}

// bundleFilesFor names the release assets of a collection, e.g. plants.json and plants.min.json.
func bundleFilesFor(collection string) []bundleFile {
	return []bundleFile{
		{name: collection + bundleExtension, indent: bundleIndent},
		{name: collection + minBundleExtension, indent: minifiedBundleIndent},
	}
}

// runBundle implements `bundle --version vYYYY.MM.N <dataset-dir>`.
func runBundle(args []string) int {
	flags := flag.NewFlagSet("bundle", flag.ContinueOnError)
	schemaPath := addSchemaFlag(flags)
	version := flags.String("version", "", "release version, e.g. v2026.10.0 (required)")
	outputDir := flags.String("out", "", "output folder (default <dataset-dir>/../dist)")
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
	bundle, err := dataset.BuildBundle(result.Records, collection, *version, time.Now())
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return exitFailure
	}
	directory := resolveOutputDir(datasetDir, *outputDir)
	if err := writeBundleFiles(bundle, directory); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return exitFailure
	}
	fmt.Printf("ok: %s with %d %s species written to %s\n", bundle.Version, bundle.Count, collection.Name, directory)
	return exitSuccess
}

// resolveOutputDir returns the override, or dist/ next to the dataset folder.
func resolveOutputDir(datasetDir, override string) string {
	if override != "" {
		return override
	}
	return datasetSibling(datasetDir, defaultOutputDir)
}

// writeBundleFiles writes every bundle file into directory, creating it if needed.
func writeBundleFiles(bundle dataset.Bundle, directory string) error {
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return err
	}
	for _, file := range bundleFilesFor(bundle.Collection) {
		if err := writeBundleFile(bundle, filepath.Join(directory, file.name), file.indent); err != nil {
			return err
		}
	}
	return nil
}

// writeBundleFile encodes the bundle into a single file.
func writeBundleFile(bundle dataset.Bundle, path, indent string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := bundle.Encode(file, indent); err != nil {
		file.Close()
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return file.Close()
}
