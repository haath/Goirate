package main

import (
	"encoding/json"
	"git.gmantaos.com/haath/Gorrent/pkg/piratebay"
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
		{[]piratebay.Torrent{}, " Title  Size  Peers \n--------------------\n"},
	}

	for _, tt := range table {
		s := getTorrentsTable(tt.in)
		if s != tt.out {
			t.Errorf("\ngot : %v\nwant: %v", s, tt.out)
		}
	}
}
