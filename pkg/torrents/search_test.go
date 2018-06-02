package torrents

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"os"
	"testing"
)

func TestSearchTorrentList(t *testing.T) {
	table := []struct {
		in  SearchFilters
		out string
	}{}

	file, err := os.Open("../../samples/piratebay_movie.html")

	if err != nil {
		t.Error(err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		t.Error(err)
		return
	}

	u, _ := url.Parse("localhost")
	scraper := pirateBayScaper{u}

	torrents := scraper.ParseSearchPage(doc)

	for _, tt := range table {
		t.Run(tt.out, func(t *testing.T) {

			torrent, _ := SearchTorrentList(torrents, tt.in)

			if torrent.Title != tt.out {
				t.Errorf("got %v, want %v", torrent.Title, tt.out)
			}
		})
	}
}
