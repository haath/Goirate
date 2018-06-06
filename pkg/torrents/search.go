package torrents

import (
	"errors"
)

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

	return nil, errors.New("No torrent found with the specified filters")
}
