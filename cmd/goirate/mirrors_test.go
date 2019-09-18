package main

import (
	"encoding/json"
	"testing"

	"gitlab.com/haath/Goirate/pkg/torrents"
)

func TestMirrorsExecute(t *testing.T) {

	var cmd MirrorsCommand
	Options.JSON = true

	output, _ := CaptureCommand(cmd.Execute)

	var mirrors []torrents.Mirror
	json.Unmarshal([]byte(output), &mirrors)

	Options.JSON = false

	_, err := CaptureCommand(cmd.Execute)

	if err != nil {
		t.Error(err)
	}
}

func TestGetMirrorsTable(t *testing.T) {
	var table = []struct {
		in  []torrents.Mirror
		out string
	}{
		{[]torrents.Mirror{}, "|   | Country | URL |\n|---|---------|-----|\n"},
		{[]torrents.Mirror{{URL: "https://pirateproxy.sh", Country: "uk", Status: true}}, "|   | Country |          URL           |\n|---|---------|------------------------|\n| x |   UK    | https://pirateproxy.sh |\n"},
	}

	for _, tt := range table {
		s := getMirrorsTable(tt.in)
		if s != tt.out {
			t.Errorf("got %v, want %v", s, tt.out)
		}
	}
}
