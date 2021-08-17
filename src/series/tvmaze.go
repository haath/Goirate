package series

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"goirate/movies"
	"goirate/utils"
)

// TVmazeCredentials holds the credentials used to authenticate
// with the TvMaze API.
type TVmazeCredentials struct {
}

// TVmazeToken holds a token which represents an authenticated session
// with the TvMaze API.
type TVmazeToken struct {
}

type TVmazeSeries struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Premiered string `json:"premiered"`
}

type apiEndpoint string

const (
	baseEndpoint     apiEndpoint = "https://api.tvmaze.com"
	searchEndpoint   apiEndpoint = baseEndpoint + "/search/shows"
	lookupEndpoint   apiEndpoint = baseEndpoint + "/lookup/shows"
	showEndpoint     apiEndpoint = baseEndpoint + "/shows/%v"
	episodesEndpoint apiEndpoint = baseEndpoint + "/shows/%v/episodes"
)

func (ep apiEndpoint) String() string {
	return string(ep)
}

// Login authenticates with the TVmaze API.
func (cred *TVmazeCredentials) Login() (TVmazeToken, error) {

	var tkn TVmazeToken

	// no login process currently implemented for TvMaze

	return tkn, nil
}

// Search will search the TVmaze for a given query string, IMDB ID, or TVmaze ID and return a list of matching shows.
func (tkn *TVmazeToken) Search(searchName string) ([]TVmazeSeries, error) {

	var searchResults []TVmazeSeries

	var searchResponse []struct {
		Show TVmazeSeries `json:"show"`
	}

	var searchURL string

	if id, err := strconv.Atoi(searchName); err == nil {

		// search query was an integer
		show, err := tkn.lookupSeriesTVmaze(id)

		return []TVmazeSeries{show}, err

	} else if movies.IsIMDbURL(searchName) {

		// search query appears to be imdb url
		imdbID, err := movies.ExtractIMDbID(searchName)
		if err != nil {
			return nil, err
		}

		show, err := tkn.lookupSeriesIMDB(imdbID)

		return []TVmazeSeries{show}, err

	} else if movies.IsIMDbID(searchName) {

		// search query appears to be imdb id
		imdbID, err := movies.FormatIMDbID(searchName)
		if err != nil {
			return nil, err
		}

		show, err := tkn.lookupSeriesIMDB(imdbID)

		return []TVmazeSeries{show}, err

	}

	searchURL = fmt.Sprintf("%v?q=%v", searchEndpoint, url.QueryEscape(searchName))

	err := tkn.apiCall(searchURL, &searchResponse)

	if err == nil && len(searchResponse) > 0 {

		for _, show := range searchResponse {

			searchResults = append(searchResults, show.Show)
		}
	}

	if len(searchResponse) == 0 {

		return nil, fmt.Errorf("show not found on TVmaze")
	}

	return searchResults, nil
}

// Search will search the TVmaze for a given query string or IMDB ID and return the top matching show.
func (tkn *TVmazeToken) SearchFirst(searchName string) (*TVmazeSeries, error) {

	shows, err := tkn.Search(searchName)

	if len(shows) > 0 {

		return &shows[0], err
	}

	return nil, err
}

// LastEpisode uses the TVmaze API to retrieve the last episode that aired
// for a particular series.
func (tkn *TVmazeToken) LastEpisode(seriesID int) (Episode, error) {

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

// NextEpisode uses the TVmaze API to make a best guess as to which episode is sequentially
// next to the one given.
func (tkn *TVmazeToken) NextEpisode(seriesID int, episode Episode) (nextEpisode Episode, err error) {

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

func (tkn *TVmazeToken) lookupSeriesTVmaze(seriesID int) (TVmazeSeries, error) {

	var lookupResponse TVmazeSeries

	searchURL := fmt.Sprintf(string(showEndpoint), seriesID)

	err := tkn.apiCall(searchURL, &lookupResponse)

	return lookupResponse, err
}

func (tkn *TVmazeToken) lookupSeriesIMDB(imdbID string) (TVmazeSeries, error) {

	var lookupResponse TVmazeSeries

	searchURL := fmt.Sprintf("%v?imdb=%v", lookupEndpoint, imdbID)

	err := tkn.apiCall(searchURL, &lookupResponse)

	return lookupResponse, err
}

func (tkn *TVmazeToken) getEpisodes(seriesID int, callback func(Episode)) error {

	var episodeSearchResponse []struct {
		Season  uint   `json:"season"`
		Episode uint   `json:"number"`
		Name    string `json:"name"`
		Airdate string `json:"airdate"`
	}

	episodesURL := fmt.Sprintf(episodesEndpoint.String(), seriesID)
	err := tkn.apiCall(episodesURL, &episodeSearchResponse)

	if err != nil {
		return err
	}

	for _, ep := range episodeSearchResponse {

		var aired *time.Time

		if ep.Airdate != "" {

			tmp, err := time.Parse("2006-01-02", ep.Airdate)

			if err != nil {
				return err
			}

			aired = &tmp
		}

		episode := Episode{
			Season:  ep.Season,
			Episode: ep.Episode,
			Title:   ep.Name,
			Aired:   aired,
		}

		callback(episode)
	}

	return nil
}

func (tkn *TVmazeToken) apiCall(url string, v interface{}) error {

	httpClient := utils.HTTPClient{}

	return httpClient.GetJSON(url, &v)
}

// EnvTVmazeCredentials returns a struct variable which contains credentials for the TVmaze API.
func EnvTVmazeCredentials() TVmazeCredentials {
	return TVmazeCredentials{}
}
