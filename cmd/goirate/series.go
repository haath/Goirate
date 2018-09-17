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
	"strconv"
	"strings"

	"git.gmantaos.com/haath/Goirate/pkg/utils"

	"git.gmantaos.com/haath/Goirate/pkg/torrents"

	"git.gmantaos.com/haath/Goirate/pkg/series"
	"github.com/BurntSushi/toml"
	"github.com/olekukonko/tablewriter"
)

// SeriesCommand is the command used to add or remove series from the watchlist
// as well as perform a scan for new episodes.
type SeriesCommand struct {
	Add    addCommand    `command:"add" description:"Add a series to the watchlist."`
	Remove removeCommand `command:"remove" alias:"rm" description:"Remove a series from the watchlist."`
	Show   showCommand   `command:"show" alias:"ls" description:"Print out the current series watchlist."`
	Scan   scanCommand   `command:"scan" description:"Perform a scan for new episodes on the existing watchlist."`
}

type addCommand struct {
	LastEpisode      string                `long:"last-episode" short:"e" description:"The last episode that came out."`
	MinQuality       torrents.VideoQuality `long:"min-quality" description:"The minimum video quality to accept when scanning for torrents of this series."`
	VerifiedUploader bool                  `long:"trusted" description:"Only accepted torrents from trusted or verified uploaders for this series."`
	Force            bool                  `long:"force" short:"f" description:"Overwrite this series if it already exists in the watchlist."`
	Args             struct {
		Title string `positional-arg-name:"<series title>"`
	} `positional-args:"1" required:"1"`
}
type removeCommand struct {
	Args struct {
		Title string `positional-arg-name:"<series title/id>"`
	} `positional-args:"1" required:"1"`
}
type showCommand struct{}
type scanCommand struct {
	torrentSearchArgs

	DryRun   bool `long:"dry-run" description:"Perform the scan for new episodes without downloading torrents, sending notifications or updating the episode numbers in the watchlist."`
	NoUpdate bool `long:"no-update" description:"Perform the scan for new episodes without updating the last episode aired in the watchlist."`
}

// Execute is the callback of the series add command.
func (cmd *addCommand) Execute(args []string) error {

	tvdbToken, err := tvdbLogin()

	if err != nil {
		return err
	}

	seriesID, seriesName, err := tvdbToken.Search(cmd.Args.Title)

	if err != nil {
		return err
	}

	var episode series.Episode

	if cmd.LastEpisode != "" {
		episode = series.ParseEpisodeString(cmd.LastEpisode)
	}

	if episode.Season == 0 && episode.Episode == 0 {
		episode, err = tvdbToken.LastEpisode(seriesID)

		if err != nil {
			return err
		}
	}

	ser := series.Series{
		ID:               seriesID,
		Title:            seriesName,
		MinQuality:       cmd.MinQuality,
		VerifiedUploader: cmd.VerifiedUploader,
		LastEpisode:      episode,
	}

	seriesList := loadSeries()

	if containsID(seriesList, ser.ID) {

		if cmd.Force {

			seriesList = remove(seriesList, ser.ID, "")

		} else {

			return fmt.Errorf("series %v already on the watchlist", ser.Title)

		}
	}

	seriesList = append(seriesList, ser)

	storeSeries(seriesList)

	log.Printf("Added series: %v", ser.Title)

	return nil
}

// Execute is the callback of the series remove command.
func (cmd *removeCommand) Execute(args []string) error {

	cmd.Args.Title = utils.NormalizeQuery(cmd.Args.Title)
	id, _ := strconv.Atoi(cmd.Args.Title)

	seriesList := loadSeries()

	seriesList = remove(seriesList, id, cmd.Args.Title)

	storeSeries(seriesList)

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

	tvdbToken, err := tvdbLogin()

	if err != nil {
		return err
	}

	var torrentList []interface{}

	seriesList := loadSeries()

	for i := range seriesList {

		ser := &seriesList[i]

		found := true

		for found && (cmd.Count == 0 || uint(len(torrentList)) < cmd.Count) {

			found, err = cmd.scanSeries(tvdbToken, ser, &torrentList)

			if err != nil {
				return err
			}

			if found && !cmd.DryRun && !cmd.NoUpdate {

				storeSeries(seriesList)

			}
		}

		if cmd.Count > 0 && uint(len(torrentList)) == cmd.Count {
			break
		}

	}

	if Options.JSON {

		torrentsJSON, err := json.MarshalIndent(torrentList, "", "   ")

		if err != nil {
			return err
		}

		log.Println(string(torrentsJSON))

	}

	return nil
}

func (cmd *scanCommand) scanSeries(tvdbToken *series.TVDBToken, ser *series.Series, torrentList *[]interface{}) (bool, error) {
	filters := cmd.GetFilters()
	filters.MinQuality = ser.MinQuality
	filters.VerifiedUploader = ser.VerifiedUploader

	nextEpisode, err := ser.NextEpisode(tvdbToken)

	if err != nil {
		return false, err
	}

	scraper, err := cmd.GetScraper(ser.SearchQuery(nextEpisode))

	if err != nil {
		return false, err
	}

	torrent, err := ser.GetTorrent(scraper, filters, nextEpisode)

	if err != nil {
		return false, err
	}

	if torrent == nil {
		return false, nil
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

	*torrentList = append(*torrentList, struct {
		Series  series.Series    `json:"series"`
		Torrent torrents.Torrent `json:"torrent"`
	}{
		Series:  *ser,
		Torrent: *torrent,
	})

	ser.LastEpisode = nextEpisode

	return true, nil
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

func tvdbLogin() (*series.TVDBToken, error) {

	cred := series.EnvTVDBCredentials()

	if cred.APIKey == "" || cred.UserKey == "" || cred.Username == "" {
		cred = Config.TVDBCredentials
	}

	if cred.APIKey == "" || cred.UserKey == "" || cred.Username == "" {

		return nil, fmt.Errorf("the series command requires valid credentials for the TVDB API to be configured at %v, which can be obtained by making a free account at https://www.thetvdb.com/", configPath())
	}

	tkn, err := cred.Login()

	if err != nil {
		return nil, errors.New("Error logging into the TVDB API.: " + err.Error())
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

func remove(seriesList []series.Series, id int, title string) []series.Series {

	var newList []series.Series

	for _, ser := range seriesList {

		normalizedTitle := utils.NormalizeQuery(ser.Title)

		if (id == 0 && title != "" && strings.Contains(normalizedTitle, title)) || ser.ID == id {

			log.Printf("Removed series: %v", ser.Title)
			continue
		}

		newList = append(newList, ser)
	}
	return newList
}
