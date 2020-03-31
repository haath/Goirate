package torrents

import (
	"errors"
	"flag"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gitlab.com/haath/goirate/pkg/utils"
)

const defaultProxySourceURL string = "https://proxybay.github.io/"

// Mirror represents a PirateBay mirror and its status.
type Mirror struct {
	URL     string `json:"url"`
	Country string `json:"country"`
	Status  bool   `json:"status"`
}

// FallbackMirror returns a default Pirate Bay mirror for when the list of mirrors is unavailable.
// The mirror chosen for this, is one that I have personally experienced to be generally available and reliable,
// for the time being, and it should in no way reflect a long-term solution for mirror availability.
func FallbackMirror() Mirror {

	return Mirror{
		URL:     "https://thepiratebay.org",
		Country: "US",
	}
}

// MirrorScraper holds the url to a torrents proxy list.
// By default the scraper will use proxybay.github.io.
type MirrorScraper struct {
	proxySourceURL string
	MirrorFilters  MirrorFilters
}

// MirrorFilters define filters for picking a Pirate Bay mirror.
type MirrorFilters struct {
	Preferred string   `toml:"preferred"`
	Whitelist []string `toml:"whitelist"`
	Blacklist []string `toml:"blacklist"`
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

	return m.parseMirrors(doc), nil
}

// GetTorrents fetches all available Pirate Bay mirrors and returns the first Pirate Bay page that it finds.
func (m *MirrorScraper) GetTorrents(query string) ([]Torrent, error) {

	mirrors, err := m.GetMirrors()

	if err != nil {
		return nil, err
	}

	_, torrents, err := m.getTorrents(mirrors, query, true)

	return torrents, err
}

// PickMirror fetches all available Pirate Bay mirrors and returns the first one that responds to HTTP get for the given query.
func (m *MirrorScraper) PickMirror(query string) (*Mirror, error) {

	mirrors, err := m.GetMirrors()

	if err != nil {
		return nil, err
	}

	mirror, _, err := m.getTorrents(mirrors, query, true)

	return mirror, err
}

// IsOk returns true if the given mirror complies with the filters.
func (m *MirrorFilters) IsOk(mirror Mirror) bool {

	if flag.Lookup("test.v") != nil && strings.Contains(mirror.URL, "thepiratebay.vin") {

		// This mirror sucks so let's exclude it from tests for now

		return false
	}

	contains := func(s []string, e string) bool {

		e = strings.ToLower(e)

		for _, a := range s {

			a = strings.ToLower(a)

			if strings.Contains(e, a) {
				return true
			}
		}
		return false
	}

	return (len(m.Blacklist) == 0 || (!contains(m.Blacklist, mirror.URL) && !contains(m.Blacklist, mirror.Country))) &&
		(len(m.Whitelist) == 0 || (contains(m.Whitelist, mirror.URL) || contains(m.Whitelist, mirror.Country)))
}

func (m *MirrorScraper) parseMirrors(doc *goquery.Document) []Mirror {

	mirrors := make([]Mirror, 0)

	doc.Find("#proxyList > tbody > tr").Each(func(i int, s *goquery.Selection) {

		site, _ := s.Find(".site a").Attr("href")
		country, _ := s.Find(".country img").Attr("alt")
		status, _ := s.Find(".status img").Attr("alt")

		country = strings.ToUpper(country)

		mirror := Mirror{site, country, status == "up"}

		if m.MirrorFilters.IsOk(mirror) {

			mirrors = append(mirrors, mirror)
		}
	})

	return mirrors
}

func (m *MirrorScraper) getTorrents(mirrors []Mirror, query string, trustSource bool) (*Mirror, []Torrent, error) {

	type torrentResponse struct {
		mirror   *Mirror
		torrents []Torrent
	}

	var workingMirror *Mirror

	if m.MirrorFilters.Preferred != "" {

		mirrors = append([]Mirror{{URL: m.MirrorFilters.Preferred}}, mirrors...)
	}

	mirrors = append(mirrors, FallbackMirror())

	timeout := 4 * time.Second

	channel := make(chan torrentResponse)

	searchMirror := func(mirror Mirror) {

		scraper := NewScraper(mirror.URL)

		torrents, err := scraper.SearchTimeout(query, timeout)

		if len(torrents) > 0 {
			workingMirror = &mirror
		}

		if os.Getenv("GOIRATE_DEBUG") == "true" {
			log.Printf("%v -> %v, %v\n", mirror.URL, len(torrents), err)
		}

		if len(torrents) > 0 {

			channel <- torrentResponse{&mirror, torrents}

		} else {

			channel <- torrentResponse{nil, nil}
		}
	}

	// Return the first mirror that responds to HTTP GET
	requestsSent := 0
	for _, mirror := range mirrors {

		if !mirror.Status && trustSource {
			continue
		}

		go searchMirror(mirror)

		requestsSent++
	}

	var allTorrents []Torrent

	for i := 0; i < requestsSent; i++ {

		resp := <-channel

		if resp.torrents != nil {
			allTorrents = append(allTorrents, resp.torrents...)
		}
	}

	if len(allTorrents) > 0 {
		return workingMirror, allTorrents, nil
	}

	if trustSource {

		return m.getTorrents(mirrors, query, false)
	}

	if workingMirror != nil {

		return workingMirror, []Torrent{}, nil
	}

	return nil, nil, errors.New("all Pirate Bay proxies seem to be unreachable")
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
