package series

import (
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

func TestNextLastEpisode(t *testing.T) {

	tkn := login(t)

	id := 261690

	expLast := Episode{6, 10}
	expNext := Episode{6, 11}

	last, err := tkn.LastEpisode(id)

	if err != nil {
		t.Error(err)
	}

	if last != expLast {
		t.Errorf("got %v, want %v", last, expLast)
	}

	next, err := tkn.NextEpisode(id, last)

	if err != nil {
		t.Error(err)
	}

	if next != expNext {
		t.Errorf("got %v, want %v", next, expNext)
	}
}
