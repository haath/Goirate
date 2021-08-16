package main

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"goirate/utils"
)

func CaptureCommand(cmd func([]string) error) (string, error) {
	log.SetFlags(0)
	var buf bytes.Buffer
	log.SetOutput(&buf)
	err := cmd(nil)
	log.SetOutput(&utils.GoirateLogger{})

	var filtered bytes.Buffer

	for _, line := range strings.Split(buf.String(), "\n") {

		if line != "" && !strings.Contains(line, "Unsolicited response received on idle HTTP channel") {

			filtered.WriteString(line + "\n")
		}
	}

	return filtered.String(), err
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
