package series

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"git.gmantaos.com/haath/Goirate/pkg/utils"
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

	searchURL := fmt.Sprintf("%v?name=%v", searchEndpoint, url.QueryEscape(searchName))

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

	callback := func(ep Episode, aired *time.Time) {

		if aired == nil || time.Now().Sub(*aired).Hours() < 1 {
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
func (tkn *TVDBToken) NextEpisode(seriesID int, episode Episode) (Episode, error) {

	nextSeasonOut := false
	seasonHasMore := false

	callback := func(ep Episode, aired *time.Time) {
		if ep.Season > episode.Season {
			nextSeasonOut = true
		}
		if ep.Season == episode.Season && ep.Episode > episode.Episode {
			seasonHasMore = true
		}
	}

	err := tkn.getEpisodes(seriesID, callback)

	if nextSeasonOut && !seasonHasMore {
		episode.Season++
		episode.Episode = 1
	} else {
		episode.Episode++
	}

	return episode, err
}

func (tkn *TVDBToken) getEpisodes(seriesID int, callback func(Episode, *time.Time)) error {

	var episodeSearchResponse struct {
		Data []struct {
			Season     uint   `json:"airedSeason"`
			Episode    uint   `json:"airedEpisodeNumber"`
			FirstAired string `json:"firstAired"`
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

			episode := Episode{
				Season:  ep.Season,
				Episode: ep.Episode,
			}

			if ep.FirstAired != "" {

				aired, err := time.Parse("2006-01-02", ep.FirstAired)

				if err != nil {
					return err
				}

				callback(episode, &aired)

			} else {

				callback(episode, nil)
			}
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
