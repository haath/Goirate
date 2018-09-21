package main

import (
	"strings"
	"testing"
	"time"

	"git.gmantaos.com/haath/Goirate/pkg/series"
	"git.gmantaos.com/haath/Goirate/pkg/torrents"
)

func TestLoadTorrentTemplate(t *testing.T) {

	now := time.Now()
	exp := now.Format("02/01/2006")

	torrents := seriesTorrents{
		Series: &series.Series{Title: "SUper special show"},
		Torrents: []seriesTorrent{
			seriesTorrent{Torrent: torrents.Torrent{MirrorURL: "localhost", TorrentURL: "my/torrent", VerifiedUploader: true}, Episode: series.Episode{Season: 1, Episode: 1}},
			seriesTorrent{Torrent: torrents.Torrent{MirrorURL: "localhost", TorrentURL: "my/torrent"}, Episode: series.Episode{Season: 1, Episode: 2}},
			seriesTorrent{Torrent: torrents.Torrent{MirrorURL: "localhost", TorrentURL: "my/torrent"}, Episode: series.Episode{Season: 1, Episode: 2, Title: "Episode Title", Aired: &now}},
		},
	}

	tmpl, err := LoadTorrentTemplate(torrents)

	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(tmpl, exp) {
		t.Errorf("Template does not contain date: %v\n%v", exp, tmpl)
	}
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
