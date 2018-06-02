package movies

import (
	"github.com/PuerkitoBio/goquery"
	"os"
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

func TestParseDuration(t *testing.T) {
	var table = []struct {
		in  string
		out int
	}{
		{"2h 23min", 143},
		{"3h", 180},
		{"47min", 47},
	}
	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {
			s := parseDuration(tt.in)
			if s != tt.out {
				t.Errorf("got %q, want %q", s, tt.out)
			}
		})
	}
}

func TestIMDbPage(t *testing.T) {

	expected := Movie{
		Title:     "Cast Away",
		Year:      2000,
		Duration:  143,
		Rating:    7.8,
		PosterURL: "https://m.media-amazon.com/images/M/MV5BN2Y5ZTU4YjctMDRmMC00MTg4LWE1M2MtMjk4MzVmOTE4YjkzXkEyXkFqcGdeQXVyNTc1NTQxODI@._V1_UX182_CR0,0,182,268_AL_.jpg",
	}

	file, err := os.Open("../../samples/imdb.html")

	if err != nil {
		t.Error(err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		t.Error(err)
		return
	}

	movie := ParseIMDbPage(doc)

	if movie.Title != expected.Title || movie.Year != expected.Year ||
		movie.Duration != expected.Duration || movie.Rating != expected.Rating ||
		movie.PosterURL != expected.PosterURL {
		t.Errorf("got: %v\nwant: %v\n", movie, expected)
	}
}
