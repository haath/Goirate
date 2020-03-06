package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func CaptureCommand(cmd func([]string) error) (string, error) {
	log.SetFlags(0)
	var buf bytes.Buffer
	log.SetOutput(&buf)
	err := cmd(nil)
	log.SetOutput(os.Stdout)

	var filtered bytes.Buffer

	for _, line := range strings.Split(buf.String(), "\n") {

		if !strings.Contains(line, "Unsolicited response received on idle HTTP channel") {

			filtered.WriteString(line)
		}
	}

	return buf.String(), err
}

func TestGetScraper(t *testing.T) {
	table := []struct {
		in       torrentSearchArgs
		outURL   string
		outError bool
	}{
		{torrentSearchArgs{}, "", false},
		{torrentSearchArgs{Mirror: "http://1.2.3.4/"}, "http://1.2.3.4/", false},
		{torrentSearchArgs{SourceURL: "http://1.2.3.4/"}, "", true},
	}

	for _, tt := range table {
		t.Run(tt.outURL, func(t *testing.T) {
			cmd := tt.in
			scraper, err := cmd.GetScraper("ubuntu")

			if !tt.outError && err != nil {
				t.Error(err)
				return
			}

			if !tt.outError && scraper == nil {
				t.Error(err)
				return
			}

			if !tt.outError && scraper.URL() != tt.outURL && tt.outURL != "" {
				t.Errorf("\ngot: %v\nwant: %v", scraper.URL(), tt.outURL)
			}
		})
	}
}

func TestValidOutputFlags(t *testing.T) {
	var table = []struct {
		label string
		in    func() torrentSearchArgs
		out   bool
	}{
		{"None", func() torrentSearchArgs { return torrentSearchArgs{} }, true},
		{"Magnet", func() torrentSearchArgs {
			cmd := torrentSearchArgs{}
			cmd.MagnetLink = true
			return cmd
		}, true},
		{"URLs", func() torrentSearchArgs {
			cmd := torrentSearchArgs{}
			cmd.TorrentURL = true
			return cmd
		}, true},
		{"Both", func() torrentSearchArgs {
			cmd := torrentSearchArgs{}
			cmd.TorrentURL = true
			cmd.MagnetLink = true
			return cmd
		}, false},
	}
	for _, tt := range table {
		t.Run(tt.label, func(t *testing.T) {
			cmd := tt.in()
			s := cmd.ValidOutputFlags()
			if s != tt.out {
				t.Errorf("\ngot: %v\nwant: %v", s, tt.out)
			}
		})
	}
}
