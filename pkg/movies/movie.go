package movies

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
)

// MovieID holds the defining properties of an IMDb movie as they appear in search results.
type MovieID struct {
	IMDbID   string `json:"imdb_id"`
	Title    string `json:"title"`
	Year     uint   `json:"year"`
	AltTitle string `json:"alt_title"`
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

// FormattedDuration returns the duration of the movie in human-readable format.
func (m Movie) FormattedDuration() string {
	hours := m.Duration / 60
	minutes := m.Duration % 60

	var buf bytes.Buffer

	if hours > 0 {
		buf.WriteString(fmt.Sprintf("%vh ", hours))
	}

	if minutes > 0 {
		buf.WriteString(fmt.Sprintf("%vmin", minutes))
	}

	return strings.TrimSpace(buf.String())
}
