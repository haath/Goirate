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
		{"walking dead", 153021, "The Walking Dead"},
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
		{261690, Episode{6, 10}},
	}

	tkn := login(t)

	for _, tt := range table {
		t.Run(fmt.Sprint(tt.in), func(t *testing.T) {

			last, err := tkn.LastEpisode(tt.in)
			if err != nil {
				t.Error(err)
			}

			if last != tt.last {
				t.Errorf("got %v, want %v", last, tt.last)
			}

		})
	}
}

func TestNextEpisode(t *testing.T) {

	table := []struct {
		in   int
		last Episode
		next Episode
	}{
		{261690, Episode{6, 10}, Episode{6, 11}},
		{121361, Episode{1, 0}, Episode{1, 1}},
	}

	tkn := login(t)

	for _, tt := range table {
		t.Run(fmt.Sprint(tt.in), func(t *testing.T) {

			next, err := tkn.NextEpisode(tt.in, tt.last)
			if err != nil {
				t.Error(err)
			}

			if next != tt.next {
				t.Errorf("got %v, want %v", next, tt.next)
			}

		})
	}
}
