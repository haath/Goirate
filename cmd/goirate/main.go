package main

import (
	"log"
	"os"

	"git.gmantaos.com/haath/Goirate/pkg/torrents"
	"github.com/jessevdk/go-flags"
)

// Options holds the command line options for the cli program
var Options struct {
	// Options
	Verbose bool   `short:"v" long:"verbose" description:"Show more information."`
	JSON    bool   `short:"j" long:"json" description:"Output in JSON format."`
	Version func() `long:"version" description:"Show the current version."`

	// Commands
	Mirrors MirrorsCommand `command:"mirrors" description:"Get a list of PirateBay mirrors."`
	Search  SearchCommand  `command:"search" alias:"s" description:"Search for torrents."`
	Movie   MovieCommand   `command:"movie" alias:"m" description:"Scrape a movie and find torrents for it."`
	IMDb    IMDbCommand    `command:"imdb" description:"Search IMDb for movies to retrieve their IMDbID."`
}

type torrentSearchArgs struct {
	Mirror     string `short:"m" long:"mirror" description:"The PirateBay mirror URL to use. By default one is chosen at runtime."`
	SourceURL  string `short:"s" long:"source" description:"Link to the list of PirateBay proxies that will be used to pick a mirror."`
	Count      uint   `short:"c" long:"count" description:"Limit the number of results."`
	MagnetLink bool   `long:"only-magnet" description:"Only output magnet links, one on each line."`
	TorrentURL bool   `long:"only-url" description:"Only output torrent urls, one on each line."`
}

type positionalArgs struct {
	Query string `positional-arg-name:"query"`
}

func (a torrentSearchArgs) GetScraper(query string) (*torrents.PirateBayScaper, error) {

	var scraper torrents.PirateBayScaper

	if a.Mirror != "" {

		scraper = torrents.NewScraper(a.Mirror)

	} else {

		var mirrorScraper torrents.MirrorScraper

		if a.SourceURL != "" {
			mirrorScraper.SetProxySourceURL(a.SourceURL)
		}

		mirror, err := mirrorScraper.PickMirror(query)

		if err != nil {
			return nil, err
		}

		scraper = torrents.NewScraper(mirror.URL)
	}

	return &scraper, nil
}

func (a *torrentSearchArgs) ValidOutputFlags() bool {
	outputFlags := 0

	if Options.JSON {
		outputFlags++
	}
	if a.MagnetLink {
		outputFlags++
	}
	if a.TorrentURL {
		outputFlags++
	}

	return outputFlags <= 1
}

func main() {

	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	parser := flags.NewParser(&Options, flags.HelpFlag|flags.PassDoubleDash|flags.PrintErrors)

	Options.Version = func() {
		log.Printf("Goirate build: %v\n", VERSION)
		os.Exit(0)
	}

	parser.Parse()
}
