package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"goirate/series"
	"goirate/torrents"
	"goirate/utils"

	"github.com/BurntSushi/toml"
	"github.com/olekukonko/tablewriter"
)

// SeriesCommand is the command used to add or remove series from the watchlist
// as well as perform a scan for new episodes.
type SeriesCommand struct {
	Add    addCommand    `command:"add" description:"Add a series to the watchlist."`
	Remove removeCommand `command:"remove" alias:"rm" description:"Remove a series from the watchlist."`
	Search searchCommand `command:"search" alias:"s" description:"Search for a series given a query string."`
	Show   showCommand   `command:"show" alias:"ls" description:"Print out the current series watchlist."`
	Scan   scanCommand   `command:"scan" description:"Perform a scan for new episodes on the existing watchlist."`
}

type addCommand struct {
	LastEpisode      string                `long:"last-episode" short:"e" description:"The last episode that came out."`
	MinQuality       torrents.VideoQuality `long:"min-quality" description:"The minimum video quality to accept when scanning for torrents of this series."`
	VerifiedUploader bool                  `long:"trusted" description:"Only accepted torrents from trusted or verified uploaders for this series."`
	Force            bool                  `long:"force" short:"f" description:"Overwrite this series if it already exists in the watchlist."`
	Show             bool                  `long:"ls" description:"Execute the show command after adding."`
	Args             struct {
		Title string `positional-arg-name:"<title | imdbID>"`
	} `positional-args:"1" required:"1"`
}
type removeCommand struct {
	Show bool `long:"ls" description:"Execute the show command after removing."`
	Args struct {
		Title string `positional-arg-name:"<title | id>"`
	} `positional-args:"1" required:"1"`
}
type searchCommand struct {
	Count int `long:"count" short:"c" description:"Limit the number of results."`
	Args  struct {
		Query string `positional-arg-name:"<title | imdb | tvmaze id>"`
	} `positional-args:"1" required:"1"`
}
type showCommand struct{}
type scanCommand struct {
	torrentSearchArgs

	DryRun   bool `long:"dry-run" description:"Perform the scan for new episodes without downloading torrents, sending notifications or updating the episode numbers in the watchlist."`
	NoUpdate bool `long:"no-update" description:"Perform the scan for new episodes without updating the last episode aired in the watchlist."`
	Quiet    bool `long:"quiet" short:"q" description:"Do not print anything to the standard output."`
	Quick    bool `long:"quick" description:"Perform a quick scan, only searching for torrents for episodes that were found on the tvmaze API."`
}

type seriesTorrent struct {
	Episode series.Episode   `json:"episode"`
	Torrent torrents.Torrent `json:"torrent"`
}
type seriesTorrents struct {
	Series   *series.Series  `json:"series"`
	Torrents []seriesTorrent `json:"torrents"`
}

// Execute is the callback of the series add command.
func (cmd *addCommand) Execute(args []string) error {

	tvmazeToken, err := tvmazeLogin()

	if err != nil {
		return err
	}

	show, err := tvmazeToken.SearchFirst(cmd.Args.Title)

	if err != nil {
		return err
	}

	var episode series.Episode

	if cmd.LastEpisode != "" {

		episode = series.ParseEpisodeString(cmd.LastEpisode)

		if episode.Season == 0 && episode.Episode == 0 {

			return fmt.Errorf("unable to parse the last episode number from: %v", cmd.LastEpisode)
		}

	} else {

		episode, err = tvmazeToken.LastEpisode(show.ID)

		if err != nil {
			return err
		}
	}

	ser := series.Series{
		ID:               show.ID,
		Title:            show.Name,
		MinQuality:       cmd.MinQuality,
		VerifiedUploader: cmd.VerifiedUploader,
		LastEpisode:      episode,
	}
	ser.Actions.Emails = []string{}

	seriesList := loadSeries()

	if containsID(seriesList, ser.ID) {

		if cmd.Force {

			remove(&seriesList, ser.ID, "")

		} else {

			return fmt.Errorf("series %v already on the watchlist", ser.Title)

		}
	}

	seriesList = append(seriesList, ser)

	storeSeries(seriesList)

	if cmd.Show {
		var showCmd showCommand
		return showCmd.Execute(args)
	}

	return nil
}

