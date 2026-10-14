package dataset

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// Slugify converts a scientific name into a kebab-case id without accents.
func Slugify(name string) string {
	// NFD splits accented letters into base letter + combining mark, so marks can be dropped.
	decomposed := norm.NFD.String(strings.ToLower(name))
	var builder strings.Builder
	pendingDash := false
	for _, character := range decomposed {
		switch {
		case unicode.Is(unicode.Mn, character):
			continue
		case isSlugCharacter(character):
			if pendingDash && builder.Len() > 0 {
				builder.WriteByte('-')
			}
			pendingDash = false
			builder.WriteRune(character)
		default:
			pendingDash = true
		}
	}
	return builder.String()
}

// isSlugCharacter reports whether the rune can appear in a slug as is.
func isSlugCharacter(character rune) bool {
	return (character >= 'a' && character <= 'z') || (character >= '0' && character <= '9')
}

// Genus returns the first part of a species id, which names its folder.
func Genus(id string) string {
	genus, _, _ := strings.Cut(id, "-")
	return genus
}
