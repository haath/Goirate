package movies

import (
	"errors"
	"fmt"
	"strconv"
)

// FormatIMDbID formats the given ID to its canonical 7-digit value
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
