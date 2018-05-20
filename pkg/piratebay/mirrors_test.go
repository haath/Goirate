package piratebay

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"os"
	"testing"
)

func TestMirrorScraper(t *testing.T) {
	var m MirrorScraper

	_, err := url.Parse(m.GetProxySourceURL())

	if m.GetProxySourceURL() == "" || err != nil {
		t.Errorf("Invalid default ProxySourceURL %v\n", m.GetProxySourceURL())
	}

	m.SetProxySourceURL("https://localhost/")
	_, err = url.Parse(m.GetProxySourceURL())

	if m.GetProxySourceURL() != "https://localhost/" || err != nil {
		t.Errorf("Error setting ProxySourceURL %v\n", m.GetProxySourceURL())
	}
}

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

func TestGetMirrors(t *testing.T) {
	var scraper MirrorScraper

	mirrors := scraper.GetMirrors()

	if len(mirrors) == 0 {
		t.Errorf("Error fetching PirateBay mirrors.\n")
	}
}
