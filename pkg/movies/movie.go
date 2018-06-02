package movies

import (
	"fmt"
	"net/url"
)

// Movie holds all the information regarding a movie on IMDb.
type Movie struct {
	IMDbID    string  `json:"imdb_id"`
	Title     string  `json:"title"`
	Year      int     `json:"year"`
	Duration  int     `json:"duration"`
	Rating    float32 `json:"rating"`
	PosterURL string  `json:"poster_url"`
}

// GetURL formats the IMDbID of the movie object and returns the full
// URL to the movie's page on IMDb.
func (m Movie) GetURL() (*url.URL, error) {
	formattedID, err := FormatIMDbID(m.IMDbID)

	if err != nil {
		return nil, err
	}

	urlString := fmt.Sprintf("https://www.imdb.com/title/tt%v/", formattedID)

	return url.Parse(urlString)
}
