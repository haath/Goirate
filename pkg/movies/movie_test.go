package movies

import (
	"testing"
)

func TestGetURL(t *testing.T) {
	var table = []struct {
		in  string
		out string
	}{
		{"123", "https://www.imdb.com/title/tt0000123/"},
		{"-123", ""},
		{"123456789", ""},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {

			movie := Movie{IMDbID: tt.in}
			url, err := movie.GetURL()

			if (url == nil && tt.out != "") || (err != nil && tt.out != "") {
				t.Errorf("got %v, want %v", url, tt.out)
			}
		})
	}

}

func TestGetMovie(t *testing.T) {
	var table = []struct {
		in  string
		out Movie
	}{
		{"tt0848228", Movie{
			Title:     "The Avengers",
			Duration:  143,
			IMDbID:    "0848228",
			Year:      2012,
			Rating:    8.1,
			PosterURL: "https://m.media-amazon.com/images/M/MV5BNDYxNjQyMjAtNTdiOS00NGYwLWFmNTAtNThmYjU5ZGI2YTI1XkEyXkFqcGdeQXVyMTMxODk2OTU@._V1_UX182_CR0,0,182,268_AL_.jpg",
		}},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {

			movie, err := GetMovie(tt.in)

			if err != nil {
				t.Error(err)
			}

			if *movie != tt.out {
				t.Errorf("got %v, want %v", movie, tt.out)
			}
		})
	}
}
