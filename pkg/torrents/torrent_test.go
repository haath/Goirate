package torrents

import (
	"strings"
	"testing"
)

var torrentURLTests = []struct {
	inMirror  string
	inTorrent string
	out       string
}{
	{"https://pirateproxy.sh/", "/torrent/22274951/The.Expanse.S03E07.PROPER.720p.HDTV.x264-AVS", "https://pirateproxy.sh/torrent/22274951/The.Expanse.S03E07.PROPER.720p.HDTV.x264-AVS"},
	{"https://pirateproxy.sh", "/torrent/22274951/The.Expanse.S03E07.PROPER.720p.HDTV.x264-AVS", "https://pirateproxy.sh/torrent/22274951/The.Expanse.S03E07.PROPER.720p.HDTV.x264-AVS"},
	{"https://pirateproxy.sh", "/torrent/19431416/Windows_10_Pro_v.1709_En-US_(64-bit)_ACTiVATED-HOBBiT", "https://pirateproxy.sh/torrent/19431416/Windows_10_Pro_v.1709_En-US_(64-bit)_ACTiVATED-HOBBiT"},
}

var torrentPeersTest = []struct {
	inSeeders int
	inLeeches int
	out       string
}{
	{1, 2, "1 / 3"},
	{0, 3, "0 / 3"},
}

var torrentSizeTest = []struct {
	in  int64
	out string
}{
	{1, "1.0 KB"},
	{1000, "1.0 MB"},
	{1234, "1.2 MB"},
	{1264, "1.3 MB"},
	{1000000, "1.0 GB"},
	{1500000, "1.5 GB"},
	{1520000, "1.5 GB"},
}

func TestFullURL(t *testing.T) {
	for _, tt := range torrentURLTests {
		t.Run(tt.inTorrent, func(t *testing.T) {
			s := Torrent{
				MirrorURL: tt.inMirror, TorrentURL: tt.inTorrent,
			}
			if s.FullURL() != tt.out {
				t.Errorf("\ngot: %q\nwant %q", s.FullURL(), tt.out)
			}
		})
	}
}

func TestPeersString(t *testing.T) {
	for _, tt := range torrentPeersTest {
		t.Run(tt.out, func(t *testing.T) {
			s := Torrent{
				Seeders: tt.inSeeders, Leeches: tt.inLeeches,
			}
			if s.PeersString() != tt.out {
				t.Errorf("\ngot: %q\nwant %q", s.PeersString(), tt.out)
			}
		})
	}
}

func TestSizeString(t *testing.T) {
	for _, tt := range torrentSizeTest {
		t.Run(tt.out, func(t *testing.T) {
			s := Torrent{Size: tt.in}
			if s.SizeString() != tt.out {
				t.Errorf("\ngot: %q\nwant %q", s.SizeString(), tt.out)
			}
		})
	}
}

func TestMarshalJSON(t *testing.T) {

	tor := Torrent{
		Title:      "some torrent",
		Size:       1200000,
		MirrorURL:  "base_url",
		TorrentURL: "torrent_url",
	}

	json, err := tor.MarshalJSON()

	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(string(json), "base_url/torrent_url") || !strings.Contains(string(json), "1.2 GB") {

		t.Errorf("error unmarshaling json: %v", string(json))
	}
}
