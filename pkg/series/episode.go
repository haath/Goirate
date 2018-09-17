package series

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Episode represents a unique episode of a series, identified by a
// pair of a season and episode number.
type Episode struct {
	Season  uint `toml:"season" json:"season"`
	Episode uint `toml:"episode" json:"episode"`
}

// ParseEpisodeString will extract the season and episode number from a string
// description, such as S03E07.
func ParseEpisodeString(episodeStr string) Episode {

	episode := Episode{1, 1}

	episodeStr = strings.ToLower(episodeStr)

	r, _ := regexp.Compile(`(?:\s*|\d*)(?:s|se|season)\s*(\d+)`)

	m := r.FindStringSubmatch(episodeStr)

	if len(m) > 0 {

		s, _ := strconv.Atoi(m[1])
		episode.Season = uint(s)
	}

	r, _ = regexp.Compile(`(?:\s+|\d+)(?:e|ep|episode)\s*(\d+)`)

	m = r.FindStringSubmatch(episodeStr)

	if len(m) > 0 {

		e, _ := strconv.Atoi(m[1])
		episode.Episode = uint(e)
	}

	return episode
}

// IsAfter returns true if this episode is sequentially after the given episode.
func (ep *Episode) IsAfter(episode Episode) bool {
	return ep.Season > episode.Season ||
		ep.Episode > episode.Episode
}

// String returns the string SxxEyy representation of an episode.
func (ep Episode) String() string {
	return fmt.Sprintf("S%02dE%02d", ep.Season, ep.Episode)
}
