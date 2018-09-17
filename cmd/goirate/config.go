package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"git.gmantaos.com/haath/Goirate/pkg/series"
	"git.gmantaos.com/haath/Goirate/pkg/torrents"
	"github.com/BurntSushi/toml"
)

// Config holds the global goirate configuration
var Config struct {
	torrents.SearchFilters
	TVDBCredentials series.TVDBCredentials `toml:"tvdb"`
}

// ConfigCommand defines the config command and holds its options.
type ConfigCommand struct {
	torrents.SearchFilters
}

// Execute is the callback of the config command.
func (cmd *ConfigCommand) Execute(args []string) error {

	ImportConfig()

	ApplyFilters(cmd.SearchFilters)

	ExportConfig()

	log.Printf("Updated configuration at %v\n", configPath())

	return nil
}

// ApplyFilters applies the filters that have been specified
// on the given filters object into the Config variable.
func ApplyFilters(filters torrents.SearchFilters) {

	applyFilters(&Config.SearchFilters, &filters)
}

// ApplyConfig applies the Config variable to the given filters object.
func ApplyConfig(filters torrents.SearchFilters) {

	applyFilters(&filters, &Config.SearchFilters)
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

	file, err := os.OpenFile(configPath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	encoder := toml.NewEncoder(file)

	encoder.Encode(Config)
}

func applyFilters(dst *torrents.SearchFilters, src *torrents.SearchFilters) {
	dst.VerifiedUploader = src.VerifiedUploader
	if src.MinQuality != "" {
		dst.MinQuality = src.MinQuality
	}
	if src.MaxQuality != "" {
		dst.MaxQuality = src.MaxQuality
	}
	if src.MinSize != "" {
		dst.MinSize = src.MinSize
	}
	if src.MaxSize != "" {
		dst.MaxSize = src.MaxSize
	}
	dst.MinSeeders = src.MinSeeders

	for _, name := range src.Uploaders.Whitelist {
		dst.Uploaders.Whitelist = append(dst.Uploaders.Whitelist, name)
	}
	for _, name := range src.Uploaders.Blacklist {
		dst.Uploaders.Blacklist = append(dst.Uploaders.Blacklist, name)
	}
}

func configPath() string {

	return path.Join(configDir(), "config.toml")
}
