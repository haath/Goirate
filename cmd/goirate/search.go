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
	torrentSearchArgs
	Args positionalArgs `positional-args:"1" required:"1"`
}

// Execute is the callback of the mirrors command.
func (m *SearchCommand) Execute(args []string) error {

	if !m.ValidOutputFlags() {
		return errors.New("too many flags specifying the kind of output")
	}

	scraper, err := m.GetScraper()

	if err != nil {
		return err
	}

	torrents, err := (*scraper).Search(m.Args.Query)

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
