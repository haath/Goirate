package series

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"goirate/movies"
	"goirate/utils"
)

// TVDBCredentials holds the credentials used to authenticate
// with the TVDB API.
type TVDBCredentials struct {
	APIKey   string `toml:"api_key" json:"apikey"`
	UserKey  string `toml:"user_key" json:"userkey"`
	Username string `toml:"username" json:"username"`
}

// TVDBToken holds a token which represents an authenticated session
// with the TVDB API.
type TVDBToken struct {
	Token string `json:"token"`
}

type apiEndpoint string

const (
	baseEndpoint     apiEndpoint = "https://api.thetvdb.com"
	loginEndpoint    apiEndpoint = baseEndpoint + "/login"
	searchEndpoint   apiEndpoint = baseEndpoint + "/search/series"
	episodesEndpoint apiEndpoint = baseEndpoint + "/series/%v/episodes"
)

func (ep apiEndpoint) String() string {
	return string(ep)
}

// Login submits a POST request to the /login API endpoint, which is
// used to obtain a JWT token for authenticating with the rest of the API.
func (cred *TVDBCredentials) Login() (TVDBToken, error) {

	var tkn TVDBToken

	err := utils.HTTPPost(loginEndpoint.String(), cred, &tkn)

	return tkn, err
}

// Search will search the TVDB for the given series name and return its ID.
func (tkn *TVDBToken) Search(searchName string) (id int, name string, err error) {

	var searchResponse struct {
		Data []struct {
			ID         int    `json:"id"`
			SeriesName string `json:"seriesName"`
		} `json:"data"`
	}

	var searchURL string

	if movies.IsIMDbURL(searchName) {

		imdbID, err := movies.ExtractIMDbID(searchName)

		if err != nil {

			return 0, "", err
		}

		searchURL = fmt.Sprintf("%v?imdbId=%v", searchEndpoint, imdbID)

	} else if movies.IsIMDbID(searchName) {

		searchName, err = movies.FormatIMDbID(searchName)

		if err != nil {

			return 0, "", err
		}

		searchURL = fmt.Sprintf("%v?imdbId=%v", searchEndpoint, searchName)
	} else {

		searchURL = fmt.Sprintf("%v?name=%v", searchEndpoint, url.QueryEscape(searchName))
	}

	err = tkn.apiCall(searchURL, &searchResponse)

	if err == nil && len(searchResponse.Data) > 0 {
		name = searchResponse.Data[0].SeriesName
		id = searchResponse.Data[0].ID
	}

	return
}

// LastEpisode uses the TVDB API to retrieve the last episode that aired
// for a particular series.
func (tkn *TVDBToken) LastEpisode(seriesID int) (Episode, error) {

	var episode Episode

	callback := func(ep Episode) {

		if !ep.HasAired() {
			return
		}

		if ep.IsAfter(episode) {
			episode = ep
		}

	}

	err := tkn.getEpisodes(seriesID, callback)

	return episode, err
}

// NextEpisode uses the TVDB API to make a best guess as to which episode is sequentially
// next to the one given.
func (tkn *TVDBToken) NextEpisode(seriesID int, episode Episode) (nextEpisode Episode, err error) {

	nextSeasonOut := false
	nextSeasonFirst := Episode{Season: episode.Season + 1, Episode: 1}

	seasonHasMore := false
	curSeasonNext := Episode{Season: episode.Season, Episode: episode.Episode + 1}

	callback := func(ep Episode) {

		if ep.Season == episode.Season+1 {

			nextSeasonOut = true

			if ep.Episode == 1 {

				nextSeasonFirst = ep
			}
		}
		if ep.Season == episode.Season && ep.Episode == episode.Episode+1 {

			seasonHasMore = true
			curSeasonNext = ep
		}
	}

	err = tkn.getEpisodes(seriesID, callback)

	if nextSeasonOut && !seasonHasMore {

		nextEpisode = nextSeasonFirst

	} else {

		nextEpisode = curSeasonNext
	}

	return
}

func (tkn *TVDBToken) getEpisodes(seriesID int, callback func(Episode)) error {

	var episodeSearchResponse struct {
		Data []struct {
			Season      uint   `json:"airedSeason"`
			Episode     uint   `json:"airedEpisodeNumber"`
			EpisodeName string `json:"episodeName"`
			FirstAired  string `json:"firstAired"`
		} `json:"data"`
		Links struct {
			Next int `json:"next"`
			Last int `json:"last"`
		}
	}

	baseURL := fmt.Sprintf(episodesEndpoint.String(), seriesID)

	pageNum := 1
	lastPage := 1

	for pageNum <= lastPage {

		url := fmt.Sprintf("%v?page=%v", baseURL, pageNum)

		err := tkn.apiCall(url, &episodeSearchResponse)

		if err != nil {
			return err
		}

		for _, ep := range episodeSearchResponse.Data {

			var aired *time.Time

			if ep.FirstAired != "" {

				tmp, err := time.Parse("2006-01-02", ep.FirstAired)

				if err != nil {
					return err
				}

				aired = &tmp
			}

			episode := Episode{
				Season:  ep.Season,
				Episode: ep.Episode,
				Title:   ep.EpisodeName,
				Aired:   aired,
			}

			callback(episode)
		}

		pageNum++
		lastPage = episodeSearchResponse.Links.Last
	}

	return nil
}

func (tkn *TVDBToken) apiCall(url string, v interface{}) error {

	httpClient := utils.HTTPClient{AuthToken: tkn.Token}

	return httpClient.GetJSON(url, &v)
}

// EnvTVDBCredentials returns a TVDBCredentials struct variable which contains
// the credentials found in the TVDB_API_KEY, TVDB_USER_KEY and TVDB_USERNAME
// environment variables.
func EnvTVDBCredentials() TVDBCredentials {
	return TVDBCredentials{
		APIKey:   os.Getenv("TVDB_API_KEY"),
		UserKey:  os.Getenv("TVDB_USER_KEY"),
		Username: os.Getenv("TVDB_USERNAME"),
	}
}
