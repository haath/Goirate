package torrents

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func OpenTestSample(sampleFile string) ([]Torrent, error) {

	file, err := os.Open(sampleFile)

	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, err
	}

	scraper := NewScraper("localhost")

	torrentList := scraper.ParseSearchPage(doc)

	return torrentList, nil
}

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

	torrents, err := OpenTestSample("../../test_samples/piratebay_movie.html")

	if err != nil {
		t.Error(err)
	}

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
				t.Errorf("error fetching multiple qualities. got %v, want %v", multi, tt.num)
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

func TestUploaderOk(t *testing.T) {

	var tor Torrent
	tor.Uploader = "someDude"

	table := []struct {
		in  SearchFilters
		out bool
	}{
		{SearchFilters{}, true},
		{SearchFilters{
			Uploaders: UploaderFilters{Blacklist: []string{"someDude"}},
		},
			false,
		},
		{SearchFilters{
			Uploaders: UploaderFilters{Whitelist: []string{"otherDude"}},
		},
			false,
		},
		{SearchFilters{
			Uploaders: UploaderFilters{Blacklist: []string{"someDude"}, Whitelist: []string{"someDude"}},
		},
			false,
		},
		{SearchFilters{
			Uploaders: UploaderFilters{Blacklist: []string{"otherDude"}, Whitelist: []string{"someDude"}},
		},
			true,
		},
	}

	for _, tt := range table {
		t.Run(fmt.Sprintf("%v %v", tt.in.Uploaders.Whitelist, tt.in.Uploaders.Blacklist), func(t *testing.T) {

			ok := tt.in.UploaderOk(tor.Uploader)

			if ok != tt.out {
				t.Errorf("got %v, want %v", ok, tt.out)
			}

		})
	}
}

func TestFilterTorrentList(t *testing.T) {

	torrentList, err := OpenTestSample("../../test_samples/piratebay_search.html")

	if err != nil {
		t.Error(err)
	}

	var table = []struct {
		in    func() SearchFilters
		count uint
		out   int
	}{
		{func() SearchFilters { return SearchFilters{} }, 0, 30},
		{func() SearchFilters {
			cmd := SearchFilters{}
			cmd.VerifiedUploader = true
			return cmd
		}, 0, 21},
		{func() SearchFilters {
			cmd := SearchFilters{}
			return cmd
		}, 1, 1},
	}

	for _, tt := range table {
		t.Run(strconv.Itoa(tt.out), func(t *testing.T) {
			filt := tt.in()

			s := filt.FilterTorrentsCount(torrentList, tt.count)
			if len(s) != tt.out {
				t.Errorf("\ngot: %v\nwant: %v", len(s), tt.out)
			}
		})
	}
}
