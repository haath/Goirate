package torrents

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
		{"https://pirateproxy.sh", "uk", false},
		{"https://thepbproxy.com", " nl", true},
		{"https://thetorrents.red", "us", true},
		{"https://thetorrents-org.prox.space", "us", true},
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

func TestGetAndPickMirror(t *testing.T) {

	var scraper MirrorScraper

	mirror, err := scraper.PickMirror()

	if err != nil {
		t.Error(err)
	}

	if mirror == nil || !mirror.Status {
		t.Errorf("Fetched dead mirror %v\n", mirror)
	}
}

func TestPickMirror(t *testing.T) {

	expected := Mirror{"https://thepbproxy.com", "nl", true}

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
	mirror, err := pickMirror(mirrors)

	if err != nil {
		t.Error(err)
	}

	if *mirror != expected {
		t.Errorf("got %v, want %v", mirror, expected)
	}
}

func TestGetMirrors(t *testing.T) {
	var scraper MirrorScraper

	mirrors := scraper.GetMirrors()

	if len(mirrors) == 0 {
		t.Errorf("Error fetching PirateBay mirrors.\n")
	}
}

func TestParseLoadTime(t *testing.T) {
	var searchTests = []struct {
		in  string
		out float32
	}{
		{"Loaded in 0.817 seconds", 0.817},
		{"Loaded in 1.850 seconds", 1.85},
		{"Loaded in 1.489 seconds", 1.489},
		{"Loaded in 3.507 seconds", 3.507},
	}

	for _, tt := range searchTests {
		t.Run(tt.in, func(t *testing.T) {
			s := parseLoadTime(tt.in)
			if s != tt.out {
				t.Errorf("got %v, want %v", s, tt.out)
			}
		})
	}
}
