package torrents

import (
	"strconv"
	"strings"
)

// PirateBayAPIResponseTorrent represents a torrents, as it is returned by the PirateBay API.
type PirateBayAPIResponseTorrent struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	InfoHash string `json:"info_hash"`
	Leechers string `json:"leechers"`
	Seeders  string `json:"seeders"`
	NumFiles string `json:"num_files"`
	Size     string `json:"size"`
	Username string `json:"username"`
	Added    string `json:"added"`
	Status   string `json:"status"`
	Category string `json:"category"`
	IMDB     string `json:"imdb"`
}

// PirateBayAPIResponse represents the response returned by the PirateBay API.
type PirateBayAPIResponse []PirateBayAPIResponseTorrent

// GetTorrents converts the response from the PirateBay API into a list of torrents.
func (response PirateBayAPIResponse) GetTorrents() []Torrent {

	var trnts []Torrent

	for _, obj := range response {

		seeders, _ := strconv.ParseInt(obj.Seeders, 10, 32)
		leechers, _ := strconv.ParseInt(obj.Leechers, 10, 32)
		verifiedUploader := strings.ToLower(obj.Status) == "vip" || strings.ToLower(obj.Status) == "trusted"
		sizeBytes, _ := strconv.ParseInt(obj.Size, 10, 32)

		torrent := Torrent{
			Title:            obj.Name,
			VideoQuality:     extractVideoQuality(obj.Name),
			VideoRelease:     ExtractVideoRelease(obj.Name),
			VerifiedUploader: verifiedUploader,
			Uploader:         obj.Username,
			Size:             sizeBytes / 1000,
			Seeders:          int(seeders),
			Leeches:          int(leechers),
		}

		trnts = append(trnts, torrent)
	}

	return trnts
}
