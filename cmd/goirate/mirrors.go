package main

import (
	"bytes"
	"encoding/json"
	"git.gmantaos.com/haath/Goirate/pkg/torrents"
	"github.com/olekukonko/tablewriter"
	"log"
	"strings"
)

// MirrorsCommand defines the mirrors command and holds its options.
type MirrorsCommand struct {
	SourceURL string `short:"s" long:"source" description:"Link to a list of PirateBay proxies. Default: proxybay.github.io"`
}

// Execute acts as the call back of the mirrors command.
func (m *MirrorsCommand) Execute(args []string) error {

	var scraper torrents.MirrorScraper

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

	log.Printf(getMirrorsTable(mirrors))

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
