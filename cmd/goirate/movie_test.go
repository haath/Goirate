package main

import (
	"encoding/json"
	"log"
	"testing"

	"git.gmantaos.com/haath/Goirate/pkg/movies"
)

func TestMovieExecute(t *testing.T) {

	var cmd MovieCommand
	Options.JSON = true

	cmd.Args.Query = "third person"

	output, err := CaptureCommand(cmd.Execute)

	var moviesJSON []movies.Movie
	json.Unmarshal([]byte(output), &moviesJSON)

	cmd.MagnetLink = true

	log.Println(cmd.Args.Query)
	output, err = CaptureCommand(cmd.Execute)

	if err == nil {
		t.Errorf("expected error")
	}

	Options.JSON = false

	output, err = CaptureCommand(cmd.Execute)

	if err != nil {
		t.Error(err)
	}

	cmd.MagnetLink = false
	cmd.Args.Query = "0848228"

	log.Println(cmd.Args.Query)
	output, err = CaptureCommand(cmd.Execute)

	if err != nil {
		t.Error(err)
	}

	cmd.Args.Query = "https://www.imdb.com/title/tt0315983/"

	log.Println(cmd.Args.Query)
	output, err = CaptureCommand(cmd.Execute)

	if err != nil {
		t.Error(err)
	}
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
