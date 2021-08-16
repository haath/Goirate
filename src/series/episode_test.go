package series

import (
	"fmt"
	"testing"
)

func TestParseEpisodeString(t *testing.T) {
	table := []struct {
		in      string
		outSes  uint
		outEp   uint
		outStr  string
		outLong string
	}{
		{"S05E12", 5, 12, "S05E12", "Season 5, Episode 12"},
		{"S 05 E 12", 5, 12, "S05E12", "Season 5, Episode 12"},
		{"S1234 E 12", 1234, 12, "S1234E12", "Season 1234, Episode 12"},
		{"Season 12 episode 5", 12, 5, "S12E05", "Season 12, Episode 5"},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {

			ep := ParseEpisodeString(tt.in)

			if ep.Season != tt.outSes || ep.Episode != tt.outEp {
				t.Errorf("got s%ve%v, want s%ve%v", ep.Season, ep.Episode, tt.outSes, tt.outEp)
			}

			if ep.String() != tt.outStr {
				t.Errorf("got %v, want %v", ep.String(), tt.outStr)
			}

			if ep.LongString() != tt.outLong {
				t.Errorf("got %v, want %v", ep.LongString(), tt.outLong)
			}
		})
	}
}

func TestIsAfter(t *testing.T) {

	table := []struct {
		aft []uint
		bef []uint
		out bool
	}{
		{[]uint{0, 0}, []uint{0, 3}, false},
		{[]uint{0, 3}, []uint{0, 3}, false},
		{[]uint{0, 4}, []uint{0, 3}, true},
		{[]uint{3, 10}, []uint{0, 12}, true},
		{[]uint{0, 12}, []uint{3, 10}, false},
	}

	for _, tt := range table {
		t.Run(fmt.Sprintf("%v > %v", tt.aft, tt.bef), func(t *testing.T) {

			aft := Episode{Season: tt.aft[0], Episode: tt.aft[1]}
			bef := Episode{Season: tt.bef[0], Episode: tt.bef[1]}

			s := aft.IsAfter(bef)

			if tt.out != s {
				t.Errorf("got %v, want %v", s, tt.out)
			}
		})
	}
}
