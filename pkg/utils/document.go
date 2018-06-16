package utils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"time"
)

// HTTPGet fetches an HTTP url and returns a goquery.Document
func HTTPGet(url string) (*goquery.Document, error) {

	timeout := time.Duration(3 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept-Language", "en")

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
