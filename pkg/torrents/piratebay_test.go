package torrents

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var urlTests = []struct {
	in  string
	out string
}{
	{"test_url_123", "test_url_123"},
	{"https://localhost", "https://localhost"},
	{"https://localhost/", "https://localhost/"},
	{"https://localhost:8080/", "https://localhost:8080/"},
}

var videoQualityTests = []struct {
	in  string
	out VideoQuality
}{
	{"The.Expanse.S02E03.PROPER.HDTV.x264-KILLERS[ettv]", Default},
	{"The.Expanse.S02E03.1080p.AMZN.WEBRip.DD5.1.HEVC.x265.sharpysword", High},
	{"The.Expanse.S02E03.720p.HDTV.x264-AVS", Medium},
	{"The.Expanse.S02E03.WEB-DL.XviD-FUM[ettv]", Default},
	{"The.Expanse.S02E03.480p.164mb.hdtv.x264-][ Static ][ 09- mp4", Low},
	{"Avengers Infinity War 2018 NEW PROPER 720p HD-CAM X264 HQ-CPG", Medium},
}

func TestNewScraper(t *testing.T) {
	for _, tt := range urlTests {
		t.Run(tt.in, func(t *testing.T) {
			s := NewScraper(tt.in)
			if s.URL() != tt.out {
				t.Errorf("got %q, want %q", s.URL(), tt.out)
			}
		})
	}
}

func TestFindScraper(t *testing.T) {
	_, err := FindScraper("ubuntu")

	if err != nil {
		t.Error(err)
	}
}

func TestSearchURL(t *testing.T) {

	searchTests := []struct {
		in  string
		out string
	}{
		{"test", "https://pirateproxy.sh/search/test"},
		{"one two", "https://pirateproxy.sh/search/one+two"},
		{"one'two", "https://pirateproxy.sh/search/one+two"},
		{"one!", "https://pirateproxy.sh/search/one"},
	}

	for _, tt := range searchTests {
		t.Run(tt.in, func(t *testing.T) {
			s := NewScraper("https://pirateproxy.sh/")
			searchURL := s.SearchURL(tt.in)
			if searchURL != tt.out {
				t.Errorf("got %q, want %q", searchURL, tt.out)
			}
		})
	}
}

func TestParseSearchPage(t *testing.T) {

	expected := Torrent{
		TorrentURL:       "/torrent/21412907/Avengers_Infinity_War_2018_NEW_PROPER_720p_HD-CAM_X264_HQ-CPG",
		Seeders:          7437,
		Leeches:          2003,
		Title:            "Avengers Infinity War 2018 NEW PROPER 720p HD-CAM X264 HQ-CPG",
		Size:             3843996,
		Magnet:           "magnet:?xt=urn:btih:bee75372b98077bfd4de8ef03eb33e9289be5cd8&dn=Avengers+Infinity+War+2018+NEW+PROPER+720p+HD-CAM+X264+HQ-CPG&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Fzer0day.ch%3A1337&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969",
		VideoQuality:     Medium,
		VerifiedUploader: true,
		Uploader:         "makintos13",
	}

	file, err := os.Open("../../samples/piratebay_search.html")

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

	if len(torrents) != 30 {
		t.Errorf("got %v, want 30", len(torrents))
	}

	tr := torrents[0]

	if tr.TorrentURL != expected.TorrentURL || tr.Title != expected.Title ||
		tr.Size != expected.Size || tr.Leeches != expected.Leeches || tr.Seeders != expected.Seeders ||
		tr.Magnet != expected.Magnet || tr.VideoQuality != expected.VideoQuality ||
		tr.VerifiedUploader != expected.VerifiedUploader || tr.Uploader != expected.Uploader {

		t.Errorf("\ngot: %v\nwant: %v\n", tr, expected)
	}
}

func TestExtractSize(t *testing.T) {

	sizeTests := []struct {
		in  string
		out int64
	}{
		{"Uploaded 04-29 04:41, Size 3.58 GiB, ULed by makintos13", 3843996},
		{"Uploaded 02-27 2014, Size 58.35 MiB, ULed by gnv65", 61184},
		{"Uploaded 10-12 2008, Size 740.35 KiB, ULed by my_name_is_bob", 758},
	}

	for _, tt := range sizeTests {
		t.Run(tt.in, func(t *testing.T) {
			s := extractSize(tt.in)
			if s != tt.out {
				t.Errorf("got %v, want %v", s, tt.out)
			}
		})
	}
}

func TestExtractUploadTime(t *testing.T) {

	yday := time.Now().AddDate(0, 0, -1)
	minago := time.Now().Add(time.Duration(-50) * time.Minute)

	table := []struct {
		in     string
		year   int
		month  time.Month
		day    int
		hour   int
		minute int
	}{
		{"Uploaded 04-29 04:41, Size 3.58 GiB, ULed by makintos13", time.Now().Year(), time.April, 29, 4, 41},
		{"Uploaded 02-27 2014, Size 58.35 MiB, ULed by gnv65", 2014, time.February, 27, 0, 0},
		{"Uploaded 10-12 2008, Size 740.35 KiB, ULed by my_name_is_bob", 2008, time.October, 12, 0, 0},
		{" Uploaded 04-27 20:41, Size 788.25 MiB, ULed by shmasti", time.Now().Year(), time.April, 27, 20, 41},
		{"Uploaded Today 08:05, Size 1.62 GiB, ULedbyAnonymous", time.Now().Year(), time.Now().Month(), time.Now().Day(), 8, 5},
		{"Uploaded Y-day 08:05, Size 1.62 GiB, ULedbyAnonymous", yday.Year(), yday.Month(), yday.Day(), 8, 5},
		{"Uploaded 50 mins ago, Size 1.62 GiB, ULedbyAnonymous", minago.Year(), minago.Month(), minago.Day(), minago.Hour(), minago.Minute()},
	}
	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {
			s := extractUploadTime(tt.in)
			if s.Year() != tt.year || s.Month() != tt.month || s.Day() != tt.day || s.Hour() != tt.hour || s.Minute() != tt.minute {
				t.Errorf("got %v, want %v", s, tt)
			}
		})
	}
}

func TestExtractVideoQuality(t *testing.T) {
	for _, tt := range videoQualityTests {
		t.Run(tt.in, func(t *testing.T) {
			s := extractVideoQuality(tt.in)
			if s != tt.out {
				t.Errorf("got %v, want %v", s, tt)
			}
		})
	}
}

func TestSearch(t *testing.T) {

	scraper, err := FindScraper("Windows 10")

	if err != nil {
		t.Errorf("Error finding scraper: %v\n", err)
	}

	torrents, err := (*scraper).Search("Windows 10")

	if err != nil {
		t.Error(err)
	}

	if len(torrents) == 0 {
		t.Errorf("Search yielded 0 torrents")
	}
}
