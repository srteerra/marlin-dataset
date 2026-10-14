package dataset

import "testing"

func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"Paracheirodon innesi":    "paracheirodon-innesi",
		"  Betta   splendens ":    "betta-splendens",
		"Corydoras sp. 'C123'":    "corydoras-sp-c123",
		"Apistogramma agassizii":  "apistogramma-agassizii",
		"Pseudotropheus démasoni": "pseudotropheus-demasoni",
	}
	for input, want := range cases {
		if got := Slugify(input); got != want {
			t.Errorf("Slugify(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestGenus(t *testing.T) {
	if got := Genus("paracheirodon-innesi"); got != "paracheirodon" {
		t.Errorf("Genus = %q, want %q", got, "paracheirodon")
	}
}
