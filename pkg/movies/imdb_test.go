package movies

import (
	"testing"
)

func TestFormatIMDbID(t *testing.T) {
	var table = []struct {
		in  string
		out string
	}{
		{"123", "0000123"},
		{"-123", ""},
		{"123456789", ""},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {
			s, _ := FormatIMDbID(tt.in)
			if s != tt.out {
				t.Errorf("got %q, want %q", s, tt.out)
			}
		})
	}
}
