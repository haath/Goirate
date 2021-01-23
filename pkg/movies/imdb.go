package movies

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"gitlab.com/haath/goirate/pkg/utils"
)

// BaseURL is the base for IMDb URLS.
const BaseURL = "https://www.imdb.com"

// FormatIMDbID formats the given ID to its canonical 7-digit value.
func FormatIMDbID(id string) (string, error) {

	id = strings.TrimLeft(id, "t")

	idNum, err := strconv.Atoi(id)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("tt%07d", idNum), nil
}

// IsIMDbID returns true if the string is in a valid IMDbID format.
func IsIMDbID(str string) bool {
	_, err := FormatIMDbID(str)
	return err == nil
}

// IsIMDbURL returns true if the string is in a valid IMDb URL.
func IsIMDbURL(str string) bool {
	_, err := ExtractIMDbID(str)
	return err == nil
}

// ExtractIMDbID will extract the IMDb ID of a movie from its URL.
// Assuming that the URL is in the format: https://www.imdb.com/title/tt0848228/
func ExtractIMDbID(url string) (string, error) {
	r, _ := regexp.Compile(`/title/tt(\w+)/?`)
	m := r.FindStringSubmatch(url)

	if len(m) > 0 && len(m[1]) >= 7 {

		return FormatIMDbID(m[1])
	}

	return "", errors.New("error extracting IMDb ID from: " + url)
}

// ParseIMDbPage extracts a movie's details from its IMDb page.
func ParseIMDbPage(doc *goquery.Document) Movie {

	posterURL, _ := doc.Find(".poster > a > img").Attr("src")
	year, _ := strconv.Atoi(doc.Find("#titleYear > a").Text())

	doc.Find("#titleYear").Remove()

	title := strings.TrimSpace(doc.Find(".title_wrapper > h1").Text())
	rating, _ := strconv.ParseFloat(doc.Find(".ratingValue strong").Text(), 32)
	duration := parseDuration(strings.TrimSpace(doc.Find(".subtext time").Text()))
	imdbURL, _ := doc.Find("link[rel='canonical']").Attr("href")
	imdbID, _ := ExtractIMDbID(imdbURL)
	genres := extractGenres(doc)

	movie := Movie{
		PosterURL: posterURL,
		Rating:    float32(int(rating*10)) / 10, // Round to one decimal, just in case,
		Duration:  duration,
		MovieID: MovieID{
			Title:  title,
			IMDbID: imdbID,
			Year:   uint(year),
		},
		Genres: genres,
	}

	origTitle := doc.Find(".originalTitle")

	if origTitle.Length() > 0 {
		origTitle.Find(".description").Remove()

		movie.AltTitle = origTitle.Text()
	}

	return movie
}

// ParseSearchPage will parse the result page of an IMDb search and return the titles
// and IDs of all movies in it.
func ParseSearchPage(doc *goquery.Document) []MovieID {

	var movies []MovieID

	doc.Find(".findSection td.result_text").Each(func(i int, row *goquery.Selection) {

		title := row.Find("a").Text()
		imdbURL, _ := row.Find("a").Attr("href")
		id, _ := ExtractIMDbID(imdbURL)

		row.Find("a").Remove()

		year := extractYear(row.Text())
		altTitle := extractAltTitle(row.Text())

		movie := MovieID{
			Title:    title,
			IMDbID:   id,
			Year:     uint(year),
			AltTitle: altTitle,
		}

		movies = append(movies, movie)
	})

	return movies
}

func extractGenres(doc *goquery.Document) []string {

	var genres []string

	doc.Find(".subtext a").Each(func(i int, item *goquery.Selection) {

		_, titleExists := item.Attr("title")

		// The only subtext with a title attribute is the release date.
		// The rest are genres.
		if !titleExists {
			genres = append(genres, item.Text())
		}
	})

	return genres
}

func searchURL(query string) string {

	searchURL, _ := url.Parse(BaseURL)
	searchURL.Path = "/find"

	params := searchURL.Query()
	params.Add("s", "tt")
	params.Add("ttype", "ft")
	params.Add("q", url.QueryEscape(query))

	searchURL.RawQuery = params.Encode()

	return searchURL.String()
}

func parseDuration(durationString string) int {

	minutes := 0

	r, _ := regexp.Compile(`(\d+)h`)
	m := r.FindStringSubmatch(durationString)

	if len(m) > 0 {
		hours, _ := strconv.Atoi(m[1])

		minutes += hours * 60
	}

	r, _ = regexp.Compile(`(\d+)min`)
	m = r.FindStringSubmatch(durationString)

	if len(m) > 0 {
		min, _ := strconv.Atoi(m[1])

		minutes += min
	}

	return minutes
}

func extractYear(searchRow string) int {
	r, _ := regexp.Compile(`\((\d+)\)`)
	m := r.FindStringSubmatch(searchRow)

	if len(m) > 0 {
		year, _ := strconv.Atoi(m[1])

		return year
	}

	return 0
}

func extractAltTitle(searchRow string) string {
	r, _ := regexp.Compile(`aka\s+"(.+)"`)
	m := r.FindStringSubmatch(searchRow)

	if len(m) > 0 {
		return m[1]
	}

	return ""
}

// GetMovie will scrape the IMDb page of the movie with the given id
// and return its details.
func GetMovie(imdbID string) (*Movie, error) {

	tmp := MovieID{
		IMDbID: imdbID,
	}

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

// Search performs a text search on IMDb, limited to movies, and returns the results.
func Search(query string) ([]MovieID, error) {
	url := searchURL(query)
	doc, err := utils.HTTPGet(url)

	if err != nil {
		return nil, err
	}

	movies, err := ParseSearchPage(doc), nil

	if len(movies) == 0 {
		err = fmt.Errorf("Movie not found on IMDB: %v", query)
	}

	return movies, err
}
