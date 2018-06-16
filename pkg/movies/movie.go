package movies

import (
	"fmt"
	"net/url"
)

// MovieID holds the defining properties of an IMDb movie as they appear in search results.
type MovieID struct {
	IMDbID string `json:"imdb_id"`
	Title  string `json:"title"`
	Year   int    `json:"year"`
}

// Movie holds all the information regarding a movie on IMDb.
type Movie struct {
	MovieID
	Duration  int     `json:"duration"`
	Rating    float32 `json:"rating"`
	PosterURL string  `json:"poster_url"`
}

// GetURL formats the IMDbID of the movie object and returns the full
// URL to the movie's page on IMDb.
func (m MovieID) GetURL() (*url.URL, error) {
	formattedID, err := FormatIMDbID(m.IMDbID)

	if err != nil {
		return nil, err
	}

	urlString := fmt.Sprintf(BaseURL+"/title/tt%v/", formattedID)

	return url.Parse(urlString)
}
