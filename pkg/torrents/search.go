package torrents

import (
	"git.gmantaos.com/haath/gobytes"
)

// SearchFilters holds conditions and filters, used to search for specific torrents.
type SearchFilters struct {
	VerifiedUploader  bool         `long:"trusted" description:"Only consider torrents where the uploader is either VIP or Trusted."`
	MinQuality        VideoQuality `long:"min-quality" description:"Minimum acceptable torrent quality (inclusive)."`
	MaxQuality        VideoQuality `long:"max-quality" description:"Maximum acceptable torrent quality (inclusive)."`
	MinSize           string       `long:"min-size" description:"Minimum acceptable torrent size."`
	MaxSize           string       `long:"max-size" description:"Maximum acceptable torrent size."`
	MinSeeders        int          `long:"min-seeders" description:"Minimum acceptable amount of seeders."`
	UploaderWhitelist []string     `long:"whitelist" description:"Add to a whitelist of uploaders, to only consider torrents from them."`
	UploaderBlacklist []string     `long:"blacklist" description:"Add to a blacklist of uploaders, to avoid torrents from them."`
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

	return (len(f.UploaderBlacklist) == 0 || !contains(f.UploaderBlacklist, uploader)) &&
		(len(f.UploaderWhitelist) == 0 || contains(f.UploaderWhitelist, uploader))
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

	return true
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

	if t, exists := trnts[High]; exists && t.Seeders > 0 && ok(t) {
		return t, nil
	}
	if t, exists := trnts[Medium]; exists && t.Seeders > 0 && ok(t) {
		return t, nil
	}
	if t, exists := trnts[Low]; exists && t.Seeders > 0 && ok(t) {
		return t, nil
	}
	if t, exists := trnts[Default]; exists && t.Seeders > 0 && ok(t) {
		return t, nil
	}

	return nil, nil
}

// SearchVideoTorrentList will find the first torrent in the list for each video quality, that also match the given filters.
// Since it returns one torrent for each known quality, the MinQuality and MaxQuality of the given filters are ignored.
// Returns nil if none are found.
func SearchVideoTorrentList(torrents []Torrent, filters SearchFilters) (map[VideoQuality]*Torrent, error) {

	trnts := make(map[VideoQuality]*Torrent)

	fetch := func(q VideoQuality) error {
		filters.MinQuality = q
		filters.MaxQuality = q

		torrent, err := SearchTorrentList(torrents, filters)

		if err != nil {
			return err
		}

		if torrent != nil {
			trnts[q] = torrent
		}

		return nil
	}

	if err := fetch(Default); err != nil {
		return nil, err
	}
	if err := fetch(Low); err != nil {
		return nil, err
	}
	if err := fetch(Medium); err != nil {
		return nil, err
	}
	if err := fetch(High); err != nil {
		return nil, err
	}

	return trnts, nil
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
