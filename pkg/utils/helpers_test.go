package utils

import (
	"testing"
)

func TestNormalizeQuery(t *testing.T) {
	table := []struct {
		in  string
		out string
	}{
		{"Spider-Man: Homecoming", "spider man homecoming"},
		{"The Hitchhiker's Guide to the Galaxy", "the hitchhiker s guide to the galaxy"},
		{"American Dad!", "american dad"},
		{"     a     lot    Of!spaces here!  ", "a lot of spaces here"},
	}

	for _, tt := range table {
		t.Run(tt.out, func(t *testing.T) {

			s := NormalizeQuery(tt.in)

			if tt.out != s {
				t.Errorf("\ngot: %v\nwant: %v\n", s, tt.out)
			}
		})
	}
}
