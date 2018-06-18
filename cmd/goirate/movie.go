package main

import (
	"encoding/json"
	"errors"
	"log"

	"git.gmantaos.com/haath/Goirate/pkg/movies"
	"git.gmantaos.com/haath/Goirate/pkg/torrents"
)

// MovieCommand defines the movie command and holds its options.
type MovieCommand struct {
	torrents.SearchFilters
	torrentSearchArgs

	Year uint                `short:"y" long:"year" description:"The release year of the movie. Used when searching for the movie by title instead of by IMDbID."`
	Args moviePositionalArgs `positional-args:"1" required:"1"`
}

type moviePositionalArgs struct {
	Query string `positional-arg-name:"imdbID/title"`
}

// Execute is the callback of the movie command.
func (m *MovieCommand) Execute(args []string) error {

	if !m.ValidOutputFlags() {
		return errors.New("too many flags specifying the kind of output")
	}

	var movie *movies.Movie
	var err error

	if movies.IsIMDbID(m.Args.Query) {

		movie, err = movies.GetMovie(m.Args.Query)

		if err != nil {
			return err
		}

	} else {

		movie, err = m.findMovie()

		if err != nil {
			return err
		}

	}

	scraper, err := m.GetScraper()
	filters := &m.SearchFilters

	if err != nil {
		return err
	}

	if Options.JSON {

		perQualityTorrents, err := movie.GetTorrents(scraper, filters)

		if err != nil {
			return err
		}

		movieObj := struct {
			movies.Movie
			Torrents []torrents.Torrent `json:"torrents"`
		}{
			*movie,
			perQualityTorrents,
		}

		movieJSON, err := json.MarshalIndent(movieObj, "", "   ")

		if err != nil {
			return err
		}

		log.Println(string(movieJSON))

	} else {

		topTorrent, err := movie.GetTorrent(scraper, filters)

		if err != nil {
			return err
		}

		if m.MagnetLink {

			log.Println(topTorrent.Magnet)

		} else if m.TorrentURL {

			log.Println(topTorrent.FullURL())

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

			log.Println("")

			if topTorrent == nil {
				log.Println("No torrent found")
			} else {

				log.Println(topTorrent.Title)

				log.Printf("URL:\t\t%v\n", topTorrent.FullURL())
				log.Printf("Seeds/Peers:\t%v\n", topTorrent.PeersString())
				log.Printf("Size:\t\t%v\n", topTorrent.SizeString())
				log.Printf("Trusted:\t%v\n", topTorrent.VerifiedUploader)
				log.Printf("Magnet:\n%v\n", topTorrent.Magnet)
			}

		}

	}

	return nil
}

func (m *MovieCommand) findMovie() (*movies.Movie, error) {

	searchResults, err := movies.Search(m.Args.Query)

	if err != nil {
		return nil, err
	}

	for _, movie := range searchResults {

		if m.Year == 0 || m.Year == movie.Year {

			return movies.GetMovie(movie.IMDbID)

		}

	}

	return nil, errors.New("movie not found")
}
