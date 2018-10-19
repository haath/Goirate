package series

import (
	"fmt"
	"strings"

	"git.gmantaos.com/haath/Goirate/pkg/torrents"
	"git.gmantaos.com/haath/Goirate/pkg/utils"
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

// SearchQuery returns the normalized title of the series along with its episode number,
// as it will be used when searching for torrents.
func (s *Series) SearchQuery(episode Episode) string {

	title := strings.Replace(s.Title, "'s", "s", -1)
	title = utils.NormalizeMediaTitle(title)

	var query string

	if episode.Season == 0 && episode.Episode == 0 {

		query = title

	} else if episode.Episode == 0 {

		query = fmt.Sprintf("%v Season %d", title, episode.Season)

	} else {

		query = fmt.Sprintf("%v %s", title, episode)
	}

	return utils.NormalizeQuery(query)
}

// GetTorrent will search The Pirate Bay and return the best torrent that complies with the given filters.
func (s *Series) GetTorrent(scraper torrents.PirateBayScaper, filters torrents.SearchFilters, episode Episode) (*torrents.Torrent, error) {

	filteredTorrents, err := s.GetTorrents(scraper, filters, episode)

	if err != nil {
		return nil, err
	}

	return torrents.PickVideoTorrent(filteredTorrents, filters)
}

// GetTorrents will attempt to find torrent for an episode of this series.
func (s *Series) GetTorrents(scraper torrents.PirateBayScaper, filters torrents.SearchFilters, episode Episode) ([]torrents.Torrent, error) {

	query := s.SearchQuery(episode)

	title := strings.Replace(s.Title, "'s", "s", -1)
	title = utils.NormalizeMediaTitle(title)

	if episode.Season == 0 && episode.Episode == 0 {

		return scraper.SearchVideoTorrents(query, filters, title)

	} else if episode.Episode == 0 {

		return scraper.SearchVideoTorrents(query, filters, title, fmt.Sprintf("Season %d", episode.Season))

	}

	return scraper.SearchVideoTorrents(query, filters, title, episode.String())
}
