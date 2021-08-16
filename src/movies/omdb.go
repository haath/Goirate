package movies

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"goirate/utils"
)

type apiEndpoint string

const (
	baseEndpoint apiEndpoint = "https://www.omdbapi.com"
)

// OMDBCredentials holds the API key for access to the OMDB API.
type OMDBCredentials struct {
	APIKey string `toml:"api_key"`
}

// IsEnabled returns true if an API key has been provided for the OMDb API.
func (omdb *OMDBCredentials) IsEnabled() bool {

	if omdb.APIKey == "" {

		envCred := EnvOMDBCredentials()
		omdb.APIKey = envCred.APIKey
	}

	return omdb.APIKey != ""
}

// GetMovie fetches a movie using the OMDB API.
func (omdb *OMDBCredentials) GetMovie(imdbID string) (*Movie, error) {

	var omdbMovieResponse struct {
		ImdbID  string `json:"imdbID"`
		Title   string `json:"Title"`
		Year    string `json:"Year"`
		Runtime string `json:"Runtime"`
		Genres  string `json:"Genre"`
		Rating  string `json:"imdbRating"`
		Poster  string `json:"Poster"`
	}

	formattedID, err := FormatIMDbID(imdbID)

	if err != nil {
		return nil, err
	}

	baseURL, err := omdb.getRequestBaseURL()

	if err != nil {
		return nil, err
	}

	reqURL := fmt.Sprintf("%v&i=%v", baseURL, formattedID)

	httpClient := utils.HTTPClient{}

	err = httpClient.GetJSON(reqURL, &omdbMovieResponse)

	if err != nil {
		return nil, err
	}

	year, _ := strconv.Atoi(omdbMovieResponse.Year)
	duration, _ := strconv.Atoi(strings.Split(omdbMovieResponse.Runtime, " ")[0])
	rating, _ := strconv.ParseFloat(omdbMovieResponse.Rating, 32)
	genres := strings.Split(omdbMovieResponse.Genres, ", ")

	movie := Movie{
		MovieID: MovieID{
			IMDbID:   omdbMovieResponse.ImdbID,
			Title:    omdbMovieResponse.Title,
			Year:     uint(year),
			AltTitle: "",
		},
		Duration:  duration,
		Rating:    float32(rating),
		PosterURL: omdbMovieResponse.Poster,
		Genres:    genres,
	}

	return &movie, nil
}

// Search searches for a movie on the OMDb API, given a string as query.
func (omdb *OMDBCredentials) Search(query string) ([]MovieID, error) {

	var omdbSearchResponse struct {
		Search []struct {
			ImdbID string `json:"imdbID"`
			Title  string `json:"Title"`
			Year   string `json:"Year"`
		} `json:"Search"`
	}

	baseURL, err := omdb.getRequestBaseURL()

	if err != nil {
		return nil, err
	}

	formattedQuery := url.QueryEscape(query)

	reqURL := fmt.Sprintf("%v&s=%v", baseURL, formattedQuery)

	httpClient := utils.HTTPClient{}

	err = httpClient.GetJSON(reqURL, &omdbSearchResponse)

	if err != nil {
		return nil, err
	}

	var movies []MovieID

	for _, movie := range omdbSearchResponse.Search {

		year, _ := strconv.Atoi(movie.Year)

		movies = append(movies, MovieID{
			IMDbID:   movie.ImdbID,
			Title:    movie.Title,
			Year:     uint(year),
			AltTitle: "",
		})
	}

	return movies, nil
}

func (omdb *OMDBCredentials) getRequestBaseURL() (string, error) {

	envCred := EnvOMDBCredentials()

	if envCred.APIKey != "" && omdb.APIKey == "" {

		omdb.APIKey = envCred.APIKey
	}

	if omdb.APIKey == "" {

		return "", fmt.Errorf("fetching movie data requires an API key for OMDB. (http://www.omdbapi.com/apikey.aspx)")
	}

	baseURL := fmt.Sprintf("%v/?apikey=%v", baseEndpoint, omdb.APIKey)

	return baseURL, nil
}

// EnvOMDBCredentials returns an OMDBCredentials struct variable which contains
// the credentials found in the TVDB_API_KEY.
func EnvOMDBCredentials() OMDBCredentials {
	return OMDBCredentials{
		APIKey: os.Getenv("GOIRATE_OMDB_API_KEY"),
	}
}
