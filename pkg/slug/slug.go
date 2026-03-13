// Package slug
package slug

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func Generate(input string) string {
	// Normalize unicode characters (é → e + accent mark)
	s := norm.NFD.String(input)

	var b strings.Builder
	prevDash := false

	for _, r := range s {

		// remove accent marks
		if unicode.Is(unicode.Mn, r) {
			continue
		}

		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(unicode.ToLower(r))
			prevDash = false

		case unicode.IsSpace(r) || r == '-' || r == '_':
			if !prevDash {
				b.WriteRune('-')
				prevDash = true
			}
		}
	}

	result := b.String()

	result = strings.Trim(result, "-")

	return result
}
