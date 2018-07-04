package torrents

import (
	"errors"
	"regexp"
	"strconv"

	"git.gmantaos.com/haath/Goirate/pkg/utils"
	"github.com/PuerkitoBio/goquery"
)

const defaultProxySourceURL string = "https://proxybay.github.io/"

// Mirror represents a PirateBay mirror and its status.
type Mirror struct {
	URL     string `json:"url"`
	Country string `json:"country"`
	Status  bool   `json:"status"`
}

// MirrorScraper holds the url to a torrents proxy list.
// By default the scraper will use proxybay.github.io.
type MirrorScraper struct {
	proxySourceURL string
}

// SetProxySourceURL overrides the URL at which MirrorScraper will attempt to fetch a list
// of Pirate Bay proxies from.
func (m *MirrorScraper) SetProxySourceURL(url string) {
	m.proxySourceURL = url
}

// GetProxySourceURL retrieves the current URL at which the scraper will attempt to fetch a list
// of Pirate Bay proxies from.
func (m *MirrorScraper) GetProxySourceURL() string {
	if m.proxySourceURL == "" {
		return defaultProxySourceURL
	}
	return m.proxySourceURL
}

// GetMirrors retrieves a list of PirateBay mirrors.
func (m *MirrorScraper) GetMirrors() ([]Mirror, error) {

	doc, err := utils.HTTPGet(m.GetProxySourceURL())

	if err != nil {
		return nil, err
	}

	return parseMirrors(doc), nil
}

// PickMirror fetches all available Pirate Bay mirrors and picks the the fastest one available.
func (m *MirrorScraper) PickMirror(query string) (*Mirror, error) {
	mirrors, err := m.GetMirrors()

	if err != nil {
		return nil, err
	}

	return pickMirror(mirrors, query, true)
}

func parseMirrors(doc *goquery.Document) []Mirror {

	mirrors := make([]Mirror, 0)

	doc.Find("#proxyList > tbody > tr").Each(func(i int, s *goquery.Selection) {
		site, _ := s.Find(".site a").Attr("href")
		country, _ := s.Find(".country img").Attr("alt")
		status, _ := s.Find(".status img").Attr("alt")

		mirror := Mirror{site, country, status == "up"}

		mirrors = append(mirrors, mirror)
	})

	return mirrors
}

func parseLoadTime(speedTitle string) float32 {
	r, _ := regexp.Compile("Loaded in (\\-?\\d+\\.\\d+) seconds")
	m := r.FindStringSubmatch(speedTitle)

	if len(m) > 0 {
		val, _ := strconv.ParseFloat(m[1], 32)

		return float32(val)
	}
	return 0.0
}

func pickMirror(mirrors []Mirror, query string, trustSource bool) (*Mirror, error) {

	// Return the first mirror that responds to HTTP GET
	for _, mirror := range mirrors {

		if !mirror.Status && trustSource {
			continue
		}

		scraper := NewScraper(mirror.URL)

		doc, err := utils.HTTPGet(scraper.SearchURL(query))

		if err == nil {

			torrents := scraper.ParseSearchPage(doc)

			if len(torrents) > 0 {
				return &mirror, nil
			}

		}
	}

	if trustSource {
		return pickMirror(mirrors, query, false)
	}

	return nil, errors.New("all Pirate Bay proxies seem to be unreachable")
}
