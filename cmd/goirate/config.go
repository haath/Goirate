package main

import (
	"io/ioutil"
	"log"
	"os/user"
	"path"

	"github.com/BurntSushi/toml"
)

// Config holds the global goirate configuration
var Config struct {
	Uploaders struct {
		Whitelist []string
		Blacklist []string
	}
}

func configPath() string {

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return path.Join(usr.HomeDir, ".goirate/config.toml")
}

// ImportConfig the configuration from the toml file onto the Config variable
func ImportConfig() {

	tomlBytes, err := ioutil.ReadFile(configPath())

	if err != nil {
		log.Fatal(err)
	}

	tomlString := string(tomlBytes)

	if _, err := toml.Decode(tomlString, &Config); err != nil {
		log.Fatal(err)
	}
}
