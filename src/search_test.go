package main

import (
	"encoding/json"
	"testing"

	"goirate/torrents"
)

func TestSearchExecute(t *testing.T) {

	var cmd SearchCommand
	cmd.Args.Query = "avengers"
	Options.JSON = true

	output, _ := CaptureCommand(cmd.Execute)

	var mirrors []torrents.Mirror
	json.Unmarshal([]byte(output), &mirrors)

	cmd.MagnetLink = true

	output, err := CaptureCommand(cmd.Execute)

	if err == nil {
		t.Error(output)
		t.Errorf("Expected error")
	}

	Options.JSON = false

	cmd.SourceURL = "http://localhost"

	output, err = CaptureCommand(cmd.Execute)

	if err == nil {
		t.Error(output)
		t.Errorf("Expected error")
	}
}

func TestGetTorrentsTable(t *testing.T) {
	var table = []struct {
		in  []torrents.Torrent
		out string
	}{
		{[]torrents.Torrent{}, " Title  Size  Seeds/Peers \n---------------------------\n"},
	}

	for _, tt := range table {
		s := getTorrentsTable(tt.in)
		if s != tt.out {
			t.Errorf("\ngot : %v\nwant: %v", s, tt.out)
		}
	}
}
