package piratebay

import (
	"testing"
)

var torrentUrlTests = []struct {
	inMirror  string
	inTorrent string
	out       string
}{
	{"https://pirateproxy.sh/", "/torrent/22274951/The.Expanse.S03E07.PROPER.720p.HDTV.x264-AVS", "https://pirateproxy.sh/torrent/22274951/The.Expanse.S03E07.PROPER.720p.HDTV.x264-AVS"},
	{"https://pirateproxy.sh", "/torrent/22274951/The.Expanse.S03E07.PROPER.720p.HDTV.x264-AVS", "https://pirateproxy.sh/torrent/22274951/The.Expanse.S03E07.PROPER.720p.HDTV.x264-AVS"},
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
	for _, tt := range torrentUrlTests {
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
