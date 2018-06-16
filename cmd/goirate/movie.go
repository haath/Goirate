package main

import (
	"git.gmantaos.com/haath/Goirate/pkg/torrents"
)

// MovieCommand defines the movie command and holds its options.
type MovieCommand struct {
	torrents.SearchFilters
	torrentSearchArgs
}

// Execute is the callback of the movie command.
func (m *MovieCommand) Execute(args []string) error {

	return nil
}
