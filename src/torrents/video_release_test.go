package torrents

import (
	"testing"
)

func TestExtractVideoRelease(t *testing.T) {

	expected := []VideoRelease{BDRip, BDRip, BDRip, DVDRip, BDRip, BDRip, DVDRip, DVDRip, "", BDRip, BDRip, BDRip, "", BDRip}

	torrentList, err := OpenTestSample("../../test_samples/piratebay_movie.html")

	if err != nil {
		t.Error(err)
	}

	for i, torrent := range torrentList {

		s := ExtractVideoRelease(torrent.Title)

		if i < len(expected) && expected[i] != s {

			t.Errorf("want %v, got %v\n", expected[i], s)
		}
	}
}
