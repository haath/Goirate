package torrents

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
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

var trackers = []string{
	"udp://tracker.coppersurfer.tk:6969/announce",
	"udp://9.rarbg.to:2920/announce",
	"udp://tracker.opentrackr.org:1337",
	"udp://tracker.internetwarriors.net:1337/announce",
	"udp://tracker.leechers-paradise.org:6969/announce",
	"udp://tracker.coppersurfer.tk:6969/announce",
	"udp://tracker.pirateparty.gr:6969/announce",
	"udp://tracker.cyberia.is:6969/announce",
}

// GetTorrents converts the response from the PirateBay API into a list of torrents.
func (response PirateBayAPIResponse) GetTorrents(mirrorURL *url.URL) []Torrent {

	var trnts []Torrent

	for _, obj := range response {

		seeders, _ := strconv.ParseInt(obj.Seeders, 10, 32)
		leechers, _ := strconv.ParseInt(obj.Leechers, 10, 32)
		verifiedUploader := strings.ToLower(obj.Status) == "vip" || strings.ToLower(obj.Status) == "trusted"
		sizeBytes, _ := strconv.ParseInt(obj.Size, 10, 32)

		mirrorSchemeHost := fmt.Sprintf("%v://%v", mirrorURL.Scheme, mirrorURL.Host)
		torrentURL := fmt.Sprintf("/description.php?id=%v", obj.ID)

		addedTimeInt, _ := strconv.ParseInt(obj.Added, 10, 64)

		torrent := Torrent{
			Title:            obj.Name,
			Size:             sizeBytes / 1000,
			Seeders:          int(seeders),
			Leeches:          int(leechers),
			VerifiedUploader: verifiedUploader,
			VideoQuality:     extractVideoQuality(obj.Name),
			VideoRelease:     ExtractVideoRelease(obj.Name),
			MirrorURL:        mirrorSchemeHost,
			TorrentURL:       torrentURL,
			Magnet:           obj.getMagnetLink(),
			UploadTime:       time.Unix(addedTimeInt, 0),
			Uploader:         obj.Username,
		}

		trnts = append(trnts, torrent)
	}

	return trnts
}

func (torrent PirateBayAPIResponseTorrent) getMagnetLink() string {

	title := strings.ReplaceAll(torrent.Name, " ", "+")
	title = strings.Replace(title, "%2B", "+", -1)

	magnet, _ := url.Parse("magnet:")
	query := magnet.Query()
	query.Add("dn", title)

	for _, tracker := range trackers {

		query.Add("tr", tracker)
	}

	magnet.RawQuery = query.Encode() + "&xt=urn:btih:" + torrent.InfoHash

	return magnet.String()
}
