package piratebay

import (
	"time"
)

// VideoQuality defines the standard video qualities for torrents.
type VideoQuality string

const (
	// HDTV acts as the default quality when no other is found in the title.
	HDTV VideoQuality = "HDTV"
	// LOW represents the 480p quality.
	LOW VideoQuality = "480p"
	// MEDIUM represents the 720p quality.
	MEDIUM VideoQuality = "720p"
	// HIGH represents the 1080p quality.
	HIGH VideoQuality = "1080p"
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
