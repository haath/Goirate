package torrents

import (
	"net/url"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
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

func TestIsOk(t *testing.T) {
	var mf MirrorFilters

	table := []struct {
		whitelist []string
		blacklist []string
		mirrorURL string
		isOk      bool
	}{
		{
			[]string{},
			[]string{"pirateproxy.mx"},
			"https://pirateproxy.mx",
			false,
		},
	}

	for _, tt := range table {
		t.Run(tt.mirrorURL, func(t *testing.T) {

			mf.Whitelist = tt.whitelist
			mf.Blacklist = tt.blacklist

			mirror := Mirror{URL: tt.mirrorURL, Country: "US", Status: true}

			isOk := mf.IsOk(mirror)

			if isOk != tt.isOk {
				t.Errorf("got %v, want %v", isOk, tt.isOk)
			}
		})
	}
}

func TestParseMirrors(t *testing.T) {

	table := []Mirror{
		{"https://knaben.xyz", "UK", false},
		{"https://thepbproxy.com", " NL", true},
		{"https://thetorrents.red", "US", true},
		{"https://thetorrents-org.prox.space", "US", true},
		{"https://cruzing.xyz", "US", true},
	}

	file, err := os.Open("../../test_samples/proxybay.html")
	if err != nil {
		t.Error(err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		t.Error(err)
		return
	}

	var mirrorScraper MirrorScraper

	mirrors := mirrorScraper.parseMirrors(doc)

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

	torrents, err := scraper.GetTorrents("ubuntu")

	if err != nil {
		t.Error(err)
	}

	if len(torrents) == 0 {
		t.Errorf("fetched dead mirror")
	}
}

func TestPickMirror(t *testing.T) {

	file, err := os.Open("../../test_samples/proxybay.html")
	if err != nil {
		t.Error(err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		t.Error(err)
		return
	}

	var mirrorScraper MirrorScraper

	mirrors := mirrorScraper.parseMirrors(doc)

	_, torrents, err := mirrorScraper.getTorrents(mirrors, "ubuntu", true)

	if err != nil {
		t.Error(err)
	}

	if len(torrents) == 0 {
		t.Errorf("fetched dead mirror")
	}
}

func TestGetMirrors(t *testing.T) {
	var scraper MirrorScraper

	mirrors, err := scraper.GetMirrors()

	if err != nil {
		t.Error(err)
	}

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
