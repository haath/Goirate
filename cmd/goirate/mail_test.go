package main

import (
	"log"
	"testing"

	"git.gmantaos.com/haath/Goirate/pkg/series"
	"git.gmantaos.com/haath/Goirate/pkg/torrents"
)

func TestLoadTorrentTemplate(t *testing.T) {

	torrents := seriesTorrents{
		Series: &series.Series{Title: "SUper special show"},
		Torrents: []seriesTorrent{
			seriesTorrent{Torrent: torrents.Torrent{MirrorURL: "localhost", TorrentURL: "my/torrent", VerifiedUploader: true}, Episode: series.Episode{Season: 1, Episode: 1}},
			seriesTorrent{Torrent: torrents.Torrent{MirrorURL: "localhost", TorrentURL: "my/torrent"}, Episode: series.Episode{Season: 1, Episode: 2}},
		},
	}

	tmpl, err := LoadTorrentTemplate(torrents)

	if err != nil {
		t.Error(err)
	}

	log.Print(tmpl)
}

func TestSendEmail(t *testing.T) {

	resetConfigs()

	ImportConfig()

	inbox := "goirate-test@mailinator.com"

	err := Config.SMTPConfig.SendEmail("Test e-mail", "<b>Some</b> body", inbox)

	if err != nil {
		t.Error(err)
	}

	resetConfigs()
}
