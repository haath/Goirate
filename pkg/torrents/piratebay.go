package torrents

import (
	"fmt"
	"log"
	"math"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"gitlab.com/haath/gobytes"
	"gitlab.com/haath/goirate/pkg/utils"
)

// PirateBayScaper holds the url of a PirateBay mirror on which to run torrent searches.
type PirateBayScaper interface {
	URL() string
	SearchURLs(query string) []string
	Search(query string) ([]Torrent, error)
	SearchTimeout(query string, timeout time.Duration) ([]Torrent, error)
	SearchVideoTorrents(query string, filters SearchFilters) ([]Torrent, error)
	ParseSearchPage(doc *goquery.Document) []Torrent
}

type pirateBayScaper struct {
	url *url.URL
}

// NewScraper initializes a new PirateBay scapper from a mirror url.
func NewScraper(mirrorURL string) PirateBayScaper {
	URL, err := url.Parse(mirrorURL)

	if err != nil {
		log.Fatalf("Invalid mirror URL: %s\n", mirrorURL)
	}

	var scraper pirateBayScaper
	scraper.url = URL
	return &scraper
}

// FindScraper will use the default MirrorScraper to find a suitable Pirate Bay mirror,
// then return a scraper for that mirror.
func FindScraper(query string) (PirateBayScaper, error) {
	var mirrorScraper MirrorScraper

	mirror, err := mirrorScraper.PickMirror(query)

	if err != nil {
		return nil, err
	}

	scraper := NewScraper(mirror.URL)

	return scraper, nil
}

func (s *pirateBayScaper) URL() string {
	return s.url.String()
}

func (s *pirateBayScaper) SearchURLs(query string) []string {

	query = utils.NormalizeQuery(query)

	var urls []string

	searchURL, _ := url.Parse(s.URL())

	// First url (legacy)
	searchURL.Path = path.Join("/search", query)
	urls = append(urls, searchURL.String())

	// Second url (piratesbaycc.com)
	searchURL.Path = "/search.php"
	queryBuilder := searchURL.Query()
	queryBuilder.Set("orderby", "99")
	queryBuilder.Set("page", "0")
	queryBuilder.Set("q", url.QueryEscape(query))
	searchURL.RawQuery = queryBuilder.Encode()
	urls = append(urls, searchURL.String())

	// third url (knaben)
	searchURL.Path = "/s/"
	queryBuilder = searchURL.Query()
	queryBuilder.Set("orderby", "99")
	queryBuilder.Set("page", "0")
	queryBuilder.Set("category", "0")
	queryBuilder.Set("q", url.QueryEscape(query))
	searchURL.RawQuery = queryBuilder.Encode()
	urls = append(urls, searchURL.String())

	return urls
}

func (s *pirateBayScaper) SearchTimeout(query string, timeout time.Duration) ([]Torrent, error) {

	return s.search(query, timeout)
}

func (s *pirateBayScaper) Search(query string) ([]Torrent, error) {

	return s.search(query, 0)
}

