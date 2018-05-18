package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
)

func main() {

	var opts struct {
		Verbose bool `short:"v" long:"verbose" description:"Show more information"`
	}

	flags.Parse(&opts)

	mirrors := GetMirrors()
	fmt.Println(len(mirrors))

}
