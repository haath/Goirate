package main

import (
	"encoding/json"
	"fmt"
	"github.com/jessevdk/go-flags"
)

func main() {

	var opts struct {
		Verbose bool `short:"v" long:"verbose" description:"Show more information"`

		Mirrors func() `short:"m" long:"mirrors" description:"Get a list of PirateBay mirrors"`
	}

	opts.Mirrors = func() {
		fmt.Println(mirrors())
	}

	flags.Parse(&opts)

}

func mirrors() string {
	mirrors := GetMirrors()
	mirrorsJSON, _ := json.MarshalIndent(mirrors, "", "   ")
	return string(mirrorsJSON)
}
