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
)
