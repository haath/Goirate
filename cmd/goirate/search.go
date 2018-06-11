package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"git.gmantaos.com/haath/Goirate/pkg/torrents"
	"github.com/olekukonko/tablewriter"
	"log"
)

// SearchCommand defines the search command and holds its options.
type SearchCommand struct {
	torrents.SearchFilters
	Args       searchArgs `positional-args:"1" required:"1"`
	Mirror     string     `short:"m" long:"mirror" description:"The PirateBay mirror URL to use. By default one is chosen at runtime."`
	SourceURL  string     `short:"s" long:"source" description:"Link to the list of PirateBay proxies that will be used to pick a mirror."`
	MagnetLink bool       `long:"only-magnet" description:"Only output magnet links, one on each line."`
	TorrentURL bool       `long:"only-url" description:"Only output torrent urls, one on each line."`
	Count      uint       `short:"c" long:"count" description:"Limit the number of results."`
}

type searchArgs struct {
	Query string `positional-arg-name:"query"`
}

// Execute acts as the call back of the mirrors command.
func (m *SearchCommand) Execute(args []string) error {

	var scraper torrents.PirateBayScaper

	if !m.validOutputFlags() {
		return errors.New("too many flags specifying the kind of output")
	}

	if m.SourceURL != "" {
		scraper = torrents.NewScraper(m.SourceURL)
	} else {
		var mirrorScraper torrents.MirrorScraper

		if m.SourceURL != "" {
			mirrorScraper.SetProxySourceURL(m.SourceURL)
		}

		mirror, err := mirrorScraper.PickMirror()

		if err != nil {
			return err
		}

		scraper = torrents.NewScraper(mirror.URL)
	}

	torrents, err := scraper.Search(m.Args.Query)

	if err != nil {
		return err
	}

	torrents = m.filterTorrentList(torrents)

	if Options.JSON {
		torrentsJSON, err := json.MarshalIndent(torrents, "", "   ")

		if err != nil {
			return err
		}

		log.Println(string(torrentsJSON))

	} else if m.MagnetLink {

		for _, torrent := range torrents {
			log.Println(torrent.Magnet)
		}

	} else if m.TorrentURL {

		for _, torrent := range torrents {
			log.Println(torrent.FullURL())
		}

	} else {

		log.Printf(getTorrentsTable(torrents))

	}

	return nil
}

func (m *SearchCommand) filterTorrentList(torrentList []torrents.Torrent) []torrents.Torrent {

	var filtered []torrents.Torrent

	for _, torrent := range torrentList {

		if !m.VerifiedUploader || torrent.VerifiedUploader {
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
	if m.MagnetLink {
		outputFlags++
	}
	if m.TorrentURL {
		outputFlags++
	}

	return outputFlags <= 1
}

func getTorrentsTable(torrents []torrents.Torrent) string {
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
