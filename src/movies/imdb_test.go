package movies

import (
	"testing"
)

func TestFormatIMDbID(t *testing.T) {
	var table = []struct {
		in  string
		out string
	}{
		{"123", "tt0000123"},
		{"-123", ""},
		{"123456789", "tt123456789"},
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

func TestIsIMDbID(t *testing.T) {
	var table = []struct {
		in  string
		out bool
	}{
		{"123", true},
		{"tt1234567", true},
		{"-123", false},
		{"123456789", true},
		{"Hail ceasar", false},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {
			s := IsIMDbID(tt.in)
			if s != tt.out {
				t.Errorf("got %v, want %v", s, tt.out)
			}
		})
	}

}

func TestIsIMDbURL(t *testing.T) {
	var table = []struct {
		in  string
		out bool
	}{
		{"123", false},
		{"https://www.imdb.com/title/tt0368226/", true},
		{"https://www.imdb.com/title/tt03/", false},
		{"https://www.imdb.com/title/tt6155194/?ref_=nm_knf_i3", true},
		{"https://www.imdb.com/title/tt0848226", true},
		{"https://m.imdb.com/title/tt0848226", true},
		{"Avengers: Age of Ultron", false},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {
			s := IsIMDbURL(tt.in)
			if s != tt.out {

				id, err := ExtractIMDbID(tt.in)

				t.Errorf("got %v, want %v. %v %v", s, tt.out, id, err)
			}
		})
	}

}

func TestExtractIMDbID(t *testing.T) {
	var table = []struct {
		in  string
		out string
	}{
		{"https://www.imdb.com/title/tt0848228/?ref_=fn_al_tt_1/", "tt0848228"},
		{"https://www.imdb.com/title/tt0848226", "tt0848226"},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {
			s, err := ExtractIMDbID(tt.in)

			if err != nil {
				t.Error(err)
			}

			if s != tt.out {
				t.Errorf("got %q, want %q", s, tt.out)
			}
		})
	}
}
