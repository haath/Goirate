package main

import (
	"log"
	"net/url"
	"path"
)

// PirateBayScaper holds the url of a PirateBay mirror on which to run torrent searches.
type PirateBayScaper interface {
	URL() string
	SearchURL(query string) string
}

type pirateBayScaper struct {
	url *url.URL
}

// NewScraper initializes a new PirateBay scapper from a mirror url.
func NewScraper(mirrorURL string) PirateBayScaper {
	URL, err := url.Parse(mirrorURL)

	if err != nil {
		log.Fatalf("Invalid mirror URL: %s\n", mirrorURL)
	}

	return pirateBayScaper{URL}
}

func (s pirateBayScaper) URL() string {
	return s.url.String()
}

func (s pirateBayScaper) SearchURL(query string) string {

	searchURL, _ := url.Parse(s.URL())
	searchURL.Path = path.Join("/search", url.QueryEscape(query))

	return searchURL.String()
}
