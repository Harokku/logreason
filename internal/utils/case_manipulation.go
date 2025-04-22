package utilities

import (
	"regexp"
	"strings"
	"unicode"
)

func ToCamelCase(s string) string {
	// Split the string by spaces and special characters
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	var result string
	for i, word := range words {
		// Convert word to title case (first letter uppercase, rest lowercase)
		word = strings.ToLower(word)
		if len(word) > 0 {
			if i == 0 {
				result += word
			} else {
				result += strings.Title(word)
			}
		}
	}
	return result
}

// ToPascalCase converts a space-separated string with potential special characters
// into a single PascalCase word. It removes non-alphabetic characters
// and capitalizes the first letter of each resulting segment.
func ToPascalCase(input string) string {
	// Regular expression to find sequences of non-alphabetic characters (delimiters)
	// This will split the string by spaces, punctuation, numbers, etc.
	delimiters := regexp.MustCompile(`[^a-zA-Z]+`)

	// Split the input string into words based on the delimiters.
	// The -1 argument means split all occurrences.
	words := delimiters.Split(input, -1)

	var resultParts []string

	// Process each word segment
	for _, word := range words {
		// Skip empty strings that can result from multiple delimiters together
		// or leading/trailing delimiters.
		if len(word) == 0 {
			continue
		}

		// Convert the word to runes for correct Unicode handling
		runes := []rune(word)

		// Capitalize the first rune (letter)
		runes[0] = unicode.ToUpper(runes[0])

		// Append the capitalized word to our results
		resultParts = append(resultParts, string(runes))
	}

	// Join the processed parts together without any separator
	return strings.Join(resultParts, "")
}
