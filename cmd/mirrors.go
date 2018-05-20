package main

import (
	"encoding/json"
	"git.gmantaos.com/haath/Gorrent/pkg/piratebay"
	"log"
)

// MirrorsCommand defines the mirrors command and holds its options.
type MirrorsCommand struct {
	SourceURL string `short:"s" long:"source" description:"Link to list of PirateBay proxies. Default: proxybay.github.io"`
}

// Execute acts as the call back of the mirrors command.
func (m *MirrorsCommand) Execute(args []string) error {

	var scraper piratebay.MirrorScraper

	scraper.SetProxySourceURL(m.SourceURL)

	mirrors := scraper.GetMirrors()

	if Options.JSON {
		mirrorsJSON, err := json.MarshalIndent(mirrors, "", "   ")

		if err != nil {
			return err
		}

		log.Println(string(mirrorsJSON))
		return nil
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
