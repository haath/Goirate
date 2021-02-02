package main

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"gitlab.com/haath/goirate/pkg/movies"
)

func TestMovieExecute(t *testing.T) {

	doMovieExecuteTest(t)
}

func TestMovieExecuteFallback(t *testing.T) {

	// Unset OMDb key to test fallback.
	curKey := Config.OMDBCredentials.APIKey
	Config.OMDBCredentials.APIKey = ""
	os.Unsetenv("GOIRATE_OMDB_API_KEY")

	doMovieExecuteTest(t)

	// Restore OMDb key
	os.Setenv("GOIRATE_OMDB_API_KEY", curKey)
}

func doMovieExecuteTest(t *testing.T) {

	var cmd MovieCommand
	Options.JSON = true

	cmd.Args.Query = "the avengers"

	output, err := CaptureCommand(cmd.Execute)

	if err != nil {
		log.Println(output)
		t.Fatal(err)
	}

	var moviesJSON []movies.Movie
	json.Unmarshal([]byte(output), &moviesJSON)

	cmd.MagnetLink = true

	Options.JSON = false

	output, err = CaptureCommand(cmd.Execute)

	if err != nil {
		log.Println(output)
		t.Fatal(err)
	}

	cmd.MagnetLink = false
	cmd.Args.Query = "0848228"

	output, err = CaptureCommand(cmd.Execute)

	if err != nil {
		log.Println(output)
		t.Fatal(err)
	}

	cmd.Args.Query = "https://www.imdb.com/title/tt0315983/"

	output, err = CaptureCommand(cmd.Execute)

	if err != nil {
		log.Println(output)
		t.Fatal(err)
	}
}

func TestSearchMovie(t *testing.T) {

	doTestSearchMovie(t)
}

func TestSearchMovieFallback(t *testing.T) {

	// Unset OMDb key to test fallback.
	curKey := Config.OMDBCredentials.APIKey
	Config.OMDBCredentials.APIKey = ""
	os.Unsetenv("GOIRATE_OMDB_API_KEY")

	doTestSearchMovie(t)

	// Restore OMDb key
	os.Setenv("GOIRATE_OMDB_API_KEY", curKey)
}

func doTestSearchMovie(t *testing.T) {
	table := []struct {
		in  MovieCommand
		out movies.Movie
	}{
		{MovieCommand{Args: moviePositionalArgs{"age of ultron"}}, movies.Movie{MovieID: movies.MovieID{IMDbID: "tt2395427"}}},
		{MovieCommand{Args: moviePositionalArgs{"avengers"}, Year: 2018}, movies.Movie{MovieID: movies.MovieID{IMDbID: "tt4154756"}}},
	}

	for _, tt := range table {
		t.Run(tt.in.Args.Query, func(t *testing.T) {

			s, err := tt.in.searchMovie()

			if err != nil {
				t.Error(err)
			}

			if s != tt.out.IMDbID {
				t.Errorf("\ngot: %v\nwant: %v", s, tt.out)
			}
		})
	}
}
