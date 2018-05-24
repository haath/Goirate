package main

import (
	"encoding/json"
	"git.gmantaos.com/haath/Gorrent/pkg/piratebay"
	"log"
)

// SearchCommand defines the search command and holds its options.
type SearchCommand struct {
	Args      searchArgs `positional-args:"1" required:"1"`
	Mirror    string     `short:"m" long:"mirror" description:"The PirateBay mirror URL to use. By default one is chosen at runtime."`
	SourceURL string     `short:"s" long:"source" description:"Link to the list of PirateBay proxies that will be used to pick a mirror."`
	Trusted   bool       `long:"trusted" description:"Only consider torrents where the uploader is either VIP or Trusted."`
}

type searchArgs struct {
	Query string `positional-arg-name:"query"`
}

// Execute acts as the call back of the mirrors command.
func (m *SearchCommand) Execute(args []string) error {

	var scraper piratebay.PirateBayScaper

	if m.SourceURL != "" {
		scraper = piratebay.NewScraper(m.SourceURL)
	} else {
		var mirrorScraper piratebay.MirrorScraper

		if m.SourceURL != "" {
			mirrorScraper.SetProxySourceURL(m.SourceURL)
		}

		mirror, err := mirrorScraper.PickMirror()

		if err != nil {
			return err
		}

		scraper = piratebay.NewScraper(mirror.URL)
	}

	torrents, err := scraper.Search(m.Args.Query)

	if err != nil {
		return err
	}

	if Options.JSON {
		torrentsJSON, err := json.MarshalIndent(torrents, "", "   ")

		if err != nil {
			return err
		}

		log.Println(string(torrentsJSON))
		return nil
	}

	return nil
}