// Execute is the callback of the series remove command.
func (cmd *removeCommand) Execute(args []string) error {

	cmd.Args.Title = utils.NormalizeQuery(cmd.Args.Title)
	id, _ := strconv.Atoi(cmd.Args.Title)

	seriesList := loadSeries()

	if !remove(&seriesList, id, cmd.Args.Title) {

		return fmt.Errorf("no series found on the watchlist matching: %v\nhint: goirate series show", cmd.Args.Title)
	}

	storeSeries(seriesList)

	if cmd.Show {
		var showCmd showCommand
		return showCmd.Execute(args)
	}

	return nil
}

// Execute is the callback of the series search command.
func (cmd *searchCommand) Execute(args []string) error {

	tvmazeToken, err := tvmazeLogin()
	if err != nil {
		return err
	}

	searchResult, err := tvmazeToken.Search(cmd.Args.Query)
	if err != nil {
		return err
	}

	if Options.JSON {

		seriesJSON, err := json.MarshalIndent(searchResult, "", "   ")

		if err != nil {
			return err
		}

		log.Println(string(seriesJSON))

	} else {

		log.Println(getSeriesSearchTable(searchResult, cmd.Count))
	}

	return nil
}

// Execute is the callback of the series show command.
func (cmd *showCommand) Execute(args []string) error {

	seriesList := loadSeries()

	if Options.JSON {

		seriesJSON, err := json.MarshalIndent(seriesList, "", "   ")

		if err != nil {
			return err
		}

		log.Println(string(seriesJSON))

	} else {

		log.Print(getSeriesTable(seriesList))

	}

	return nil
}

// Execute is the callback of the series scan command.
func (cmd *scanCommand) Execute(args []string) error {

	if cmd.Quiet {
		disableOutput()
	}

	tvmazeToken, err := tvmazeLogin()

	if err != nil {
		return err
	}

	var torrentList []seriesTorrents

	seriesList := loadSeries()

	for i := range seriesList {

		ser := &seriesList[i]

		found := true

		for found && (cmd.Count == 0 || seriesTorrentCount(torrentList) < cmd.Count) {

			found, err = cmd.scanSeries(tvmazeToken, ser, &torrentList)

			if err != nil && os.Getenv("GOIRATE_DEBUG") == "true" {
				log.Println(err)
			}
		}
	}

	if !cmd.DryRun && !cmd.NoUpdate {

		storeSeries(seriesList)
	}

	if Options.JSON {

		torrentsJSON, err := json.MarshalIndent(torrentList, "", "   ")

		if err != nil {
			return err
		}

		log.Println(string(torrentsJSON))
	}

	if !cmd.DryRun {
		err = cmd.handleSeriesTorrents(torrentList)

		if err != nil {
			return err
		}
	}

	if cmd.Quiet {
		enableOutput()
	}

	return nil
}

func (cmd *scanCommand) scanSeries(tvmazeToken *series.TVmazeToken, ser *series.Series, torrentList *[]seriesTorrents) (bool, error) {

	filters := cmd.GetFilters()

	if ser.MinQuality != "" {
		filters.MinQuality = ser.MinQuality
	}
	filters.VerifiedUploader = filters.VerifiedUploader || ser.VerifiedUploader

	nextEpisode, err := ser.NextEpisode(tvmazeToken)

	if err != nil {

		return false, err
	}

	if cmd.Quick && !nextEpisode.HasAired() {

		return false, nil
	}

	if !cmd.MagnetLink && !Options.JSON && !cmd.TorrentURL {

		log.Printf("Searching for: %s %s\n", ser.Title, nextEpisode)
	}

	allTorrents, err := ser.GetTorrents(*filters, nextEpisode)

	if err != nil {

		return false, err
	}

	torrent, err := torrents.PickVideoTorrent(allTorrents, *filters)

	if err != nil || torrent == nil {

		return false, err
	}

	if cmd.MagnetLink {

		log.Println(torrent.Magnet)

	} else if Options.JSON {

		// Do nothing, the torrentList will be updated and printed by the calling func

	} else if cmd.TorrentURL {

		log.Println(torrent.FullURL())

	} else {

		log.Printf("Torrent found for: %s %s\n%s\n%s\n\n", ser.Title, nextEpisode, torrent.FullURL(), torrent.Magnet)
	}

	appendSeriesTorrent(torrentList, ser, nextEpisode, *torrent)

	ser.LastEpisode = nextEpisode

	return true, err
}

