package movies

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"gitlab.com/haath/goirate/pkg/torrents"
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
	Duration  int      `json:"duration"`
	Rating    float32  `json:"rating"`
	PosterURL string   `json:"poster_url"`
	Genres    []string `json:"genres"`
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

// GetGenresString returns the movie's genres as a comma-separated string.
func (m Movie) GetGenresString() string {

	return strings.Join(m.Genres, ", ")
}

// GetSearchTerms returns a list of distinct substrings that should be present in the title of a torrent for this movie.
func (m Movie) GetSearchTerms(useAltTitle bool) []string {

	searchTerms := []string{m.Title}

	if useAltTitle {
		searchTerms[0] = m.AltTitle
	}

	if m.Year != 0 {

		searchTerms = append(searchTerms, strconv.Itoa(int(m.Year)))
	}

	return searchTerms
}

// GetSearchQuery returns the string that should be used in the query, when searching for torrents
// for this movie.
func (m Movie) GetSearchQuery(useAltTitle bool) string {

	searchQuery := m.Title

	if useAltTitle {
		searchQuery = m.AltTitle
	}

	if m.Year != 0 {

		searchQuery = fmt.Sprintf("%s %v", searchQuery, m.Year)
	}

	return searchQuery
}

// GetTorrent will search The Pirate Bay and return the best torrent that complies with the given filters.
func (m Movie) GetTorrent(filters torrents.SearchFilters) (*torrents.Torrent, error) {

	filteredTorrents, err := m.GetTorrents(filters)

	if err != nil {
		return nil, err
	}

	return torrents.PickVideoTorrent(filteredTorrents, filters)
}

// GetTorrents will search The Pirate Bay for torrents of this movie that comply with the given filters.
// It will return one torrent for each video quality.
func (m Movie) GetTorrents(filters torrents.SearchFilters) ([]torrents.Torrent, error) {

	filters.SearchTerms = m.GetSearchTerms(false)
	trnts, err := filters.SearchVideoTorrents(m.GetSearchQuery(false))

	if m.AltTitle != "" {

		filters.SearchTerms = m.GetSearchTerms(true)
		altTitleTorrents, err := filters.SearchVideoTorrents(m.GetSearchQuery(true))

		if err != nil {
			return nil, err
		}

		trnts = append(trnts, altTitleTorrents...)
	}

	torrentsQualityMap, _ := torrents.SearchVideoTorrentList(trnts, filters)
	var perQualitySlice []torrents.Torrent
	for _, value := range torrentsQualityMap {
		perQualitySlice = append(perQualitySlice, *value)
	}
	trnts = perQualitySlice

	return trnts, err
}
