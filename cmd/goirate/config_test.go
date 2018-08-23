package main

import (
	"reflect"
	"testing"
)

func TestImportExport(t *testing.T) {

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
}
