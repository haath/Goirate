package series

import (
	"fmt"
	"strings"

	"goirate/torrents"
	"goirate/utils"
)

// Series holds the title of a series along
// with the next episode expected to come out.
type Series struct {
	ID               int                    `toml:"id" json:"id"`
	Title            string                 `toml:"title" json:"title"`
	MinQuality       torrents.VideoQuality  `toml:"min_quality" json:"min_quality"`
	VerifiedUploader bool                   `toml:"only_trusted" json:"only_trusted"`
	LastEpisode      Episode                `toml:"last_episode" json:"last_episode"`
	Actions          utils.WatchlistActions `toml:"actions" json:"actions"`
}

// NextEpisode uses the TVDB API to make a best guess as to which is the next episode
// to this series' LastEpisode.
func (s *Series) NextEpisode(tkn *TVDBToken) (Episode, error) {

	return tkn.NextEpisode(s.ID, s.LastEpisode)
}

// GetSearchQuery returns the normalized title of the series along with its episode number,
// as it will be used when searching for torrents.
func (s *Series) GetSearchQuery(episode Episode) string {

	title := s.getNormalizedTitle()

	var searchQuery string

	if episode.Season == 0 && episode.Episode == 0 {

		searchQuery = title

	} else if episode.Episode == 0 {

		searchQuery = fmt.Sprintf("%v Season %d", title, episode.Season)

	} else {

		searchQuery = fmt.Sprintf("%v %s", title, episode)
	}

	return utils.NormalizeQuery(searchQuery)
}

// GetSearchTerms returns a list of strings that should exit in a torrent for the given episode of the series.
func (s *Series) GetSearchTerms(episode Episode) []string {

	title := s.getNormalizedTitle()

	searchTerms := []string{title}

	if episode.Season == 0 && episode.Episode == 0 {

		// Add nothing.

	} else if episode.Episode == 0 {

		searchTerms = append(searchTerms, fmt.Sprintf("Season %d", episode.Season))

	} else {

		searchTerms = append(searchTerms, episode.String())
	}

	return searchTerms
}

// GetTorrent will search The Pirate Bay and return the best torrent that complies with the given filters.
func (s *Series) GetTorrent(filters torrents.SearchFilters, episode Episode) (*torrents.Torrent, error) {

	filteredTorrents, err := s.GetTorrents(filters, episode)

	if err != nil {
		return nil, err
	}

	return torrents.PickVideoTorrent(filteredTorrents, filters)
}

// GetTorrents will attempt to find a torrent for an episode of this series.
func (s *Series) GetTorrents(filters torrents.SearchFilters, episode Episode) ([]torrents.Torrent, error) {

	searchQuery := s.GetSearchQuery(episode)
	filters.SearchTerms = s.GetSearchTerms(episode)

	return filters.SearchVideoTorrents(searchQuery)
}

func (s *Series) getNormalizedTitle() string {

	title := strings.Replace(s.Title, "'s", "s", -1)
	title = utils.NormalizeMediaTitle(title)
	return title
}
