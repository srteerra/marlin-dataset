package dataset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

const jsonExtension = ".json"

// listFiles returns every regular file under dir in lexical order.
func listFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && !isHidden(entry.Name()) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// isHidden skips OS and editor artifacts such as .DS_Store or .gitkeep.
func isHidden(name string) bool {
	return strings.HasPrefix(name, ".")
}

// loadRecord parses one file and runs the per-file checks.
// A nil record means the file could not be checked beyond the returned issues.
func loadRecord(schema *jsonschema.Schema, datasetDir, path string) (*Record, []Issue, error) {
	if filepath.Ext(path) != jsonExtension {
		return nil, []Issue{{File: path, Reason: reasonWrongExtension}}, nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	instance, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
	if err != nil {
		reason := fmt.Sprintf("%s: %v", reasonInvalidJSON, err)
		return nil, []Issue{{File: path, Reason: reason}}, nil
	}
	issues, err := schemaIssues(schema, path, instance)
	if err != nil || len(issues) > 0 {
		return nil, issues, err
	}
	record := Record{Path: path}
	if err := json.Unmarshal(raw, &record.Species); err != nil {
		return nil, nil, fmt.Errorf("%s: %w", path, err)
	}
	return &record, recordIssues(datasetDir, record), nil
}
