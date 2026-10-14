package dataset

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var schemaMessagePrinter = message.NewPrinter(language.English)

// compileSchema loads and compiles the JSON schema at path.
func compileSchema(path string) (*jsonschema.Schema, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return jsonschema.NewCompiler().Compile(absolutePath)
}

// schemaIssues validates a decoded JSON document and returns one issue per failing rule.
func schemaIssues(schema *jsonschema.Schema, file string, instance any) ([]Issue, error) {
	err := schema.Validate(instance)
	if err == nil {
		return nil, nil
	}
	var validationError *jsonschema.ValidationError
	if !errors.As(err, &validationError) {
		return nil, fmt.Errorf("%s: %w", file, err)
	}
	return leafIssues(file, validationError), nil
}

// leafIssues flattens the error tree; only leaves carry a concrete reason.
func leafIssues(file string, validationError *jsonschema.ValidationError) []Issue {
	if len(validationError.Causes) == 0 {
		return []Issue{{
			File:   file,
			Field:  strings.Join(validationError.InstanceLocation, "."),
			Reason: validationError.ErrorKind.LocalizedString(schemaMessagePrinter),
		}}
	}
	var issues []Issue
	for _, cause := range validationError.Causes {
		issues = append(issues, leafIssues(file, cause)...)
	}
	return issues
}
