package main

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gitlab.com/haath/goirate/pkg/torrents"
)

// MirrorsCommand defines the mirrors command and holds its options.
type MirrorsCommand struct {
	SourceURL string `short:"s" long:"source" description:"Link to a list of PirateBay proxies. Default: proxybay.github.io"`
}

// Execute is the callback of the mirrors command.
func (m *MirrorsCommand) Execute(args []string) error {

	var scraper torrents.MirrorScraper

	scraper.SetProxySourceURL(m.SourceURL)

	mirrors, err := scraper.GetMirrors()

	if err != nil {
		return err
	}

	if Options.JSON {
		mirrorsJSON, err := json.MarshalIndent(mirrors, "", "   ")

		if err != nil {
			return err
		}

		log.Println(string(mirrorsJSON))
	} else {
		log.Printf(getMirrorsTable(mirrors))
	}

	return nil
}

func getMirrorsTable(mirrors []torrents.Mirror) string {
	buf := bytes.NewBufferString("")

	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{" ", "Country", "URL"})
	table.SetColumnAlignment([]int{tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_DEFAULT})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoFormatHeaders(false)

	for _, mirror := range mirrors {
		status := "x"
		if !mirror.Status {
			status = " "
		}

		table.Append([]string{status, strings.ToUpper(mirror.Country), mirror.URL})
	}

	table.Render()

	return buf.String()
}
