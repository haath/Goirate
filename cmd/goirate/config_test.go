package main

import (
	"reflect"
	"testing"

	"gitlab.com/haath/Goirate/pkg/torrents"
)

func TestImportExpor(t *testing.T) {

	resetConfigs()

	whitelist := []string{"allowed_user1", "allowed_user2"}
	blacklist := []string{"banned_user", "bad_boye"}

	Config.Uploaders.Whitelist = whitelist
	Config.Uploaders.Blacklist = blacklist

	ExportConfig()

	Config.Uploaders.Whitelist = []string{}
	Config.Uploaders.Blacklist = []string{}

	ImportConfig()

	if !reflect.DeepEqual(Config.Uploaders.Whitelist, whitelist) {
		t.Errorf("\ngot %v\nwant %v", Config.Uploaders.Whitelist, whitelist)
	}

	resetConfigs()
}

func TestExecute(t *testing.T) {

	resetConfigs()

	var cmd ConfigCommand

	cmd.MaxQuality = torrents.Medium
	cmd.MinQuality = torrents.Low
	cmd.VerifiedUploader = true
	cmd.MinSize = "12 GB"
	cmd.Uploaders.Whitelist = []string{"allowed_user1", "allowed_user2"}
	cmd.Uploaders.Blacklist = []string{"banned_user", "bad_boye"}

	_, err := CaptureCommand(cmd.Execute)

	if err != nil {
		t.Fatal(err)
	}

	ImportConfig()

	if !reflect.DeepEqual(Config.SearchFilters, cmd.SearchFilters) {
		t.Errorf("\ngot %v\nwant %v", Config, cmd)
	}

	resetConfigs()
}

func resetConfigs() {

	Config.SearchFilters = torrents.SearchFilters{}
	Config.Uploaders.Whitelist = []string{}
	Config.Uploaders.Blacklist = []string{}

	ExportConfig()
}
