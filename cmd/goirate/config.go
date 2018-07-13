package main

import (
	"io/ioutil"
	"log"
	"os"
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

// ImportConfig the configuration from config.toml onto the Config variable
func ImportConfig() {

	if _, err := os.Stat(configPath()); err == nil {

		tomlBytes, err := ioutil.ReadFile(configPath())

		if err != nil {
			log.Fatal(err)
		}

		tomlString := string(tomlBytes)

		if _, err := toml.Decode(tomlString, &Config); err != nil {
			log.Fatal(err)
		}

	}
}

// ExportConfig writes the current configuration to the config.toml file
func ExportConfig() {

	os.MkdirAll(path.Dir(configPath()), os.ModePerm)

	file, err := os.OpenFile(configPath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)

	if err != nil {
		log.Fatal(err)
	}

	encoder := toml.NewEncoder(file)

	encoder.Encode(Config)
}
