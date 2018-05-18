package main

import (
	"github.com/PuerkitoBio/goquery"
	"os"
	"testing"
)

func TestParseMirrors(t *testing.T) {

	file, err := os.Open("../samples/proxybay.html")
	if err != nil {
		t.Error(err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		t.Error(err)
		return
	}

	mirrors := parseMirrors(doc)

	if len(mirrors) != 16 {
		t.Errorf("Expected to parse 16 mirrors. Found %d", len(mirrors))
	}
}
