package main

import (
	"encoding/json"
	"testing"

	imdb "goirate/movies"
)

func TestMovieSearchExecute(t *testing.T) {

	doTestMovieSearchExecute(t)
}

func doTestMovieSearchExecute(t *testing.T) {

	var cmd MovieSearchCommand
	Options.JSON = true
	cmd.Args.Query = "avengers"

	output, err := CaptureCommand(cmd.Execute)

	if err != nil {
		t.Error(output)
		t.Error(err)
	}

	var movies []imdb.MovieID
	json.Unmarshal([]byte(output), &movies)

	Options.JSON = false

	output, err = CaptureCommand(cmd.Execute)

	if err != nil {
		t.Error(output)
		t.Error(err)
	}
}

func TestGetMoviesTable(t *testing.T) {
	var table = []struct {
		in  []imdb.MovieID
		out string
	}{
		{[]imdb.MovieID{}, "| IMDb ID | Title | Year |\n|---------|-------|------|\n"},
		{[]imdb.MovieID{{IMDbID: "asdf", Title: "Super Awesome film", Year: 2056}}, "| IMDb ID |       Title        | Year |\n|---------|--------------------|------|\n| asdf    | Super Awesome film | 2056 |\n"},
	}

	for _, tt := range table {
		s := getMoviesTable(tt.in)
		if s != tt.out {
			t.Errorf("got %v, want %v", s, tt.out)
		}
	}
}