func (s *pirateBayScaper) ParseSearchPage(doc *goquery.Document) []Torrent {

	var torrents []Torrent

	doc.Find("#searchResult > tbody > tr").Each(func(i int, row *goquery.Selection) {

		if row.Find("td").Length() < 2 {
			// Hit the pagination row
			return
		}

		cells := []*goquery.Selection{
			row.Find("td").First(),
			row.Find("td").Next().First(),
			row.Find("td").Next().Next().First(),
			row.Find("td").Next().Next().Next().First(),
		}

		description := cells[1].Find(".detDesc").Text()
		description = strings.Replace(description, "&nbsp;", " ", -1)
		description = strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}, description)

		title := cells[1].Find(".detName > .detLink").Text()
		urlPathString, _ := cells[1].Find(".detName > .detLink").Attr("href")
		urlPath, _ := url.Parse(urlPathString)
		magnet, _ := cells[1].ChildrenFiltered("a").First().Attr("href")
		seeders, _ := strconv.Atoi(cells[2].Text())
		leeches, _ := strconv.Atoi(cells[3].Text())
		verified := row.Find("img[title='VIP'], img[title='Trusted']").Length() > 0
		uploader := cells[1].Find(".detDesc > a.detDesc").Text()

		size := extractSize(description)
		uploadTime := extractUploadTime(description)
		quality := extractVideoQuality(title)
		releaseType := ExtractVideoRelease(title)

		torrentURLPath := urlPath.Path
		if urlPath.RawQuery != "" {
			torrentURLPath += fmt.Sprintf("?%s", urlPath.RawQuery)
		}

		torrent := Torrent{
			Title:            title,
			Size:             size,
			Seeders:          seeders,
			Leeches:          leeches,
			VerifiedUploader: verified,
			VideoQuality:     quality,
			VideoRelease:     releaseType,
			TorrentURL:       torrentURLPath,
			Magnet:           magnet,
			UploadTime:       uploadTime,
			MirrorURL:        s.URL(),
			Uploader:         uploader,
		}

		torrents = append(torrents, torrent)
	})

	return torrents
}

func (s *pirateBayScaper) SearchVideoTorrents(query string, filters SearchFilters) ([]Torrent, error) {

	return s.SearchVideoTorrentsTimeout(query, filters, 0)
}

func (s *pirateBayScaper) SearchVideoTorrentsTimeout(query string, filters SearchFilters, timeout time.Duration) ([]Torrent, error) {

	query = utils.NormalizeQuery(query)

	trnts, err := s.SearchTimeout(query, timeout)

	if err != nil {
		return nil, err
	}

	filtered := filters.FilterTorrents(trnts)

	perQuality, err := SearchVideoTorrentList(filtered, filters)

	if err != nil {
		return nil, err
	}

	var perQualitySlice []Torrent
	for _, value := range perQuality {
		perQualitySlice = append(perQualitySlice, *value)
	}

	if os.Getenv("GOIRATE_DEBUG") == "true" {
		log.Printf("%d => %d => %d", len(trnts), len(filtered), len(perQuality))
	}

	return perQualitySlice, nil
}

func (s *pirateBayScaper) GetNextPageURL(doc *goquery.Document) string {

	a := doc.Find("img[alt='Next']").Parent()

	relative, _ := a.Attr("href")

	return path.Join(s.URL(), relative)
}

func extractSize(description string) int64 {

	r, _ := regexp.Compile(`Size\s*(.+)\s*GiB`)
	m := r.FindStringSubmatch(description)

	if len(m) > 0 {
		gb, _ := strconv.ParseFloat(strings.TrimSpace(m[1]), 32)

		return int64(math.Round(gb * gobytes.GiB.KBytes()))
	}

	r, _ = regexp.Compile(`Size\s*(.+)\s*MiB`)
	m = r.FindStringSubmatch(description)

	if len(m) > 0 {
		mb, _ := strconv.ParseFloat(strings.TrimSpace(m[1]), 32)

		return int64(math.Round(mb * gobytes.MiB.KBytes()))
	}

	r, _ = regexp.Compile(`Size\s*(.+)\s*KiB`)
	m = r.FindStringSubmatch(description)

	if len(m) > 0 {
		kb, _ := strconv.ParseFloat(strings.TrimSpace(m[1]), 32)

		return int64(math.Round(kb * gobytes.KiB.KBytes()))
	}

	return 0.0
}

