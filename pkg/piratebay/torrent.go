package piratebay

import (
	"time"
)

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

// Torrent holds all the information regarding a torrent.
type Torrent struct {
	Title            string       `json:"title"`
	Size             int          `json:"size"` // In kilobytes
	Seeders          int          `json:"seeders"`
	Leechers         int          `json:"leechers"`
	VerifiedUploader bool         `json:"verified_uploader"`
	VideoQuality     VideoQuality `json:"video_quality"`
	URL              string       `json:"url"`
	Magnet           string       `json:"magnet"`
	UploadTime       time.Time    `json:"upload_time"`
}
