package main

import (
	"encoding/json"
	"testing"

	"git.gmantaos.com/haath/Goirate/pkg/series"
)

func TestStoreLoadSeries(t *testing.T) {

	storeSeries([]series.Series{})

	ser := series.Series{
		Title:       "Super awesome show",
		LastEpisode: series.Episode{Season: 5, Episode: 11},
		MinQuality:  "720p",
	}

	storeSeries([]series.Series{ser})

	stored := loadSeries()

	if len(stored) != 1 {
		t.Errorf("Stored 1 series by loaded %v", len(stored))
	}

	if stored[0] != ser {
		t.Errorf("\ngot %v\nwant %v\n", stored[0], ser)
	}
}

func TestSeriesCommands(t *testing.T) {

	storeSeries([]series.Series{})

	var addCmd addCommand
	addCmd.Force = true
	addCmd.Args.Title = "the americans 2013"

	expID := 261690
	expEp := series.Episode{Season: 6, Episode: 10}
	expID2 := 280619

	addCmd.Execute([]string{})
	addCmd.Execute([]string{})

	addCmd.Args.Title = "the expanse"
	addCmd.Execute([]string{})

	stored := loadSeries()

	if len(stored) != 2 {
		t.Errorf("Stored 1 series by loaded %v", len(stored))
	}

	if stored[0].ID != expID || stored[0].LastEpisode != expEp {
		t.Errorf("\ngot %v %v\nwant %v\n", stored[0], stored[0].LastEpisode, expID)
	}

	if stored[1].ID != expID2 {
		t.Errorf("\ngot %v\nwant %v\n", stored[1].ID, expID2)
	}

	Options.JSON = true

	var showCmd showCommand

	jsonOut, err := CaptureCommand(showCmd.Execute)

	if err != nil {
		t.Error(err)
	}

	var printedSeriesList []series.Series

	err = json.Unmarshal([]byte(jsonOut), &printedSeriesList)

	if err != nil {
		t.Error(err)
	}

	if len(printedSeriesList) != 2 {
		t.Errorf("expected to print 2 series, instead got:\n%v", jsonOut)
	}

	var rmCmd removeCommand
	rmCmd.Args.Title = "americans"
	rmCmd.Execute([]string{})
	rmCmd.Args.Title = "280619"
	rmCmd.Execute([]string{})

	stored = loadSeries()

	if len(stored) != 0 {
		t.Errorf("expected to have deleted all series, instead got: %v", stored)
	}
}

func TestCapitalize(t *testing.T) {

	table := []struct {
		in  string
		out string
	}{
		{"enLo thEre MatEy", "Enlo There Matey"},
		{"stay away from me booty", "Stay Away From Me Booty"},
		{"ONCE A PIRATE", "Once a Pirate"},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {

			s := capitalize(tt.in)

			if s != tt.out {
				t.Errorf("\ngot: %v\nwant: %v", s, tt.out)
			}
		})
	}
}
