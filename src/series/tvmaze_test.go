package series

import (
	"fmt"
	"testing"
)

func login(t *testing.T) TVmazeToken {

	cred := EnvTVmazeCredentials()

	tkn, err := cred.Login()

	if err != nil {
		t.Error(err)
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
		{"expanse", 1825, "The Expanse"},
		{"strike back", 804, "Strike Back"},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {

			s, err := tkn.SearchFirst(tt.in)

			if err != nil {
				t.Error(err)
			}

			if s.ID != tt.out {
				t.Errorf("got %v, want %v", s.ID, tt.out)
			}

			if s.Name != tt.outName {
				t.Errorf("got %v, want %v", s.Name, tt.outName)
			}
		})
	}
}

func TestLastEpisode(t *testing.T) {

	table := []struct {
		in   int
		last Episode
	}{
		{157, Episode{Season: 6, Episode: 10, Title: "START"}},
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
