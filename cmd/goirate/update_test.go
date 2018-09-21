package main

import "testing"

func TestMoreRecentThan(t *testing.T) {

	table := []struct {
		first  string
		second string
		out    bool
	}{
		{"1.2.3", "1.2.4", false},
		{"1.2.4", "1.2.4", false},
		{"1.2.5", "1.2.4", true},
		{"v1.2.5", "1.2.4", true},
		{"1.2.5", "v1.2.4", true},
	}

	for _, tt := range table {
		t.Run(tt.first+" > "+tt.second, func(t *testing.T) {

			first, _ := parseVersion(tt.first)
			second, _ := parseVersion(tt.second)

			s := first.moreRecentThan(second)

			if s != tt.out {
				t.Errorf("\ngot: %v\nwant: %v", s, tt.out)
			}
		})
	}
}
