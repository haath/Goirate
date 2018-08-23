package utils

import "testing"

func TestHTTPGet(t *testing.T) {
	_, err := HTTPGet("https://www.google.com/")

	if err != nil {
		t.Errorf("HTTPGet status code: %s", err.Error())
	}
}

func TestHTTPGetError(t *testing.T) {
	_, err := HTTPGet("1.2.3.4")

	if err == nil {
		t.Errorf("HTTPGet status code: %s", err.Error())
	}
}

func TestHTTPGetStatus(t *testing.T) {

	res, err := HTTPGet("https://www.reddit.com/asdf")

	if err == nil {
		t.Errorf("Expected 404, got: %v", res)
	}
}

func TestGetFileDocument(t *testing.T) {

	expected := "Cast Away (2000) - IMDbTryIMDbProFree"

	doc, err := GetFileDocument("../../samples/imdb.html")

	if err != nil {
		t.Error(err)
	}

	title := doc.Find("title").Text()

	if title != expected {
		t.Errorf("got %v want %v", title, expected)
	}
}
