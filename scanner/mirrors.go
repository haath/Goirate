package main

import (
	"encoding/json"
	"git.gmantaos.com/haath/Gorrent/shared"
	"github.com/PuerkitoBio/goquery"
	"log"
)

const proxybayURL string = "https://proxybay.github.io/"

// Mirror represents a PirateBay mirror and its status.
type Mirror struct {
	URL     string `json:"url"`
	Country string `json:"country"`
	Status  bool   `json:"status"`
}

// MirrorsCommand defines the mirrors command and holds its options.
type MirrorsCommand struct {
}

// Execute acts as the call back of the mirrors command.
func (m *MirrorsCommand) Execute(args []string) error {
	mirrors := GetMirrors()

	if Options.JSON {
		mirrorsJSON, _ := json.MarshalIndent(mirrors, "", "   ")
		log.Println(mirrorsJSON)
	}

	for _, mirror := range mirrors {
		status := "x"
		if !mirror.Status {
			status = " "
		}

		log.Printf("[%s] %s %s\n", status, mirror.Country, mirror.URL)
	}

	return nil
}

// GetMirrors retrieves a list of PirateBay mirrors.
func GetMirrors() []Mirror {

	doc, _ := shared.HTTPGet(proxybayURL)

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
