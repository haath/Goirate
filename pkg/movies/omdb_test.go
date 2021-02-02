package movies

import (
	"reflect"
	"testing"
)

func TestGetMovieOMDB(t *testing.T) {
	var table = []struct {
		in  string
		out Movie
	}{
		{"tt0848228", Movie{
			MovieID: MovieID{
				Title:  "The Avengers",
				IMDbID: "tt0848228",
				Year:   2012,
			},
			Duration:  143,
			Rating:    8.1,
			PosterURL: "https://m.media-amazon.com/images/M/MV5BNDYxNjQyMjAtNTdiOS00NGYwLWFmNTAtNThmYjU5ZGI2YTI1XkEyXkFqcGdeQXVyMTMxODk2OTU@._V1_SX300.jpg",
			Genres:    []string{"Action", "Adventure", "Sci-Fi"},
		}},
	}

	omdb := EnvOMDBCredentials()

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {

			movie, err := omdb.GetMovie(tt.in)

			if err != nil {
				t.Error(err)
			}

			tt.out.Rating = movie.Rating // Hard-coded tests are bad

			if movie == nil || !reflect.DeepEqual(*movie, tt.out) {
				t.Errorf("\ngot %v\nwant %v", movie, tt.out)
			}
		})
	}
}

func TestSearchOMDB(t *testing.T) {

	expected := MovieID{
		Title:  "Avengers: Age of Ultron",
		Year:   2015,
		IMDbID: "tt2395427",
	}

	omdb := EnvOMDBCredentials()

	movies, err := omdb.Search("age of ultron")

	if err != nil {
		t.Error(err)
	}

	m := movies[0]

	if m != expected {
		t.Errorf("\ngot %v\nwant %v\n", m, expected)
	}
}
