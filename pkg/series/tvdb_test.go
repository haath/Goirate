package series

import (
	"fmt"
	"testing"
)

func login(t *testing.T) TVDBToken {

	cred := EnvTVDBCredentials()

	tkn, err := cred.Login()

	if err != nil {
		t.Error(err)
	}

	if tkn.Token == "" {
		t.Errorf("Got empty token")
	}

	return tkn
}

func TestSearch(t *testing.T) {

	tkn := login(t)

	table := []struct {
		in      string
		out     int
		outName string
	}{
		{"expanse", 280619, "The Expanse"},
		{"strike back", 148581, "Strike Back"},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {

			s, n, err := tkn.Search(tt.in)

			if err != nil {
				t.Error(err)
			}

			if s != tt.out {
				t.Errorf("got %v, want %v", s, tt.out)
			}

			if n != tt.outName {
				t.Errorf("got %v, want %v", n, tt.outName)
			}
		})
	}
}

func TestLastEpisode(t *testing.T) {

	table := []struct {
		in   int
		last Episode
	}{
		{261690, Episode{Season: 6, Episode: 10, Title: "START"}},
	}

	tkn := login(t)

	for _, tt := range table {
		t.Run(fmt.Sprint(tt.in), func(t *testing.T) {

			last, err := tkn.LastEpisode(tt.in)
			if err != nil {
				t.Error(err)
			}

			if last.String() != tt.last.String() || last.Title != tt.last.Title {
				t.Errorf("got %v, want %v", last.Aired, tt.last)
			}

		})
	}
}
