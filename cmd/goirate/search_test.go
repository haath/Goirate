package main

import (
	"encoding/json"
	"git.gmantaos.com/haath/Goirate/pkg/torrents"
	"github.com/PuerkitoBio/goquery"
	"os"
	"strconv"
	"testing"
)

func TestSearchExecute(t *testing.T) {

	var cmd SearchCommand
	Options.JSON = true

	output := CaptureCommand(func() { cmd.Execute([]string{"avengers"}) })

	var mirrors []torrents.Mirror
	json.Unmarshal([]byte(output), &mirrors)

	Options.JSON = false
}

func TestGetTorrentsTable(t *testing.T) {
	var table = []struct {
		in  []torrents.Torrent
		out string
	}{
		{[]torrents.Torrent{}, " Title  Size  Seeds/Peers \n--------------------------\n"},
	}

	for _, tt := range table {
		s := getTorrentsTable(tt.in)
		if s != tt.out {
			t.Errorf("\ngot : %v\nwant: %v", s, tt.out)
		}
	}
}

func TestFilterTorrentList(t *testing.T) {

	file, err := os.Open("../../samples/piratebay_search.html")

	if err != nil {
		t.Error(err)
	}

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		t.Error(err)
	}

	scraper := torrents.NewScraper("localhost")

	torrents := scraper.ParseSearchPage(doc)

	var table = []struct {
		in  SearchCommand
		out int
	}{
		{SearchCommand{}, 30},
		{SearchCommand{Trusted: true}, 21},
	}

	for _, tt := range table {
		t.Run(strconv.Itoa(tt.out), func(t *testing.T) {
			s := tt.in.filterTorrentList(torrents)
			if len(s) != tt.out {
				t.Errorf("\ngot: %v\nwant: %v", len(s), tt.out)
			}
		})
	}
}

func TestValidOutputFlags(t *testing.T) {
	var table = []struct {
		label string
		in    SearchCommand
		out   bool
	}{
		{"None", SearchCommand{}, true},
		{"Magnet", SearchCommand{MagnetLinks: true}, true},
		{"URLs", SearchCommand{TorrentURLs: true}, true},
		{"Both", SearchCommand{TorrentURLs: true, MagnetLinks: true}, false},
	}
	for _, tt := range table {
		t.Run(tt.label, func(t *testing.T) {
			s := tt.in.validOutputFlags()
			if s != tt.out {
				t.Errorf("\ngot: %v\nwant: %v", s, tt.out)
			}
		})
	}
}
