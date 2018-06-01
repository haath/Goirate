package main

import (
	"github.com/jessevdk/go-flags"
	"log"
)

// Options holds the command line options for the cli program
var Options struct {
	// Options
	Verbose bool `short:"v" long:"verbose" description:"Show more information"`
	JSON    bool `short:"j" long:"json" description:"Output in JSON format"`

	// Commands
	Mirrors MirrorsCommand `command:"mirrors" alias:"m" description:"Get a list of PirateBay mirrors"`
	Search  SearchCommand  `command:"search" alias:"s" description:"Search for torrents"`
}

func main() {

	log.SetFlags(0)

	parser := flags.NewParser(&Options, flags.HelpFlag|flags.PassDoubleDash|flags.PrintErrors)

	parser.Parse()
}
