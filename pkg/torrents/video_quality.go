package torrents

// VideoQuality defines the standard video qualities for torrents.
type VideoQuality string

const (
	// Default acts as the default quality when no other is found in the title.
	Default VideoQuality = "HDTV"

	// Low represents the 480p quality.
	Low VideoQuality = "480p"

	// Medium represents the 720p quality.
	Medium VideoQuality = "720p"

	// High represents the 1080p quality.
	High VideoQuality = "1080p"

	// UHD represents the 2160p quality.
	UHD VideoQuality = "2160p"
)

// WorseThan will return true if the quality passed as an argument is
// worse than this one.
func (q VideoQuality) WorseThan(quality VideoQuality) bool {
	return q.numeric() < quality.numeric()
}

// BetterThan will return true if the quality passed as an argument is
// better than this one.
func (q VideoQuality) BetterThan(quality VideoQuality) bool {
	return q.numeric() > quality.numeric()
}

func (q VideoQuality) numeric() int {
	switch q {
	case Default:
		return 0
	case Low:
		return 1
	case Medium:
		return 2
	case High:
		return 3
	case UHD:
		return 4
	default:
		return 0
	}
}
