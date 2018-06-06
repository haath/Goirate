package utils

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

// HTTPGet fetches an HTTP url and returns a goquery.Document
func HTTPGet(url string) (*goquery.Document, error) {
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, err
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc, err
}
