package main

import (
	"encoding/json"
	"testing"

	imdb "git.gmantaos.com/haath/Goirate/pkg/movies"
)

func TestMovieSearchExecute(t *testing.T) {

	var cmd MovieSearchCommand
	Options.JSON = true
	cmd.Args.Query = "avengers"

	output, _ := CaptureCommand(cmd.Execute)

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
