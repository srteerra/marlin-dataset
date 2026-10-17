package dataset

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strings"
	"time"
)

// schemaKey only helps editors, so it is left out of the bundle.
const schemaKey = "$schema"

// versionPattern matches release tags such as v2026.10.0.
var versionPattern = regexp.MustCompile(`^v\d{4}\.\d{2}\.\d+$`)

// Bundle is the release file holding every species of one dataset version.
type Bundle struct {
	Version     string                       `json:"version"`
	GeneratedAt time.Time                    `json:"generated_at"`
	Count       int                          `json:"count"`
	Fish        []map[string]json.RawMessage `json:"fish"`
}

// BuildBundle collects every record sorted by id, without the $schema key.
func BuildBundle(records []Record, version string, generatedAt time.Time) (Bundle, error) {
	if !versionPattern.MatchString(version) {
		return Bundle{}, fmt.Errorf("invalid version %q, expected vYYYY.MM.N", version)
	}
	sorted := slices.Clone(records)
	slices.SortFunc(sorted, func(left, right Record) int {
		return strings.Compare(left.Species.ID, right.Species.ID)
	})
	fish := make([]map[string]json.RawMessage, 0, len(sorted))
	for _, record := range sorted {
		fields, err := bundleFields(record)
		if err != nil {
			return Bundle{}, err
		}
		fish = append(fish, fields)
	}
	return Bundle{Version: version, GeneratedAt: generatedAt.UTC(), Count: len(fish), Fish: fish}, nil
}

// bundleFields decodes a record keeping each value exactly as written (e.g. pH 6.0 stays 6.0).
func bundleFields(record Record) (map[string]json.RawMessage, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(record.Raw, &fields); err != nil {
		return nil, fmt.Errorf("%s: %w", record.Path, err)
	}
	delete(fields, schemaKey)
	return fields, nil
}

// Encode writes the bundle as JSON; an empty indent produces the minified form.
func (bundle Bundle) Encode(writer io.Writer, indent string) error {
	encoder := json.NewEncoder(writer)
	// Source URLs often contain &, which would otherwise be written as &.
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", indent)
	return encoder.Encode(bundle)
}
