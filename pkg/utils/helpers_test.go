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

func TestNormalizeMediaTitle(t *testing.T) {
	table := []struct {
		in  string
		out string
	}{
		{"House of Cards (US)", "House of Cards"},
		{"The Americans (2013)", "The Americans"},
		{"The Americans [2013]", "The Americans"},
		{"The {12 34} Americans [2013]", "The Americans"},
		{"The   {12 34} Americans   [2013]", "The Americans"},
	}

	for _, tt := range table {
		t.Run(tt.out, func(t *testing.T) {

			s := NormalizeMediaTitle(tt.in)

			if tt.out != s {
				t.Errorf("\ngot: %v\nwant: %v\n", s, tt.out)
			}
		})
	}

}

func TestOverridenBy(t *testing.T) {

	table := []struct {
		gen  OptionalBoolean
		spec OptionalBoolean
		out  bool
	}{
		{Default, Default, false},
		{Default, False, false},
		{Default, True, true},
		{False, Default, false},
		{False, False, false},
		{False, True, true},
		{True, Default, true},
		{True, False, false},
		{True, True, true},
	}

	for _, tt := range table {
		t.Run(string(tt.gen+tt.spec), func(t *testing.T) {

			s := tt.gen.OverridenBy(tt.spec)

			if tt.out != s {
				t.Errorf("\ngot: %v\nwant: %v\n", s, tt.out)
			}
		})
	}
}
