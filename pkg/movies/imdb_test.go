package movies

import (
	"reflect"
	"testing"

	"gitlab.com/haath/goirate/pkg/utils"
)

func TestFormatIMDbID(t *testing.T) {
	var table = []struct {
		in  string
		out string
	}{
		{"123", "tt0000123"},
		{"-123", ""},
		{"123456789", "tt123456789"},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {
			s, _ := FormatIMDbID(tt.in)
			if s != tt.out {
				t.Errorf("got %q, want %q", s, tt.out)
			}
		})
	}
}

func TestIsIMDbID(t *testing.T) {
	var table = []struct {
		in  string
		out bool
	}{
		{"123", true},
		{"tt1234567", true},
		{"-123", false},
		{"123456789", true},
		{"Hail ceasar", false},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {
			s := IsIMDbID(tt.in)
			if s != tt.out {
				t.Errorf("got %v, want %v", s, tt.out)
			}
		})
	}

}

func TestIsIMDbURL(t *testing.T) {
	var table = []struct {
		in  string
		out bool
	}{
		{"123", false},
		{"https://www.imdb.com/title/tt0368226/", true},
		{"https://www.imdb.com/title/tt03/", false},
		{"https://www.imdb.com/title/tt6155194/?ref_=nm_knf_i3", true},
		{"https://www.imdb.com/title/tt0848226", true},
		{"https://m.imdb.com/title/tt0848226", true},
		{"Avengers: Age of Ultron", false},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {
			s := IsIMDbURL(tt.in)
			if s != tt.out {

				id, err := ExtractIMDbID(tt.in)

				t.Errorf("got %v, want %v. %v %v", s, tt.out, id, err)
			}
		})
	}

}

func TestParseDuration(t *testing.T) {
	var table = []struct {
		in  string
		out int
	}{
		{"2h 23min", 143},
		{"3h", 180},
		{"47min", 47},
	}
	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {
			s := parseDuration(tt.in)
			if s != tt.out {
				t.Errorf("got %q, want %q", s, tt.out)
			}
		})
	}
}

func TestParseIMDbPage(t *testing.T) {

	expected := Movie{
		MovieID: MovieID{
			IMDbID:   "tt0162222",
			Title:    "Cast Away",
			Year:     2000,
			AltTitle: "Third Person",
		},
		Duration:  143,
		Rating:    7.8,
		PosterURL: "https://m.media-amazon.com/images/M/MV5BN2Y5ZTU4YjctMDRmMC00MTg4LWE1M2MtMjk4MzVmOTE4YjkzXkEyXkFqcGdeQXVyNTc1NTQxODI@._V1_UX182_CR0,0,182,268_AL_.jpg",
		Genres:    []string{"Adventure", "Drama", "Romance"},
	}

	doc, err := utils.GetFileDocument("../../test_samples/imdb.html")

	if err != nil {
		t.Error(err)
		return
	}

	movie := ParseIMDbPage(doc)

	if !reflect.DeepEqual(movie, expected) {
		t.Errorf("got: %v\nwant: %v\n", movie, expected)
	}
}

func TestExtractInfo(t *testing.T) {
	table := []struct {
		in       string
		year     int
		altTitle string
	}{
		{"Fu chou zhe (1976) aka \"Avengers\" ", 1976, "Avengers"},
		{" <a href=\"/title/tt4154756/?ref_=fn_ft_tt_3\">Avengers: Infinity War</a> (2018) ", 2018, ""},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {

			s := extractYear(tt.in)
			a := extractAltTitle(tt.in)

			if s != tt.year {
				t.Errorf("\ngot %v\nwant %v\n", s, tt.year)
			}

			if a != tt.altTitle {
				t.Errorf("\ngot %v\nwant %v\n", s, tt.altTitle)
			}
		})
	}
}

func TestParseSearchPage(t *testing.T) {

	table := []struct {
		index int
		movie MovieID
	}{
		{0, MovieID{"tt0848228", "The Avengers", 2012, ""}},
		{1, MovieID{"tt0164450", "Fu chou zhe", 1976, "Avengers"}},
		{22, MovieID{"tt8277574", "To Avenge", 0, ""}},
		{83, MovieID{"tt0199812", "Ninja Operation 6: Champion on Fire", 1987, "Ninja Avengers"}},
	}

	doc, err := utils.GetFileDocument("../../test_samples/imdb_search.html")

	if err != nil {
		t.Error(err)
		return
	}

	movies := ParseSearchPage(doc)

	if len(movies) != 200 {
		t.Errorf("Expected 200, got %v\n", len(movies))
	}

	for _, tt := range table {
		t.Run(tt.movie.Title, func(t *testing.T) {

			s := movies[tt.index]
			m := tt.movie

			if s != m {
				t.Errorf("\ngot %v\nwant %v\n", s, m)
			}
		})
	}
}

func TestExtractIMDbID(t *testing.T) {
	var table = []struct {
		in  string
		out string
	}{
		{"https://www.imdb.com/title/tt0848228/?ref_=fn_al_tt_1/", "tt0848228"},
		{"https://www.imdb.com/title/tt0848226", "tt0848226"},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {
			s, err := ExtractIMDbID(tt.in)

			if err != nil {
				t.Error(err)
			}

			if s != tt.out {
				t.Errorf("got %q, want %q", s, tt.out)
			}
		})
	}
}

func TestSearchURL(t *testing.T) {
	var table = []struct {
		in  string
		out string
	}{
		{"avengers", "https://www.imdb.com/find?q=avengers&s=tt&ttype=ft"},
		{"Avengers: Age of Ultron", "https://www.imdb.com/find?q=Avengers%253A%2BAge%2Bof%2BUltron&s=tt&ttype=ft"},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {
			s := searchURL(tt.in)

			if s != tt.out {
				t.Errorf("got %q, want %q", s, tt.out)
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
			MovieID: MovieID{
				Title:  "The Avengers",
				IMDbID: "tt0848228",
				Year:   2012,
			},
			Duration:  143,
			Rating:    8.1,
			PosterURL: "https://m.media-amazon.com/images/M/MV5BNDYxNjQyMjAtNTdiOS00NGYwLWFmNTAtNThmYjU5ZGI2YTI1XkEyXkFqcGdeQXVyMTMxODk2OTU@._V1_UX182_CR0,0,182,268_AL_.jpg",
			Genres:    []string{"Action", "Adventure", "Sci-Fi"},
		}},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {

			movie, err := GetMovie(tt.in)

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

func TestSearch(t *testing.T) {

	expected := MovieID{
		Title:  "Avengers: Age of Ultron",
		Year:   2015,
		IMDbID: "tt2395427",
	}

	movies, err := Search("age of ultron")

	if err != nil {
		t.Error(err)
	}

	m := movies[0]

	if m != expected {
		t.Errorf("\ngot %v\nwant %v\n", m, expected)
	}
}
