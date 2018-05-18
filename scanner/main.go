package main

import "github.com/jessevdk/go-flags"

func Sum(x int, y int) int {
	return x + y
}

func main() {

	var opts struct {
		Verbose bool `short:"v" long:"verbose" description:"Show more information"`
	}

	flags.Parse(&opts)
}
