package torrents

import (
	"git.gmantaos.com/haath/gobytes"
	"strings"
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

// SearchTorrentList will return the best torrent in the list that matches the given filters,
// returning nil if none is found.
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
