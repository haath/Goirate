package torrents

import (
	"log"
	"math"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"git.gmantaos.com/haath/Goirate/pkg/utils"
	"git.gmantaos.com/haath/gobytes"
	"github.com/PuerkitoBio/goquery"
)

// PirateBayScaper holds the url of a PirateBay mirror on which to run torrent searches.
type PirateBayScaper interface {
	URL() string
	SearchURL(query string) string
	Search(query string) ([]Torrent, error)
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
	return scraper
}

// FindScraper will use the default MirrorScraper to find a suitable Pirate Bay mirror,
// then return a scraper for that mirror.
func FindScraper() (*PirateBayScaper, error) {
	var mirrorScraper MirrorScraper

	mirror, err := mirrorScraper.PickMirror()

	if err != nil {
		return nil, err
	}

	scraper := NewScraper(mirror.URL)

	return &scraper, nil
}

func (s pirateBayScaper) URL() string {
	return s.url.String()
}

func (s pirateBayScaper) SearchURL(query string) string {

	query = utils.NormalizeQuery(query)

	searchURL, _ := url.Parse(s.URL())
	searchURL.Path = path.Join("/search", url.QueryEscape(query))

	return searchURL.String()
}

func (s pirateBayScaper) Search(query string) ([]Torrent, error) {
	doc, err := utils.HTTPGet(s.SearchURL(query))

	if err != nil {
		return nil, err
	}

	return s.ParseSearchPage(doc), nil
}

func (s pirateBayScaper) ParseSearchPage(doc *goquery.Document) []Torrent {

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
		urlString, _ := cells[1].Find(".detName > .detLink").Attr("href")
		URL, _ := url.Parse(urlString)
		magnet, _ := cells[1].ChildrenFiltered("a").First().Attr("href")
		seeders, _ := strconv.Atoi(cells[2].Text())
		leeches, _ := strconv.Atoi(cells[3].Text())
		verified := row.Find("img[title='VIP'], img[title='Trusted']").Length() > 0
		uploader := cells[1].Find(".detDesc > a.detDesc").Text()

		size := extractSize(description)
		uploadTime := extractUploadTime(description)
		quality := extractVideoQuality(title)

		torrent := Torrent{
			Title:            title,
			Size:             size,
			Seeders:          seeders,
			Leeches:          leeches,
			VerifiedUploader: verified,
			VideoQuality:     quality,
			TorrentURL:       URL.Path,
			Magnet:           magnet,
			UploadTime:       uploadTime,
			MirrorURL:        s.URL(),
			Uploader:         uploader,
		}

		torrents = append(torrents, torrent)
	})

	return torrents
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
	r, err := regexp.Compile(`Uploaded\s*(\d\d)-(\d\d)\s*(\d\d):(\d\d)`)

	if err != nil {
		log.Fatalf("Error extracting upload time: %v\n", err)
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

	if len(m) == 0 {
		log.Printf("Failed at parsing size from: %v\n", description)
	}

	yday := time.Now().AddDate(0, 0, -1)

	day := yday.Day()
	month := yday.Month()
	hour, _ := strconv.Atoi(m[1])
	minute, _ := strconv.Atoi(m[2])
	year := yday.Year()

	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)
}

func extractVideoQuality(title string) VideoQuality {
	title = strings.ToLower(title)
	if strings.Contains(title, string(High)) {
		return High
	} else if strings.Contains(title, string(Medium)) {
		return Medium
	} else if strings.Contains(title, string(Low)) {
		return Low
	}
	return Default
}
