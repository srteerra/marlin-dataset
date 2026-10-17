// Command marlin-dataset validates the fish dataset and bundles it for releases.
package main

import (
	"fmt"
	"os"
)

// Process exit codes.
const (
	exitSuccess = 0
	exitFailure = 1
	exitUsage   = 2
)

const usage = `usage: go run . <command> [options] <fish-dir>

commands:
  validate   validate every file in the dataset
  bundle     validate, then write fish.json and fish.min.json for a release

examples:
  go run . validate ../fish
  go run . bundle --version v2026.10.0 ../fish`

// main runs the CLI and exits with its status code.
func main() {
	os.Exit(run(os.Args[1:]))
}

// run dispatches the subcommand and returns the process exit code.
func run(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, usage)
		return exitUsage
	}
	switch args[0] {
	case "validate":
		return runValidate(args[1:])
	case "bundle":
		return runBundle(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s\n", args[0], usage)
		return exitUsage
	}
}
