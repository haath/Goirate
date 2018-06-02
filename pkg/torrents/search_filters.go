package torrents

import (
	"git.gmantaos.com/haath/gobytes"
)

// SearchFilters holds conditions and filters, used to search for specific torrents.
type SearchFilters struct {
	Query            string
	VerifiedUploader bool
	MinQuality       VideoQuality
	MaxQuality       VideoQuality
	MinSize          string
	MaxSize          string
	MinSeeders       int
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
