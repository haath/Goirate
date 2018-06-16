package main

import (
	"encoding/json"
	imdb "git.gmantaos.com/haath/Goirate/pkg/movies"
	"testing"
)

func TestIMDbExecute(t *testing.T) {

	var cmd IMDbCommand
	Options.JSON = true
	cmd.Args.Query = "avengers"

	output := CaptureCommand(func() { cmd.Execute(nil) })

	var movies []imdb.MovieID
	json.Unmarshal([]byte(output), &movies)

	Options.JSON = false
}

func TestGetMoviesTable(t *testing.T) {
	var table = []struct {
		in  []imdb.MovieID
		out string
	}{
		{[]imdb.MovieID{}, "| IMDb ID | Title | Year |\n|---------|-------|------|\n"},
		{[]imdb.MovieID{imdb.MovieID{IMDbID: "asdf", Title: "Super Awesome film", Year: 2056}}, "| IMDb ID |       Title        | Year |\n|---------|--------------------|------|\n| asdf    | Super Awesome film | 2056 |\n"},
	}

	for _, tt := range table {
		s := getMoviesTable(tt.in)
		if s != tt.out {
			t.Errorf("got %v, want %v", s, tt.out)
		}
	}
}
