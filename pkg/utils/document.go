package utils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"time"
)

// HTTPGet fetches an HTTP url and returns a goquery.Document
func HTTPGet(url string) (*goquery.Document, error) {

	timeout := time.Duration(3 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	res, err := client.Get(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("http status code: %v", res.StatusCode)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Printf("Error parsing html: %v\n", err)
		return doc, err
	}

	return doc, err
}
