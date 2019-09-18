package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/olekukonko/tablewriter"
	imdb "gitlab.com/haath/goirate/pkg/movies"
)

// MovieSearchCommand is the command used to search for movies on IMDb.
type MovieSearchCommand struct {
	Year  uint16 `short:"y" long:"year" description:"The release year to limit the movie search to."`
	Count uint   `short:"c" long:"count" description:"Limit the number of results."`

	Args positionalArgs `positional-args:"1" required:"1"`
}

// Execute is the callback of the movie command.
func (c *MovieSearchCommand) Execute(args []string) error {

	movies, err := imdb.Search(c.Args.Query)

	if err != nil {
		return err
	}

	if c.Count > 0 && uint(len(movies)) > c.Count {
		movies = movies[:c.Count]
	}

	if Options.JSON {
		moviesJSON, err := json.MarshalIndent(movies, "", "   ")

		if err != nil {
			return err
		}

		log.Println(string(moviesJSON))
	} else {
		log.Printf(getMoviesTable(movies))
	}

	return nil
}

func getMoviesTable(movies []imdb.MovieID) string {
	buf := bytes.NewBufferString("")

	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"IMDb ID", "Title", "Year"})
	table.SetColumnAlignment([]int{tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_DEFAULT})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoFormatHeaders(false)

	for _, movie := range movies {

		table.Append([]string{movie.IMDbID, movie.Title, fmt.Sprint(movie.Year)})
	}

	table.Render()

	return buf.String()
}
