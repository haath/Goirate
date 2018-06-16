package main

import (
	"encoding/json"
	"git.gmantaos.com/haath/Goirate/pkg/movies"
	"git.gmantaos.com/haath/Goirate/pkg/torrents"
	"log"
)

// MovieCommand defines the movie command and holds its options.
type MovieCommand struct {
	torrents.SearchFilters
	torrentSearchArgs

	Args moviePositionalArgs `positional-args:"1" required:"1"`
}

type moviePositionalArgs struct {
	IMDbID string `positional-arg-name:"imdbID"`
}

// Execute is the callback of the movie command.
func (m *MovieCommand) Execute(args []string) error {

	movie, err := movies.GetMovie(m.Args.IMDbID)

	if err != nil {
		return err
	}

	if Options.JSON {

		movieJSON, err := json.MarshalIndent(movie, "", "   ")

		if err != nil {
			return err
		}

		log.Println(string(movieJSON))

	} else {

		log.Printf("%v\n\n", movie.Title)

		if movie.AltTitle != "" {
			log.Printf("Orig. Title:\t%v\n", movie.AltTitle)
		}

		log.Printf("IMDbID:\t\t%v\n", movie.IMDbID)
		log.Printf("Year:\t\t%v\n", movie.Year)
		log.Printf("Rating:\t\t%v\n", movie.Rating)

		if movie.FormattedDuration() != "" {
			log.Printf("Duration:\t%v\n", movie.FormattedDuration())
		}

		if movie.PosterURL != "" {
			log.Printf("Poster:\t\t%v\n", movie.PosterURL)
		}

	}

	return nil
}
