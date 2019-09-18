package movies

import (
	"testing"

	"gitlab.com/haath/goirate/pkg/torrents"
)

func TestGetURL(t *testing.T) {
	var table = []struct {
		in  string
		out string
	}{
		{"123", "https://www.imdb.com/title/tt0000123/"},
		{"-123", ""},
		{"123456789", ""},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {

			movie := MovieID{IMDbID: tt.in}
			url, err := movie.GetURL()

			if (url == nil && tt.out != "") || (err != nil && tt.out != "") {
				t.Errorf("got %v, want %v", url, tt.out)
			}
		})
	}
}

func TestFormattedDuration(t *testing.T) {
	table := []struct {
		in  int
		out string
	}{
		{143, "2h 23min"},
		{180, "3h"},
		{47, "47min"},
	}
	for _, tt := range table {
		t.Run(tt.out, func(t *testing.T) {
			m := Movie{Duration: tt.in}
			s := m.FormattedDuration()
			if s != tt.out {
				t.Errorf("got %q, want %q", s, tt.out)
			}
		})
	}
}

func TestGetTorrent(t *testing.T) {
	table := []struct {
		in Movie
	}{
		{Movie{MovieID: MovieID{Title: "third person"}}},
		//{Movie{MovieID: MovieID{Title: "Το τρίτο πρόσωπο", AltTitle: "Third Person"}}},
		{Movie{MovieID: MovieID{Title: "the loft"}}},
		{Movie{MovieID: MovieID{Title: "the loft", Year: 2014}}},
	}

	filters := torrents.SearchFilters{}

	for _, tt := range table {
		t.Run(tt.in.Title, func(t *testing.T) {

			scraper, err := torrents.FindScraper(tt.in.SearchQuery())

			if err != nil {
				t.Error(err)
			}

			tor, err := tt.in.GetTorrent(scraper, filters)

			if err != nil {
				t.Error(err)
			}

			if tor == nil {
				t.Errorf("No torrent found for: %v", tt.in.Title)
			}

		})
	}
}
