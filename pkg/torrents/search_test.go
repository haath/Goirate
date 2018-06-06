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
	}{
		{SearchFilters{}, "Cast Away (2000) 1080p BrRip x264 - 1.10GB - YIFY"},
		{SearchFilters{MaxSize: "1 GB"}, "Cast Away (2000) 720p BrRip x264 - 950MB - YIFY"},
		{SearchFilters{MinSize: "3 GB"}, "Cast.Away.2000.1080p.BluRay.x264.AC3-ETRG"},
		{SearchFilters{MaxQuality: Medium}, "Cast Away (2000) 720p BrRip x264 - 950MB - YIFY"},
		{SearchFilters{MinSeeders: 500}, ""},
	}

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

			torrent, err := SearchTorrentList(torrents, tt.in)

			if tt.out != "" && (torrent == nil || err != nil) {
				t.Error(err)
				return
			}

			if tt.out != "" && torrent.Title != tt.out {
				t.Errorf("\ngot: %v\nwant: %v\n", torrent.Title, tt.out)
			}
		})
	}
}

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

			s := normalizeQuery(tt.in)

			if tt.out != s {
				t.Errorf("\ngot: %v\nwant: %v\n", s, tt.out)
			}
		})
	}
}
