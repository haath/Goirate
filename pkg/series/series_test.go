package series

import (
	"fmt"
	"testing"

	"git.gmantaos.com/haath/Goirate/pkg/torrents"
)

func TestSearchQuery(t *testing.T) {
	table := []struct {
		in  string
		ep  Episode
		out string
	}{
		{"Scraping The Barrel", Episode{4, 5}, "scraping the barrel s04e05"},
		{"Scraping The Barrel", Episode{4, 0}, "scraping the barrel season 4"},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {

			series := Series{Title: tt.in}

			out := series.SearchQuery(tt.ep)

			if out != tt.out {
				t.Errorf("got %v, want %v", out, tt.out)
			}

		})
	}
}

func TestGetTorrent(t *testing.T) {
	table := []struct {
		in Series
		ep Episode
	}{
		{Series{Title: "Game of Thrones"}, Episode{1, 1}},
		{Series{Title: "Game of Thrones"}, Episode{2, 0}},
	}

	filters := torrents.SearchFilters{}

	for _, tt := range table {
		t.Run(tt.in.Title, func(t *testing.T) {

			scraper, err := torrents.FindScraper(tt.in.SearchQuery(tt.ep))

			if err != nil {
				t.Error(err)
			}

			tor, err := tt.in.GetTorrent(scraper, filters, tt.ep)

			if err != nil {
				t.Error(err)
			}

			if tor == nil {
				t.Errorf("No torrent found for: %v", tt.in.Title)
			}

		})
	}
}

func TestNextEpisode(t *testing.T) {

	table := []struct {
		in   int
		last Episode
		next Episode
	}{
		{261690, Episode{6, 10}, Episode{6, 11}},
		{121361, Episode{1, 0}, Episode{1, 1}},
		{255316, Episode{5, 24}, Episode{6, 1}},
	}

	tkn := login(t)

	for _, tt := range table {
		t.Run(fmt.Sprint(tt.in), func(t *testing.T) {

			ser := Series{ID: tt.in, LastEpisode: tt.last}

			next, err := ser.NextEpisode(&tkn)
			if err != nil {
				t.Error(err)
			}

			if next != tt.next {
				t.Errorf("got %v, want %v", next, tt.next)
			}

		})
	}
}
