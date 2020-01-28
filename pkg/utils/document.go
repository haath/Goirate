package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// HTTPClient holds parameters for performing HTTP requests in the library.
type HTTPClient struct {
	// Timeout is the time limit for requests.
	Timeout   time.Duration
	AuthToken string
}

// Get fetches an HTTP url and returns a goquery.Document.
// It will also set the appropriate headers to make sure the pages are returned in English.
func (c *HTTPClient) Get(url string) (*goquery.Document, error) {

	client := http.Client{
		Timeout: c.Timeout,
	}

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Accept-Language", "en-US,en;q=0.8,gd;q=0.6")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	request.Header.Set("X-FORWARDED-FOR", "165.234.102.177")
	request.Close = true

	if c.AuthToken != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer: %v", c.AuthToken))
	}

	res, err := client.Do(request)

	if err != nil {

		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {

		return nil, fmt.Errorf("http: %v -> %v", url, res.StatusCode)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {

		return doc, err
	}

	return doc, err
}

// GetJSON executes an HTTP get request on the given url by serializing the
// given object into JSON.
func (c *HTTPClient) GetJSON(url string, resp interface{}) error {

	client := http.Client{
		Timeout: c.Timeout,
	}

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Accept-Language", "en-US,en;q=0.8,gd;q=0.6")
	request.Header.Set("Accept", "application/json")

	if c.AuthToken != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.AuthToken))
	}

	res, err := client.Do(request)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, resp)

	return err
}

// Post executes an HTTP post request on the given url by serializing the
// given object into JSON.
func (c *HTTPClient) Post(url string, req interface{}, resp interface{}) error {

	client := http.Client{
		Timeout: c.Timeout,
	}

	jsonBytes, err := json.Marshal(req)

	if err != nil {
		return err
	}

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	if c.AuthToken != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer: %v", c.AuthToken))
	}

	res, err := client.Do(request)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, resp)

	return err
}

// HTTPGet fetches an HTTP url and returns a goquery.Document.
// It will also set the appropriate headers to make sure the pages are returned in English.
func HTTPGet(url string) (*goquery.Document, error) {
	var client HTTPClient
	return client.Get(url)
}

// HTTPGetJSON fetches an HTTP url and deserializes the JSON response into resp.
// It will also set the appropriate headers to make sure the pages are returned in English.
func HTTPGetJSON(url string, resp interface{}) error {
	var client HTTPClient
	return client.GetJSON(url, &resp)
}

// HTTPPost executes an HTTP post request on the given url by serializing the
// given object into JSON.
func HTTPPost(url string, req interface{}, resp interface{}) error {
	var client HTTPClient
	err := client.Post(url, req, &resp)
	return err
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
