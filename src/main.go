package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"os/user"
	"path"

	"goirate/torrents"
	"goirate/utils"

	"github.com/jessevdk/go-flags"
)

// Options holds the command line options for the cli program
var Options struct {
	// Options
	JSON    bool   `short:"j" long:"json" description:"Output in JSON format."`
	Version func() `long:"version" description:"Print the current version and exit."`

	// Commands
	Config      ConfigCommand      `command:"config" description:"Update the configuration files."`
	Mirrors     MirrorsCommand     `command:"mirrors" description:"Get a list of PirateBay mirrors."`
	Search      SearchCommand      `command:"search" description:"Search for torrents."`
	Series      SeriesCommand      `command:"series" alias:"s" description:"Manage the series watchlist or perform a scan."`
	Movie       MovieCommand       `command:"movie" alias:"m" description:"Scrape a movie and find torrents for it."`
	MovieSearch MovieSearchCommand `command:"movie-search" description:"Search OMDB for movies to retrieve their IMDbID and release year."`
	Update      UpdateCommand      `command:"update" alias:"u" description:"Update the tool to the latest version."`
}

type torrentSearchArgs struct {
	torrents.SearchFilters
	Mirror     string `short:"m" long:"mirror" description:"The PirateBay mirror URL to use. By default one is chosen at runtime."`
	SourceURL  string `short:"s" long:"source" description:"Link to the list of PirateBay proxies that will be used to pick a mirror."`
	Count      uint   `short:"c" long:"count" description:"Limit the number of results."`
	MagnetLink bool   `long:"only-magnet" description:"Only output magnet links, one on each line."`
	TorrentURL bool   `long:"only-url" description:"Only output torrent urls, one on each line."`
}

type positionalArgs struct {
	Query string `positional-arg-name:"query"`
}

func (a torrentSearchArgs) GetTorrents(query string) ([]torrents.Torrent, error) {

	if !a.ValidOutputFlags() {
		return nil, errors.New("too many flags specifying the kind of output")
	}

	if a.Mirror != "" {

		scraper := torrents.NewScraper(a.Mirror)
		return scraper.Search(query)
	}

	mirrorScraper := GetMirrorScraper()

	if a.SourceURL != "" {
		mirrorScraper.SetProxySourceURL(a.SourceURL)
	}

	return mirrorScraper.GetTorrents(query)
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

func (a torrentSearchArgs) GetFilters() *torrents.SearchFilters {

	ApplyConfig(&a.SearchFilters)

	if a.Mirror != "" {
		a.SearchFilters.MirrorURL = a.Mirror
	}
	if a.SourceURL != "" {
		a.SearchFilters.ProxyListURL = a.SourceURL
	}

	return &a.SearchFilters
}

func configDir() string {

	usr, usrErr := user.Current()

	var dir string

	if flag.Lookup("test.v") != nil {

		dir = path.Join(usr.HomeDir, ".goirate.test")

	} else if os.Getenv("GOIRATE_DIR") != "" {

		dir = os.Getenv("GOIRATE_DIR")

	} else if usrErr != nil {

		// Being here usually means we can't produce the current user.
		// In this case the '~' path will most likely also not be around.
		// With the crontab in mind, we'll default this case to a directory
		// in the current working directory of shell. Which in the case of cron,
		// should be the user's home folder anyway.
		dir = ".goirate"

	} else {

		dir = path.Join(usr.HomeDir, ".goirate")
	}

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	return dir
}

func main() {

	log.SetFlags(0)
	log.SetOutput(&utils.GoirateLogger{})

	ImportConfig()

	parser := flags.NewParser(&Options, flags.HelpFlag|flags.PassDoubleDash|flags.PrintErrors)

	Options.Version = func() {
		log.Printf("Goirate build: %v\n", Version)
		os.Exit(0)
	}

	_, err := parser.Parse()

	if err != nil {
		os.Exit(1)
	}
}
