package torrents

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestSearchTorrentList(t *testing.T) {
	table := []struct {
		in  SearchFilters
		out string
		num int
	}{
		{SearchFilters{}, "Cast Away (2000) 1080p BrRip x264 - 1.10GB - YIFY", 4},
		{SearchFilters{MaxSize: "1 GB"}, "Cast Away (2000) 720p BrRip x264 - 950MB - YIFY", 3},
		{SearchFilters{MinSize: "3 GB"}, "Cast.Away.2000.1080p.BluRay.x264.AC3-ETRG", 3},
		{SearchFilters{MaxQuality: Medium}, "Cast Away (2000) 720p BrRip x264 - 950MB - YIFY", 4},
		{SearchFilters{MaxQuality: Low, MinQuality: Low}, "Cast.Away.2000.480p.DVDRip.XviD-ViEW", 4},
		{SearchFilters{MinSeeders: 500}, "", 0},
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
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {

			torrent, err := SearchTorrentList(torrents, tt.in)

			if tt.out != "" && (torrent == nil || err != nil) {
				t.Error(err)
				return
			}

			if tt.out != "" && torrent.Title != tt.out {
				t.Errorf("\ngot: %v\nwant: %v\n", torrent.Title, tt.out)
			}

			multi, err := SearchVideoTorrentList(torrents, tt.in)

			if tt.out != "" && (torrent == nil || err != nil) {
				t.Error(err)
				return
			}

			if tt.out != "" && len(multi) != tt.num {
				t.Errorf("error fetching multiple qualities: %v", multi)
			}

			best, err := PickVideoTorrent(torrents, tt.in)

			if tt.out != "" && err != nil {
				t.Error(err)
				return
			}

			if tt.out != "" && best == nil {
				t.Error("torrent not found")
				return
			}

			if tt.out != "" && best.Title != tt.out {
				t.Errorf("\ngot: %v\nwant: %v\n", best.Title, tt.out)
			}
		})
	}
}

func TestSizeKB(t *testing.T) {
	table := []struct {
		in  string
		out int64
	}{
		{"64 KB", 64},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {

			f := SearchFilters{MinSize: tt.in, MaxSize: tt.in}

			s1, err := f.MinSizeKB()

			if err != nil {
				t.Error(err)
			}

			if s1 != tt.out {
				t.Errorf("got %v, want %v", s1, tt.out)
			}

			s2, err := f.MaxSizeKB()

			if err != nil {
				t.Error(err)
			}

			if s2 != tt.out {
				t.Errorf("got %v, want %v", s2, tt.out)
			}
		})
	}
}
