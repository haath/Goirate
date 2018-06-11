package movies

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
	"strings"
)

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
	r, _ := regexp.Compile(`https://www.imdb.com/title/tt(\w+)/?`)
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
		Title:     title,
		PosterURL: posterURL,
		Year:      year,
		Rating:    float32(int(rating*10)) / 10, // Round to one decimal, just in case,
		Duration:  duration,
		IMDbID:    imdbID,
	}
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
