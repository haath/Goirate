package main

import (
	"git.gmantaos.com/haath/Gorrent/shared"
	"github.com/PuerkitoBio/goquery"
	"log"
	"math"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// PirateBayScaper holds the url of a PirateBay mirror on which to run torrent searches.
type PirateBayScaper interface {
	URL() string
	SearchURL(query string) string
	OnlyVerifiedUploader(enabled bool)
}

type pirateBayScaper struct {
	url                  *url.URL
	onlyVerifiedUploader bool
}

// NewScraper initializes a new PirateBay scapper from a mirror url.
func NewScraper(mirrorURL string) PirateBayScaper {
	URL, err := url.Parse(mirrorURL)

	if err != nil {
		log.Fatalf("Invalid mirror URL: %s\n", mirrorURL)
	}

	var scraper pirateBayScaper
	scraper.url = URL
	return scraper
}

func (s pirateBayScaper) URL() string {
	return s.url.String()
}

func (s pirateBayScaper) OnlyVerifiedUploader(enabled bool) {
	s.onlyVerifiedUploader = enabled
}

func (s pirateBayScaper) SearchURL(query string) string {

	searchURL, _ := url.Parse(s.URL())
	searchURL.Path = path.Join("/search", url.QueryEscape(query))

	return searchURL.String()
}

func (s pirateBayScaper) parseSearchPage(doc *goquery.Document) []shared.Torrent {

	var torrents []shared.Torrent

	doc.Find("#searchResult > tbody > tr").Each(func(i int, s *goquery.Selection) {

		cells := s.Find("td")

		description := cells.Next().Find(".detDesc").Text()
		description = strings.Replace(description, "&nbsp;", " ", -1)
		description = strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}, description)

		title := cells.Next().Find(".detName > .detLink").Text()
		URL, _ := cells.Next().Find(".detName > .detLink").Attr("href")
		magnet, _ := cells.Next().Find("> a").Attr("href")
		seeders, _ := strconv.Atoi(cells.Next().Next().Text())
		leechers, _ := strconv.Atoi(cells.Next().Next().Next().Text())
		verified := s.Find("img[title='VIP'], img[title='Trusted']").Length() > 0

		size := extractSize(description)
		uploadTime := extractUploadTime(description)

		quality := shared.LOW

		torrent := shared.Torrent{
			Title: title, Size: size, Seeders: seeders,
			Leechers: leechers, VerifiedUploader: verified,
			VideoQuality: quality, URL: URL, Magnet: magnet,
			UploadTime: uploadTime,
		}

		torrents = append(torrents, torrent)
	})

	return torrents
}

func extractSize(description string) int {

	r, _ := regexp.Compile("^.+, Size (.+) GiB")
	m := r.FindStringSubmatch(description)

	if len(m) > 0 {
		gb, _ := strconv.ParseFloat(m[len(m)-1], 32)

		return int(math.Round(gb * 1000000))
	}

	r, _ = regexp.Compile("^.+, Size (.+) MiB")
	m = r.FindStringSubmatch(description)

	if len(m) > 0 {
		gb, _ := strconv.ParseFloat(m[len(m)-1], 32)

		return int(math.Round(gb * 1000))
	}

	r, _ = regexp.Compile("^.+, Size (.+) KiB")
	m = r.FindStringSubmatch(description)

	if len(m) > 0 {
		gb, _ := strconv.ParseFloat(m[len(m)-1], 32)

		return int(math.Round(gb))
	}

	return 0.0
}

func extractUploadTime(description string) time.Time {

	r, err := regexp.Compile(`Uploaded\s*(\d\d)-(\d\d)\s*(\d\d):(\d\d)`)

	if err != nil {
		log.Fatal(err)
	}

	m := r.FindStringSubmatch(description)

	if len(m) > 0 {
		day, _ := strconv.Atoi(m[2])
		month, _ := strconv.Atoi(m[1])
		hour, _ := strconv.Atoi(m[3])
		minute, _ := strconv.Atoi(m[4])
		year := time.Now().Year()

		return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)
	}

	r, _ = regexp.Compile(`Uploaded\s*(\d\d)-(\d\d)\s*(\d{4})`)
	m = r.FindStringSubmatch(description)

	log.Println(description)
	day, _ := strconv.Atoi(m[2])
	month, _ := strconv.Atoi(m[1])
	hour := 0
	minute := 0
	year, _ := strconv.Atoi(m[3])

	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)
}
