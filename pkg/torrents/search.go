package torrents

import (
	"strings"

	"git.gmantaos.com/haath/gobytes"
)

// SearchFilters holds conditions and filters, used to search for specific torrents.
type SearchFilters struct {
	VerifiedUploader bool         `long:"trusted" description:"Only consider torrents where the uploader is either VIP or Trusted."`
	MinQuality       VideoQuality `long:"min-quality" description:"Minimum acceptable torrent quality (inclusive)."`
	MaxQuality       VideoQuality `long:"max-quality" description:"Maximum acceptable torrent quality (inclusive)."`
	MinSize          string       `long:"min-size" description:"Minimum acceptable torrent size."`
	MaxSize          string       `long:"max-size" description:"Maximum acceptable torrent size."`
	MinSeeders       int          `long:"min-seeders" description:"Minimum acceptable amount of seeders."`
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

	maxSize, err := filters.MaxSizeKB()

	if err != nil {
		return nil, err
	}

	minSize, err := filters.MinSizeKB()

	if err != nil {
		return nil, err
	}

	for _, t := range torrents {

		if filters.VerifiedUploader && !t.VerifiedUploader {
			continue
		}

		if (t.Size > maxSize && maxSize > 0) ||
			(t.Size < minSize && minSize > 0) {
			continue
		}

		if (filters.MinQuality != "" && t.VideoQuality.WorseThan(filters.MinQuality)) ||
			(filters.MaxQuality != "" && t.VideoQuality.BetterThan(filters.MaxQuality)) {
			continue
		}

		return &t, nil
	}

	return nil, nil
}

func normalizeQuery(query string) string {

	replaces := []struct {
		old string
		new string
	}{
		{"-", " "},
		{"'", " "},
		{".", " "},
		{"_", " "},
		{":", ""},
		{"!", ""},
		{"(", ""},
		{")", ""},
	}

	query = strings.TrimSpace(query)

	for _, rep := range replaces {
		query = strings.Replace(query, rep.old, rep.new, -1)
	}

	return query
}
