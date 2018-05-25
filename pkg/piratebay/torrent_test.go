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
