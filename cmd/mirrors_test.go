package main

import (
	"encoding/json"
	"git.gmantaos.com/haath/Gorrent/pkg/piratebay"
	"testing"
)

func TestMirrorsExecute(t *testing.T) {

	var cmd MirrorsCommand
	Options.JSON = true

	output := CaptureCommand(func() { cmd.Execute(nil) })

	var mirrors []piratebay.Mirror
	json.Unmarshal([]byte(output), &mirrors)

	Options.JSON = false
}

func TestGetMirrorsTable(t *testing.T) {
	var table = []struct {
		in  []piratebay.Mirror
		out string
	}{
		{[]piratebay.Mirror{}, "|   | Country | URL |\n|---|---------|-----|\n"},
		{[]piratebay.Mirror{piratebay.Mirror{URL: "https://pirateproxy.sh", Country: "uk", Status: true}}, "|   | Country |          URL           |\n|---|---------|------------------------|\n| x |   UK    | https://pirateproxy.sh |\n"},
	}

	for _, tt := range table {
		s := getMirrorsTable(tt.in)
		if s != tt.out {
			t.Errorf("got %v, want %v", s, tt.out)
		}
	}
}
