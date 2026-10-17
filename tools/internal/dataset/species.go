package dataset

import "encoding/json"

// Status values allowed in the dataset.
const (
	StatusDraft    = "draft"
	StatusReviewed = "reviewed"
)

// Range is a { min, max } interval.
type Range struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// Species holds the fields the validation rules need.
type Species struct {
	ID             string   `json:"id"`
	ScientificName string   `json:"scientific_name"`
	TemperatureC   *Range   `json:"temperature_c"`
	PH             *Range   `json:"ph"`
	HardnessDGH    *Range   `json:"hardness_dgh"`
	Sources        []string `json:"sources"`
	Status         string   `json:"status"`
}

// Record is a species loaded from disk.
// Raw keeps the file as written, so the bundle carries every field, not only the ones in Species.
type Record struct {
	Path    string
	Species Species
	Raw     json.RawMessage
}
