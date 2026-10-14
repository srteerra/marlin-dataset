package dataset

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

// Incompatibility points to another species that should not share a tank.
type Incompatibility struct {
	ID     string `json:"id"`
	Reason string `json:"reason"`
}

// Species holds the fields the validation rules need.
type Species struct {
	ID               string            `json:"id"`
	ScientificName   string            `json:"scientific_name"`
	TemperatureC     *Range            `json:"temperature_c"`
	PH               *Range            `json:"ph"`
	HardnessDGH      *Range            `json:"hardness_dgh"`
	IncompatibleWith []Incompatibility `json:"incompatible_with"`
	Sources          []string          `json:"sources"`
	Status           string            `json:"status"`
}

// Record is a species loaded from disk.
type Record struct {
	Path    string
	Species Species
}
