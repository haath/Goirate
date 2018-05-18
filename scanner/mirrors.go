package main

import (
	"git.gmantaos.com/haath/Gorrent/shared"
	"github.com/PuerkitoBio/goquery"
)

const url string = "https://proxybay.github.io/"

// Mirror represents a PirateBay mirror and its status.
type Mirror struct {
	url     string
	country string
	status  bool
}

// GetMirrors retrieves a list of PirateBay mirrors.
func GetMirrors() []Mirror {

	doc, _ := shared.HTTPGet(url)

	return parseMirrors(doc)
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
