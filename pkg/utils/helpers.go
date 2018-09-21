package utils

import "strings"

// OptionalBoolean defines a boolean constant that can also be undefined.
type OptionalBoolean string

const (
	// Default represents an undefined value.
	Default OptionalBoolean = ""
	// True represents a positive value.
	True OptionalBoolean = "true"
	// False represents a negative value.
	False OptionalBoolean = "false"
)

// WatchlistAction defines actions to be taken upon discovering a new torrent,
// along with any parameters regarding said action.
type WatchlistAction struct {
	SendEmail OptionalBoolean `toml:"email" json:"email"`
	Emails    []string        `toml:"notify" json:"notify"`
	Download  OptionalBoolean `toml:"download" json:"download"`
}

// OverridenBy returns true if one of this or the other action is true, or if the other action is true.
func (opt OptionalBoolean) OverridenBy(other OptionalBoolean) bool {

	return other == True || (other != False && opt == True)
}

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
	query = strings.TrimSpace(query)

	return query
}
