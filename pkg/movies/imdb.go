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

	idNum, err := strconv.Atoi(id)

	if err != nil {
		return "", err
	}

	if idNum < 0 || idNum > 9999999 {
		return "", errors.New("Invalid IMDb ID: " + id)
	}

	return fmt.Sprintf("%07d", idNum), nil
}

// ParseIMDbPage extracts a movie's details from its IMDb page.
func ParseIMDbPage(doc *goquery.Document) Movie {

	posterURL, _ := doc.Find(".poster > a > img").Attr("src")
	year, _ := strconv.Atoi(doc.Find("#titleYear > a").Text())

	doc.Find("#titleYear").Remove()

	title := strings.TrimSpace(doc.Find(".title_wrapper > h1").Text())
	rating, _ := strconv.ParseFloat(doc.Find(".ratingValue span[itemprop='ratingValue']").Text(), 32)
	duration := parseDuration(strings.TrimSpace(doc.Find("time[itemprop='duration']").Text()))

	return Movie{
		Title:     title,
		PosterURL: posterURL,
		Year:      year,
		Rating:    float32(int(rating*10)) / 10, // Round to one decimal, just in case,
		Duration:  duration,
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
