package movies

import (
	"errors"
	"fmt"
	"git.gmantaos.com/haath/Goirate/pkg/utils"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"regexp"
	"strconv"
	"strings"
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

	if idNum < 0 || idNum > 9999999 {
		return "", errors.New("Invalid IMDb ID: " + id)
	}

	return fmt.Sprintf("%07d", idNum), nil
}

// ExtractIMDbID will extract the IMDb ID of a movie from its URL.
// Assuming that the URL is in the format: https://www.imdb.com/title/tt0848228/
func ExtractIMDbID(url string) (string, error) {
	r, _ := regexp.Compile(`/title/tt(\w+)/?`)
	m := r.FindStringSubmatch(url)

	if len(m) > 0 {
		id, _ := FormatIMDbID(m[1])
		return id, nil
	}

	return "", errors.New("error extracting IMDb ID from: " + url)
}

// ParseIMDbPage extracts a movie's details from its IMDb page.
func ParseIMDbPage(doc *goquery.Document) Movie {

	posterURL, _ := doc.Find(".poster > a > img").Attr("src")
	year, _ := strconv.Atoi(doc.Find("#titleYear > a").Text())

	doc.Find("#titleYear").Remove()

	title := strings.TrimSpace(doc.Find(".title_wrapper > h1").Text())
	rating, _ := strconv.ParseFloat(doc.Find(".ratingValue span[itemprop='ratingValue']").Text(), 32)
	duration := parseDuration(strings.TrimSpace(doc.Find("time[itemprop='duration']").Text()))
	imdbURL, _ := doc.Find("link[rel='canonical']").Attr("href")
	imdbID, _ := ExtractIMDbID(imdbURL)

	return Movie{
		PosterURL: posterURL,
		Rating:    float32(int(rating*10)) / 10, // Round to one decimal, just in case,
		Duration:  duration,
		MovieID: MovieID{
			Title:  title,
			IMDbID: imdbID,
			Year:   year,
		},
	}
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

		movie := MovieID{
			Title:  title,
			IMDbID: id,
			Year:   year,
		}

		movies = append(movies, movie)
	})

	return movies
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