func (cmd *scanCommand) handleSeriesTorrents(seriesTorrentsList []seriesTorrents) error {

	for _, seriesTorrents := range seriesTorrentsList {

		/*
			Send e-mails, groupins episode torrents per series
		*/
		if Config.Watchlist.SendEmail.OverridenBy(seriesTorrents.Series.Actions.SendEmail) {

			// Send e-mail

			notify := Config.Watchlist.Emails

			if len(seriesTorrents.Series.Actions.Emails) > 0 {

				notify = seriesTorrents.Series.Actions.Emails
			}

			if len(notify) == 0 {

				return fmt.Errorf("sending e-mails is enabled, but no recipients are specified")
			}

			log.Printf("Sending e-mail to: %s\n", notify)

			body, err := LoadSeriesTemplate(seriesTorrents)

			if err != nil {
				return err
			}

			var subject string

			if len(seriesTorrents.Torrents) > 1 {

				subject = fmt.Sprintf("Episodes out for %s (%s)", seriesTorrents.Series.Title, episodeRangeString(seriesTorrents))

			} else {

				subject = fmt.Sprintf("Episode out for %s (%s)", seriesTorrents.Series.Title, seriesTorrents.Torrents[0].Episode)
			}

			err = Config.SMTPConfig.SendEmail(subject, body, notify...)

			if err != nil {
				return err
			}
		}

		/*
			Loop over individual torrents to send each of them to qBittorrent for download
		*/
		for _, seriesTorrent := range seriesTorrents.Torrents {

			if Config.Watchlist.Download.OverridenBy(seriesTorrents.Series.Actions.Download) {

				// Send the torrent to the qBittorrent daemon for download

				qbt, err := Config.QBittorrentConfig.GetClient()

				if err != nil {
					return err
				}

				downloadPath := Config.DownloadDir.Series

				if Config.KodiMediaPaths {

					downloadPath = path.Join(
						downloadPath,
						seriesTorrents.Series.Title,
						fmt.Sprintf("Season %d", seriesTorrent.Episode.Season),
					)
				}

				log.Printf("Downloading: %s %s (%s)\n", seriesTorrents.Series.Title, seriesTorrent.Episode, downloadPath)

				err = qbt.AddTorrent(seriesTorrent.Torrent.Magnet, downloadPath)

				if err != nil {
					return err
				}

			}

		}

	}

	return nil
}

func appendSeriesTorrent(torrentList *[]seriesTorrents, ser *series.Series, episode series.Episode, torrent torrents.Torrent) {

	serTorrent := seriesTorrent{Torrent: torrent, Episode: episode}

	for i := range *torrentList {

		item := (*torrentList)[i]

		if item.Series.ID == ser.ID {

			(*torrentList)[i].Torrents = append(item.Torrents, serTorrent)

			return
		}
	}

	*torrentList = append(*torrentList, seriesTorrents{
		Series:   ser,
		Torrents: []seriesTorrent{serTorrent},
	})
}

func seriesTorrentCount(torrentList []seriesTorrents) uint {

	var count uint
	for _, seriesTorrents := range torrentList {

		count += uint(len(seriesTorrents.Torrents))
	}
	return count
}

