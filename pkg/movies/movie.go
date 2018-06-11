package movies

import (
	"fmt"
	"git.gmantaos.com/haath/Goirate/pkg/utils"
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

// GetMovie will scrape the IMDb page of the movie with the given id
// and return its details.
func GetMovie(imdbID string) (*Movie, error) {
	var tmp Movie
	tmp.IMDbID = imdbID

	url, err := tmp.GetURL()

	if err != nil {
		return nil, err
	}

	doc, err := utils.HTTPGet(url.String())

	if err != nil {
		return nil, err
	}

	movie := ParseIMDbPage(doc)

	return &movie, nil
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
