package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"git.gmantaos.com/haath/Goirate/pkg/piratebay"
	"github.com/olekukonko/tablewriter"
	"log"
)

// SearchCommand defines the search command and holds its options.
type SearchCommand struct {
	Args        searchArgs `positional-args:"1" required:"1"`
	Mirror      string     `short:"m" long:"mirror" description:"The PirateBay mirror URL to use. By default one is chosen at runtime."`
	SourceURL   string     `short:"s" long:"source" description:"Link to the list of PirateBay proxies that will be used to pick a mirror."`
	Trusted     bool       `long:"trusted" description:"Only consider torrents where the uploader is either VIP or Trusted."`
	MagnetLinks bool       `long:"magnet" description:"Only output magnet links, one on each line."`
	TorrentURLs bool       `long:"urls" description:"Only output torrent urls, one on each line."`
	Count       uint       `short:"c" long:"count" description:"Limit the number of results."`
}

type searchArgs struct {
	Query string `positional-arg-name:"query"`
}

// Execute acts as the call back of the mirrors command.
func (m *SearchCommand) Execute(args []string) error {

	var scraper piratebay.PirateBayScaper

	if !m.validOutputFlags() {
		return errors.New("too many flags specifying the kind of output")
	}

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

	torrents = m.filterTorrentList(torrents)

	if err != nil {
		return err
	}

	if Options.JSON {
		torrentsJSON, err := json.MarshalIndent(torrents, "", "   ")

		if err != nil {
			return err
		}

		log.Println(string(torrentsJSON))

	} else if m.MagnetLinks {

		for _, torrent := range torrents {
			log.Println(torrent.Magnet)
		}

	} else if m.TorrentURLs {

		for _, torrent := range torrents {
			log.Println(torrent.FullURL())
		}

	} else {

		log.Printf(getTorrentsTable(torrents))

	}

	return nil
}

func (m *SearchCommand) filterTorrentList(torrents []piratebay.Torrent) []piratebay.Torrent {

	var filtered []piratebay.Torrent

	for _, torrent := range torrents {

		if !m.Trusted || torrent.VerifiedUploader {
			filtered = append(filtered, torrent)
		}

		if m.Count > 0 && uint(len(filtered)) >= m.Count {
			break
		}
	}

	return filtered
}

func (m *SearchCommand) validOutputFlags() bool {
	outputFlags := 0

	if Options.JSON {
		outputFlags++
	}
	if m.MagnetLinks {
		outputFlags++
	}
	if m.TorrentURLs {
		outputFlags++
	}

	return outputFlags <= 1
}

func getTorrentsTable(torrents []piratebay.Torrent) string {
	buf := bytes.NewBufferString("")

	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"Title", "Size", "Seeds/Peers"})
	table.SetColumnAlignment([]int{tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_DEFAULT})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: false, Bottom: false})
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowLine(true)
	table.SetAutoFormatHeaders(false)
	table.SetAutoWrapText(false)

	for _, torrent := range torrents {

		table.Append([]string{torrent.Title + "\n" + torrent.FullURL(), torrent.SizeString(), torrent.PeersString()})
	}

	table.Render()

	return buf.String()
}
