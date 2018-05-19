package piratebay

import (
	"github.com/PuerkitoBio/goquery"
	"os"
	"testing"
)

func TestParseMirrors(t *testing.T) {

	table := []Mirror{
		{"https://pirateproxy.sh", "uk", true},
		{"https://thepbproxy.com", " nl", true},
		{"https://thepiratebay.red", "us", true},
		{"https://thepiratebay-org.prox.space", "us", true},
		{"https://cruzing.xyz", "us", true},
	}

	file, err := os.Open("../../samples/proxybay.html")
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
		t.Errorf("Expected to parse 16 mirrors. Found %d.\n", len(mirrors))
	}

	for i := range table {
		e := table[i]
		a := mirrors[i]

		if e.URL != a.URL || e.Country != e.Country || e.Status != a.Status {
			t.Errorf("Wrong mirror parsing. Expected %v, got %v.\n", e, a)
		}
	}
}
