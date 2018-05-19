package main

import (
	"encoding/json"
	"git.gmantaos.com/haath/Gorrent/pkg/piratebay"
	"testing"
)

func TestExecute(t *testing.T) {

	var cmd MirrorsCommand
	Options.JSON = true

	output := CaptureCommand(func() { cmd.Execute(nil) })

	var mirrors []piratebay.Mirror
	json.Unmarshal([]byte(output), &mirrors)

	Options.JSON = false
}
