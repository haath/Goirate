package main

import (
	"encoding/json"
	"fmt"
	"log"

	"gitlab.com/haath/goirate/pkg/movies"
	"gitlab.com/haath/goirate/pkg/torrents"
)

// MovieCommand defines the movie command and holds its options.
type MovieCommand struct {
	torrentSearchArgs

	Year      uint                `short:"y" long:"year" description:"The release year of the movie. Used when searching for the movie by title instead of by IMDbID."`
	Download  bool                `short:"d" long:"download" description:"Send the movie to the qBittorret client for download using the RPC configuration."`
	NoTorrent bool                `long:"no-torrent" description:"Do not search for torrents."`
	Args      moviePositionalArgs `positional-args:"1" required:"1"`
}

type moviePositionalArgs struct {
	Query string `positional-arg-name:"<imdb url | imdbID | title>"`
}

// Execute is the callback of the movie command.
func (m *MovieCommand) Execute(args []string) error {

	movie, err := m.getMovie()

	if err != nil {
		return err
	}

	var perQualityTorrents []torrents.Torrent
	var topTorrent *torrents.Torrent

	if !m.NoTorrent {

		perQualityTorrents, err = movie.GetTorrents(*m.GetFilters())

		if len(perQualityTorrents) > 0 {

			topTorrent, err = torrents.PickVideoTorrent(perQualityTorrents, *m.GetFilters())

			if err != nil {
				return err
			}

			if m.Download && topTorrent != nil {

				// Send the torrent to the qBittorrent daemon for download
				err = m.downloadMovieTorrent(movie, topTorrent)

				if err != nil {
					return err
				}
			}
		}
	}

	if Options.JSON {

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

		if m.MagnetLink && !m.NoTorrent {

			if topTorrent == nil {
				return fmt.Errorf("no torrent found for: %v", movie.Title)
			}

			log.Println(topTorrent.Magnet)

		} else if m.TorrentURL && !m.NoTorrent {

			if topTorrent == nil {
				return fmt.Errorf("no torrent found for: %v", movie.Title)
			}

			log.Println(topTorrent.FullURL())

		} else {

			log.Printf("%v\n", movie.Title)

			if movie.AltTitle != "" {
				log.Printf("Orig. Title:\t%v\n", movie.AltTitle)
			}

			log.Printf("IMDbID:\t\t%v\n", movie.IMDbID)
			log.Printf("Year:\t\t%v\n", movie.Year)
			log.Printf("Rating:\t\t%v\n", movie.Rating)
			log.Printf("Genres:\t\t%v\n", movie.GetGenresString())

			if movie.FormattedDuration() != "" {
				log.Printf("Duration:\t%v\n", movie.FormattedDuration())
			}

			if movie.PosterURL != "" {
				log.Printf("Poster:\t\t%v\n", movie.PosterURL)
			}

			log.Println("")

			if !m.NoTorrent {

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

	}

	return nil
}

func (m *MovieCommand) getMovie() (*movies.Movie, error) {

	var err error
	var imdbID string

	omdb := Config.OMDBCredentials

	if movies.IsIMDbID(m.Args.Query) {

		imdbID, err = movies.FormatIMDbID(m.Args.Query)

	} else if movies.IsIMDbURL(m.Args.Query) {

		imdbID, err = movies.ExtractIMDbID(m.Args.Query)

	} else {

		imdbID, err = m.searchMovie()
	}

	if err != nil {

		return nil, err
	}

	var movie *movies.Movie

	if omdb.IsEnabled() {

		movie, err = omdb.GetMovie(imdbID)

	} else {

		// OMDb API key not configured, fall back to IMDb.
		movie, err = movies.GetMovie(imdbID)
	}

	return movie, err
}

func (m *MovieCommand) searchMovie() (string, error) {

	var searchResults []movies.MovieID
	var err error

	omdb := Config.OMDBCredentials

	if omdb.IsEnabled() {

		searchResults, err = omdb.Search(m.Args.Query)

	} else {

		searchResults, err = movies.Search(m.Args.Query)
	}

	if err != nil {
		return "", err
	}

	for _, movie := range searchResults {

		if m.Year == 0 || m.Year == movie.Year {

			return movie.IMDbID, nil
		}
	}

	return "", fmt.Errorf("movie not found: %v", m.Args.Query)
}

func (m *MovieCommand) downloadMovieTorrent(movie *movies.Movie, torrent *torrents.Torrent) error {

	qbt, err := Config.QBittorrentConfig.GetClient()

	if err != nil {
		return err
	}

	downloadPath := Config.DownloadDir.Movies

	if !Options.JSON && !m.MagnetLink && !m.TorrentURL {

		log.Printf("Downloading: %s (%s)\n", movie.Title, downloadPath)
	}

	return qbt.AddTorrent(torrent.Magnet, downloadPath)
}
