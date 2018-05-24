package main

import (
	"encoding/json"
	"git.gmantaos.com/haath/Gorrent/pkg/piratebay"
	"testing"
)

func TestSearchExecute(t *testing.T) {

	var cmd SearchCommand
	Options.JSON = true

	output := CaptureCommand(func() { cmd.Execute([]string{"avengers"}) })

	var mirrors []piratebay.Mirror
	json.Unmarshal([]byte(output), &mirrors)

	Options.JSON = false
}
