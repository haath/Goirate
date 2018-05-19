package main

import (
	"encoding/json"
	"git.gmantaos.com/haath/Gorrent/pkg/piratebay"
	"log"
)

// MirrorsCommand defines the mirrors command and holds its options.
type MirrorsCommand struct {
}

// Execute acts as the call back of the mirrors command.
func (m *MirrorsCommand) Execute(args []string) error {
	mirrors := piratebay.GetMirrors()

	if Options.JSON {
		mirrorsJSON, _ := json.MarshalIndent(mirrors, "", "   ")
		log.Println(mirrorsJSON)
	}

	for _, mirror := range mirrors {
		status := "x"
		if !mirror.Status {
			status = " "
		}

		log.Printf("[%s] %s %s\n", status, mirror.Country, mirror.URL)
	}

	return nil
}
