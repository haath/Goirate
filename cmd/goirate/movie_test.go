package main

import (
	"encoding/json"
	"git.gmantaos.com/haath/Goirate/pkg/movies"
	"testing"
)

func TestMovieExecute(t *testing.T) {

	var cmd MovieCommand
	Options.JSON = true

	cmd.Args.Query = "third person"

	output := CaptureCommand(func() { cmd.Execute(nil) })

	var moviesJSON []movies.Movie
	json.Unmarshal([]byte(output), &moviesJSON)

	Options.JSON = false

	CaptureCommand(func() { cmd.Execute(nil) })
}

func TestFindMovie(t *testing.T) {
	table := []struct {
		in  MovieCommand
		out movies.Movie
	}{
		{MovieCommand{Args: moviePositionalArgs{"age of ultron"}}, movies.Movie{MovieID: movies.MovieID{IMDbID: "2395427"}}},
		{MovieCommand{Args: moviePositionalArgs{"avengers"}, Year: 2018}, movies.Movie{MovieID: movies.MovieID{IMDbID: "4154756"}}},
	}

	for _, tt := range table {
		t.Run(tt.in.Args.Query, func(t *testing.T) {

			s, err := tt.in.findMovie()

			if err != nil {
				t.Error(err)
			}

			if s.IMDbID != tt.out.IMDbID {
				t.Errorf("\ngot: %v\nwant: %v", s, tt.out)
			}
		})
	}
}
