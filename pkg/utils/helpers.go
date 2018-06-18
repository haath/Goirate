package utils

import "strings"

// NormalizeQuery will appropriate replace special characters in a title as to normalize it for better comparisons.
func NormalizeQuery(query string) string {

	replaces := []struct {
		old string
		new string
	}{
		{"-", " "},
		{"'", " "},
		{".", " "},
		{"_", " "},
		{":", " "},
		{"!", " "},
		{"(", " "},
		{")", " "},
	}

	query = strings.ToLower(query)
	query = strings.TrimSpace(query)

	// Remove special characters
	for _, rep := range replaces {
		query = strings.Replace(query, rep.old, rep.new, -1)
	}

	// Remove extra spaces
	query = strings.Join(strings.Fields(query), " ")

	return query
}
