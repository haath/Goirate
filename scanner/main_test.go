package main

import (
	"encoding/json"
	"testing"
)

func TestMirrors(t *testing.T) {
	var response []Mirror

	mirrorsJSON := mirrors()

	if err := json.Unmarshal([]byte(mirrorsJSON), &response); err != nil {
		t.Error(err)
	}
}