func extractUploadTime(description string) time.Time {

	/*
		First check the MM-DD HH:mm format
	*/
	r, _ := regexp.Compile(`Uploaded\s*(\d\d)-(\d\d)\s*(\d\d):(\d\d)`)
	m := r.FindStringSubmatch(description)

	if len(m) > 0 {
		day, _ := strconv.Atoi(m[2])
		month, _ := strconv.Atoi(m[1])
		hour, _ := strconv.Atoi(m[3])
		minute, _ := strconv.Atoi(m[4])
		year := time.Now().Year()

		return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)
	}

	/*
		Next check the MM-DD YYYY format
	*/
	r, _ = regexp.Compile(`Uploaded\s*(\d\d)-(\d\d)\s*(\d{4})`)
	m = r.FindStringSubmatch(description)

	if len(m) > 0 {
		day, _ := strconv.Atoi(m[2])
		month, _ := strconv.Atoi(m[1])
		hour := 0
		minute := 0
		year, _ := strconv.Atoi(m[3])

		return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)
	}

	/*
		Check the Today YYYY format
	*/
	r, _ = regexp.Compile(`Uploaded\s*Today\s*(\d\d):(\d\d)`)
	m = r.FindStringSubmatch(description)

	if len(m) > 0 {
		day := time.Now().Day()
		month := time.Now().Month()
		hour, _ := strconv.Atoi(m[1])
		minute, _ := strconv.Atoi(m[2])
		year := time.Now().Year()

		return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)
	}

	/*
		Finally, check the Y-day YYYY format
	*/
	r, _ = regexp.Compile(`Uploaded\s*Y-day\s*(\d\d):(\d\d)`)
	m = r.FindStringSubmatch(description)

	if len(m) > 0 {
		yday := time.Now().AddDate(0, 0, -1)

		day := yday.Day()
		month := yday.Month()
		hour, _ := strconv.Atoi(m[1])
		minute, _ := strconv.Atoi(m[2])
		year := yday.Year()

		return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)
	}

	r, _ = regexp.Compile(`Uploaded\s*(\d+)\s*mins\s*ago`)
	m = r.FindStringSubmatch(description)

	if len(m) == 0 {
		return time.Time{}
	}

	minutes, _ := strconv.Atoi(m[1])
	minutesAgo := time.Now().Add(time.Duration(-minutes) * time.Minute)

	return minutesAgo
}

func extractVideoQuality(title string) VideoQuality {

	quality := Default
	title = utils.NormalizeQuery(title)
	words := strings.Fields(title)

	containsQualityTerms := func(terms ...string) bool {
		for _, w1 := range words {
			for _, w2 := range terms {
				if w1 == w2 {
					return true
				}
			}
		}
		return false
	}

	if containsQualityTerms(string(UHD), "4k", "uhd", "ultrahd") {
		quality = UHD
	} else if containsQualityTerms(string(High)) {
		quality = High
	} else if containsQualityTerms(string(Medium)) {
		quality = Medium
	} else if containsQualityTerms(string(Low)) {
		quality = Low
	}
	return quality
}

func (s *pirateBayScaper) search(query string, timeout time.Duration) ([]Torrent, error) {

	searchURLs := s.SearchURLs(query)

	torrentsChannel := make(chan []Torrent)
	errorChannel := make(chan error)

	for _, searchURL := range searchURLs {

		searchURLformatted := strings.Replace(searchURL, "%2B", "+", -1)

		if os.Getenv("GOIRATE_DEBUG") == "true" {
			log.Printf("Search url: %s\n", searchURLformatted)
		}

		go func() {

			client := utils.HTTPClient{
				Timeout: timeout,
			}

			doc, err := client.Get(searchURLformatted)

			if err == nil {

				torrentsChannel <- s.ParseSearchPage(doc)
				errorChannel <- nil

			} else {

				torrentsChannel <- nil
				errorChannel <- err

			}
		}()
	}

	var allTorrents []Torrent
	var allError error

	for i := 0; i < len(searchURLs); i++ {

		torrents := <-torrentsChannel
		err := <-errorChannel

		if err != nil {
			allError = err
		}

		if torrents != nil {
			allTorrents = append(allTorrents, torrents...)
		}
	}

	return allTorrents, allError
}
