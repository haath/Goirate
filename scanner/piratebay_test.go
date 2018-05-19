package main

import (
	"testing"
)

var urlTests = []struct {
	in  string
	out string
}{
	{"test_url_123", "test_url_123"},
	{"https://localhost", "https://localhost"},
	{"https://localhost/", "https://localhost/"},
	{"https://localhost:8080/", "https://localhost:8080/"},
}

var searchTests = []struct {
	in  string
	out string
}{
	{"test", "https://pirateproxy.sh/search/test"},
	{"one two", "https://pirateproxy.sh/search/one+two"},
	{"one'two", "https://pirateproxy.sh/search/one%2527two"},
	{"one!", "https://pirateproxy.sh/search/one%2521"},
}

func TestNewScraper(t *testing.T) {
	for _, tt := range urlTests {
		t.Run(tt.in, func(t *testing.T) {
			s := NewScraper(tt.in)
			if s.URL() != tt.out {
				t.Errorf("got %q, want %q", s.URL(), tt.out)
			}
		})
	}
}

func TestSearchURL(t *testing.T) {
	for _, tt := range searchTests {
		t.Run(tt.in, func(t *testing.T) {
			s := NewScraper("https://pirateproxy.sh/")
			searchURL := s.SearchURL(tt.in)
			if searchURL != tt.out {
				t.Errorf("got %q, want %q", searchURL, tt.out)
			}
		})
	}
}
