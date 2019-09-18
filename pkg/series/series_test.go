package series

import (
	"fmt"
	"testing"

	"gitlab.com/haath/goirate/pkg/torrents"
)

func TestSearchQuery(t *testing.T) {
	table := []struct {
		in  string
		ep  Episode
		out string
	}{
		{"Scraping The Barrel", Episode{Season: 4, Episode: 5}, "scraping the barrel s04e05"},
		{"Scraping The Barrel", Episode{Season: 4, Episode: 0}, "scraping the barrel season 4"},
		{"House of Cards (US)", Episode{Season: 3, Episode: 3}, "house of cards s03e03"},
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
		{Series{Title: "Game of Thrones"}, Episode{Season: 1, Episode: 1}},
		{Series{Title: "Game of Thrones"}, Episode{Season: 2, Episode: 0}},
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
		{261690, Episode{Season: 6, Episode: 10}, Episode{Season: 6, Episode: 11}},
		{121361, Episode{Season: 1, Episode: 0}, Episode{Season: 1, Episode: 1, Title: "Winter Is Coming"}},
		{255316, Episode{Season: 5, Episode: 24}, Episode{Season: 6, Episode: 1, Title: "An Infinite Capacity for Taking Pains"}},
	}

	tkn := login(t)

	for _, tt := range table {
		t.Run(fmt.Sprint(tt.in), func(t *testing.T) {

			ser := Series{ID: tt.in, LastEpisode: tt.last}

			next, err := ser.NextEpisode(&tkn)
			if err != nil {
				t.Error(err)
			}

			if next.String() != tt.next.String() {
				t.Errorf("got %v, want %v", next.String(), tt.next.String())
			}

			if next.Title != tt.next.Title {
				t.Errorf("got %v, want %v", next.Title, tt.next.Title)
			}

		})
	}
}
