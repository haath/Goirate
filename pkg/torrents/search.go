package torrents

import (
	"sort"

	"gitlab.com/haath/gobytes"
)

// UploaderFilters holds filters regarding the acceptance of a torrent's uploader.
type UploaderFilters struct {
	Whitelist []string `long:"whitelist" description:"Add to a whitelist of uploaders, to only consider torrents from them." toml:"whitelist"`
	Blacklist []string `long:"blacklist" description:"Add to a blacklist of uploaders, to avoid torrents from them." toml:"blacklist"`
}

// SearchFilters holds conditions and filters, used to search for specific torrents.
type SearchFilters struct {
	VerifiedUploader bool            `long:"trusted" description:"Only consider torrents where the uploader is either VIP or Trusted." toml:"trusted"`
	MinQuality       VideoQuality    `long:"min-quality" description:"Minimum acceptable torrent quality (inclusive)." toml:"min-quality"`
	MaxQuality       VideoQuality    `long:"max-quality" description:"Maximum acceptable torrent quality (inclusive)." toml:"max-quality"`
	MinSize          string          `long:"min-size" description:"Minimum acceptable torrent size." toml:"min-size"`
	MaxSize          string          `long:"max-size" description:"Maximum acceptable torrent size." toml:"max-size"`
	MinSeeders       int             `long:"min-seeders" description:"Minimum acceptable amount of seeders." toml:"min-seeders"`
	Uploaders        UploaderFilters `toml:"uploaders"`
}

// MinSizeKB returns the specified minimum size in kilobytes.
func (f SearchFilters) MinSizeKB() (int64, error) {
	var v gobytes.ByteSize

	err := v.UnmarshalText([]byte(f.MinSize))

	return int64(v.KBytes()), err
}

// MaxSizeKB returns the specified minimum size in kilobytes.
func (f SearchFilters) MaxSizeKB() (int64, error) {
	var v gobytes.ByteSize

	err := v.UnmarshalText([]byte(f.MaxSize))

	return int64(v.KBytes()), err
}

// UploaderOk will return true if the given uploader's name is acceptable according to the blacklist and whitelist
// of the filters.
func (f SearchFilters) UploaderOk(uploader string) bool {

	contains := func(s []string, e string) bool {
		for _, a := range s {
			if a == e {
				return true
			}
		}
		return false
	}

	return (len(f.Uploaders.Blacklist) == 0 || !contains(f.Uploaders.Blacklist, uploader)) &&
		(len(f.Uploaders.Whitelist) == 0 || contains(f.Uploaders.Whitelist, uploader))
}

// IsOk returns true if the given torrent complies with the filters.
func (f SearchFilters) IsOk(torrent *Torrent) bool {

	maxSize, _ := f.MaxSizeKB()

	minSize, _ := f.MinSizeKB()

	if f.VerifiedUploader && !torrent.VerifiedUploader {
		return false
	}

	if (torrent.Size > maxSize && maxSize > 0) ||
		(torrent.Size < minSize && minSize > 0) {
		return false
	}

	if (f.MinQuality != "" && torrent.VideoQuality.WorseThan(f.MinQuality)) ||
		(f.MaxQuality != "" && torrent.VideoQuality.BetterThan(f.MaxQuality)) {
		return false
	}

	if !f.UploaderOk(torrent.Uploader) {
		return false
	}

	if torrent.Seeders < f.MinSeeders {
		return false
	}

	return true
}

// FilterTorrents filters the given list of torrents, returning only the ones that
// comply with the filters.
func (f SearchFilters) FilterTorrents(torrents []Torrent) []Torrent {

	var filtered []Torrent

	for _, torrent := range torrents {

		if f.IsOk(&torrent) {

			filtered = append(filtered, torrent)
		}
	}

	return filtered
}

// FilterTorrentsCount filters the given list of torrents, returning only the ones that
// comply with the filters, while also limiting the result to the number specified by count.
func (f SearchFilters) FilterTorrentsCount(torrents []Torrent, count uint) []Torrent {

	var filtered []Torrent

	for _, torrent := range f.FilterTorrents(torrents) {

		filtered = append(filtered, torrent)

		if count > 0 && uint(len(filtered)) >= count {
			break
		}
	}

	return filtered
}

// PickVideoTorrent functions similar to SearchTorrentList(), but instead returns the torrent with the best available video quality
// with at least one seeder.
func PickVideoTorrent(torrents []Torrent, filters SearchFilters) (*Torrent, error) {

	trnts, err := SearchVideoTorrentList(torrents, filters)

	if err != nil {
		return nil, err
	}

	ok := func(t *Torrent) bool {
		return (filters.MaxQuality == "" || !t.VideoQuality.BetterThan(filters.MaxQuality)) &&
			(filters.MinQuality == "" || !t.VideoQuality.WorseThan(filters.MinQuality))
	}

	var torrent *Torrent

	if t, exists := trnts[UHD]; exists && t.Seeders > 0 && ok(t) {

		torrent = t

	} else if t, exists := trnts[High]; exists && t.Seeders > 0 && ok(t) {

		torrent = t

	} else if t, exists := trnts[Medium]; exists && t.Seeders > 0 && ok(t) {

		torrent = t

	} else if t, exists := trnts[Low]; exists && t.Seeders > 0 && ok(t) {

		torrent = t

	} else if t, exists := trnts[Default]; exists && t.Seeders > 0 && ok(t) {

		torrent = t
	}

	return torrent, nil
}

// SearchVideoTorrentList will find the first torrent in the list for each video quality, that also match the given filters.
// Since it returns one torrent for each known quality, the MinQuality and MaxQuality of the given filters are ignored.
// Returns nil if none are found.
func SearchVideoTorrentList(torrents []Torrent, filters SearchFilters) (map[VideoQuality]*Torrent, error) {

	sort.Sort(sortBySeeds(torrents))

	trnts := make(map[VideoQuality]*Torrent)

	fetch := func(q VideoQuality) error {

		filters.MinQuality = q
		filters.MaxQuality = q

		torrent, err := SearchTorrentList(torrents, filters)

		if err == nil && torrent != nil {
			trnts[q] = torrent
		}

		return err
	}

	var err error

	if err == nil {
		err = fetch(Default)
	}
	if err == nil {
		err = fetch(Low)
	}
	if err == nil {
		err = fetch(Medium)
	}
	if err == nil {
		err = fetch(High)
	}
	if err == nil {
		err = fetch(UHD)
	}

	return trnts, err
}

// SearchTorrentList will return the first torrent in the list that matches the given filters, returning nil if none is found.
func SearchTorrentList(torrents []Torrent, filters SearchFilters) (*Torrent, error) {

	for _, t := range torrents {

		if filters.IsOk(&t) {

			return &t, nil

		}
	}

	return nil, nil
}
