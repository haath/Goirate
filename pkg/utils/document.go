package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// HTTPGet fetches an HTTP url and returns a goquery.Document.
// It will also set the appropriate headers to make sure the pages are returned in English.
func HTTPGet(url string) (*goquery.Document, error) {

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept-Language", "en-US,en;q=0.8,gd;q=0.6")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Set("X-FORWARDED-FOR", "165.234.102.177")

	res, err := client.Do(req)

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

// GetFileDocument opens an HTML file and returns a GoQuery document from that file.
func GetFileDocument(filePath string) (*goquery.Document, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(file)

	if err != nil {
		return nil, err
	}

	return doc, nil
}
