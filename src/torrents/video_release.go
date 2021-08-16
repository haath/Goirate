package torrents

import (
	"sort"
	"strings"
)

/*
	Source: https://en.wikipedia.org/wiki/Pirated_movie_release_types#Release_formats
*/

// VideoRelease defines the release type of a video torrent.
type VideoRelease string

const (
	// Cam is a copy made in a cinema using a camcorder.
	Cam VideoRelease = "Cam"

	// Telesync is a bootleg recording of a film recorded in a movie theater.
	Telesync VideoRelease = "Telesync"

	// Workprint is a copy made from an unfinished version of a film produced by the studio.
	Workprint VideoRelease = "Workprint"

	// Telecine is a opy captured from a film print using a machine that transfers the movie from its analog reel to digital format.
	Telecine VideoRelease = "Telecine"

	// PPVRip come from Pay-Per-View sources.
	PPVRip VideoRelease = "Pay-Per-View Rip"

	// Screener are early DVD or BD releases of the theatrical version of a film, typically sent to movie reviewers.
	Screener VideoRelease = "Screener"

	// DDC stands for Digital Distribution Copy and is basically the same as a Screener, but sent digitally to companies instead of via the postal system.
	DDC VideoRelease = "Digital Distribution Copy"

	// R5 is a studio produced unmastered telecine put out quickly and cheaply to compete against telecine piracy in Russia.
	R5 VideoRelease = "R5"

	// DVDRip is a final retail version of a film, typically released before it is available outside its originating region.
	DVDRip VideoRelease = "DVD-Rip"

	// DVDR is a final retail version of a film in DVD format, generally a complete copy from the original DVD.
	DVDR VideoRelease = "DVD-R"

	// TVRip is a capture source from an analog capture card (coaxial/composite/s-video connection).
	TVRip VideoRelease = "HDTV, PDTV or DSRip"

	// VODRip stands for Video-On-Demand Rip and is recorded or captured from an On-Demand service such as through a cable or satellite TV service.
	VODRip VideoRelease = "VODRip"

	// WEBDL is a file losslessly ripped from a streaming service.
	WEBDL VideoRelease = "WEB-DL"

	// WEBRip is captured from a streaming service, similarly to WEBRip, but by recording the web video stream instead of directly downloading the video.
	WEBRip VideoRelease = "WEBRip"

	// WEBCap is a rip created by capturing video from a DRM-enabled streaming service.
	WEBCap VideoRelease = "WEBCap"

	// BDRip are encoded directly from a Blu-ray disc.
	BDRip VideoRelease = "Blu-ray"
)

var releaseLabels = map[VideoRelease][]string{
	Cam:       {"CAMRip", "CAM"},
	Telesync:  {"TS", "HDTS", "TELESYNC", "PDVD", "PreDVDRip"},
	Workprint: {"WP", "WORKPRINT"},
	Telecine:  {"TC", "HDTC", "TELECINE"},
	PPVRip:    {"PPV", "PPVRip"},
	Screener:  {"SCR", "SCREENER", "DVDSCR", "DVDSCREENER", "BDSCR"},
	DDC:       {"DDC"},
	R5:        {"R5", "R5.LINE", "R5.AC3.5.1.HQ"},
	DVDRip:    {"DVDRip", "DVDMux", "Xvid"},
	DVDR:      {"DVDR", "DVD-Full", "Full-Rip", "ISO rip", "lossless rip", "untouched rip", "DVD-5", "DVD-9"},
	TVRip:     {"DSR", "DSRip", "SATRip", "DTHRip", "DVBRip", "HDTV", "PDTV", "DTVRip", "TVRip", "HDTVRip"},
	VODRip:    {"VODRip", "VODR"},
	WEBDL:     {"WEBDL", "WEB DL", "WEB-DL", "HDRip", "WEB-DLRip"},
	WEBRip:    {"WEBRip", "WEB Rip", "WEB-Rip", "WEB"},
	WEBCap:    {"WEB-Cap", "WEBCAP", "WEB Cap"},
	BDRip:     {"Blu-Ray", "BluRay", "BLURAY", "BDRip", "BRRip", "BDMV", "BDR", "BD25", "BD50", "BD5", "BD9", "BR-rip"},
}

// ExtractVideoRelease parses a torrent's title and returns its video release type, if it exists.
func ExtractVideoRelease(torrentTitle string) VideoRelease {

	torrentTitle = strings.ToLower(torrentTitle)

	labelReleaseMap := map[string]VideoRelease{}
	var sortedLaebels []string

	// Create a mapping of label -> release
	for release, labels := range releaseLabels {

		for _, label := range labels {

			labelReleaseMap[label] = release
			sortedLaebels = append(sortedLaebels, label)
		}
	}

	// Sort labels in descending order
	sort.Slice(sortedLaebels, func(i, j int) bool {
		return len(sortedLaebels[i]) > len(sortedLaebels[j])
	})

	for _, label := range sortedLaebels {

		if strings.Contains(torrentTitle, strings.ToLower(label)) {

			return labelReleaseMap[label]
		}
	}

	return ""
}
