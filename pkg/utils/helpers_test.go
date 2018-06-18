package utils

import (
	"testing"
)

func TestNormalizeQuery(t *testing.T) {
	table := []struct {
		in  string
		out string
	}{
		{"Spider-Man: Homecoming", "Spider Man Homecoming"},
		{"The Hitchhiker's Guide to the Galaxy", "The Hitchhiker s Guide to the Galaxy"},
		{"American Dad!", "American Dad"},
	}

	for _, tt := range table {
		t.Run(tt.out, func(t *testing.T) {

			s := NormalizeQuery(tt.in)

			if tt.out != s {
				t.Errorf("\ngot: %v\nwant: %v\n", s, tt.out)
			}
		})
	}
}
