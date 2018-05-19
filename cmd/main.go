package main

import (
	"github.com/jessevdk/go-flags"
)

// Options holds the command line options for the cli program
var Options struct {
	// Options
	Verbose bool `short:"v" long:"verbose" description:"Show more information"`
	JSON    bool `short:"j" long:"json" description:"Output in JSON format"`

	// Commands
	Mirrors MirrorsCommand `command:"mirrors" alias:"m" description:"Get a list of PirateBay mirrors"`
}

func main() {

	parser := flags.NewParser(&Options, flags.HelpFlag|flags.PassDoubleDash|flags.PrintErrors)

	parser.Parse()
}