func episodeRangeString(seriesTorrents seriesTorrents) string {

	min := series.Episode{Season: ^uint(0), Episode: ^uint(0)}
	max := series.Episode{Season: 0, Episode: 0}

	if len(seriesTorrents.Torrents) == 0 {

		return ""
	}

	for _, seriesTorrent := range seriesTorrents.Torrents {

		if min.IsAfter(seriesTorrent.Episode) {

			min = seriesTorrent.Episode
		}

		if seriesTorrent.Episode.IsAfter(max) {

			max = seriesTorrent.Episode
		}
	}

	if min.Season == max.Season && min.Episode == max.Episode {

		return min.String()
	}

	return fmt.Sprintf("%s - %s", min, max)
}

func loadSeries() []series.Series {

	var seriesList struct {
		Series []series.Series `toml:"series"`
	}

	if _, err := os.Stat(seriesConfigPath()); err == nil {

		tomlBytes, err := ioutil.ReadFile(seriesConfigPath())

		if err != nil {
			log.Fatal(err)
		}

		tomlString := string(tomlBytes)

		if _, err := toml.Decode(tomlString, &seriesList); err != nil {
			log.Fatal(err)
		}

	}

	sort.Slice(seriesList.Series, func(i, j int) bool {
		return seriesList.Series[i].Title < seriesList.Series[j].Title
	})

	return seriesList.Series
}

func storeSeries(seriesList []series.Series) {

	file, err := os.OpenFile(seriesConfigPath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	encoder := toml.NewEncoder(file)

	encoder.Encode(struct {
		Series []series.Series `toml:"series"`
	}{seriesList})
}

func seriesConfigPath() string {

	return path.Join(configDir(), "series.toml")
}

func capitalize(str string) string {

	// Function replacing words (assuming lower case input)
	replace := func(word string) string {
		switch word {
		case "with", "in", "a", "to", "of":
			return word
		}
		return strings.Title(word)
	}

	r := regexp.MustCompile(`\w+`)
	str = r.ReplaceAllStringFunc(strings.ToLower(str), replace)

	return str
}

func getSeriesTable(seriesList []series.Series) string {
	buf := bytes.NewBufferString("")

	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"ID", "Series", "Season", "Last Episode", "Min. Quality"})
	table.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoFormatHeaders(false)

	for _, series := range seriesList {

		table.Append([]string{strconv.Itoa(series.ID), series.Title,
			fmt.Sprint(series.LastEpisode.Season), fmt.Sprint(series.LastEpisode.Episode), string(series.MinQuality)})
	}

	table.Render()

	return buf.String()
}

func tvmazeLogin() (*series.TVmazeToken, error) {

	cred := series.EnvTVmazeCredentials()

	tkn, err := cred.Login()

	if err != nil {
		return nil, errors.New("Error logging into the TVmaze API.: " + err.Error())
	}

	return &tkn, nil
}

func containsID(seriesList []series.Series, id int) bool {

	for _, ser := range seriesList {

		if ser.ID == id {
			return true
		}
	}

	return false
}

func remove(seriesList *[]series.Series, id int, title string) bool {

	var newList []series.Series
	var found bool

	for _, ser := range *seriesList {

		normalizedTitle := utils.NormalizeQuery(ser.Title)

		if (id == 0 && title != "" && strings.Contains(normalizedTitle, title)) || ser.ID == id {

			found = true
			continue
		}

		newList = append(newList, ser)
	}
	*seriesList = newList
	return found
}

func getSeriesSearchTable(searchResult []series.TVmazeSeries, count int) string {
	buf := bytes.NewBufferString("")

	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"ID", "Title", "Premiered"})
	table.SetColumnAlignment([]int{tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_CENTER})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoFormatHeaders(false)

	for i, show := range searchResult {

		if count > 0 && i >= count {
			break
		}

		var premierYear string
		premierDate, err := time.Parse("2006-01-02", show.Premiered)
		if err == nil {
			premierYear = strconv.Itoa(premierDate.Year())
		}

		table.Append([]string{fmt.Sprint(show.ID), show.Name, premierYear})
	}

	table.Render()

	return buf.String()
}

func disableOutput() {
	var buf bytes.Buffer
	log.SetOutput(&buf)
}

func enableOutput() {
	log.SetOutput(os.Stdout)
}
