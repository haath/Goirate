package movies

import (
	"errors"
	"fmt"
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

	if idNum < 0 {
		return "", errors.New("negative number in IMDB ID")
	}

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
