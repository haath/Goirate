package main

import (
	"encoding/json"
	"git.gmantaos.com/haath/Gorrent/pkg/piratebay"
	"github.com/PuerkitoBio/goquery"
	"os"
	"strconv"
	"testing"
)

func TestSearchExecute(t *testing.T) {

	var cmd SearchCommand
	Options.JSON = true

	output := CaptureCommand(func() { cmd.Execute([]string{"avengers"}) })

	var mirrors []piratebay.Mirror
	json.Unmarshal([]byte(output), &mirrors)

	Options.JSON = false
}

func TestGetTorrentsTable(t *testing.T) {
	var table = []struct {
		in  []piratebay.Torrent
		out string
	}{
		{[]piratebay.Torrent{}, " Title  Size  Seeds/Peers \n--------------------------\n"},
	}

	for _, tt := range table {
		s := getTorrentsTable(tt.in)
		if s != tt.out {
			t.Errorf("\ngot : %v\nwant: %v", s, tt.out)
		}
	}
}

func TestFilterTorrentList(t *testing.T) {

	file, err := os.Open("../samples/piratebay_search.html")

	if err != nil {
		t.Error(err)
	}

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		t.Error(err)
	}

	scraper := piratebay.NewScraper("localhost")

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
