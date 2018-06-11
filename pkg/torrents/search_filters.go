package torrents

import (
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
