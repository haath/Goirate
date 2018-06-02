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
