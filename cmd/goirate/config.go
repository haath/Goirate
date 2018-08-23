package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"

	"git.gmantaos.com/haath/Goirate/pkg/torrents"
	"github.com/BurntSushi/toml"
)

// Config holds the global goirate configuration
var Config struct {
	VerifiedUploader bool                  `toml:"trusted"`
	MinQuality       torrents.VideoQuality `toml:"min-quality"`
	MaxQuality       torrents.VideoQuality `toml:"max-quality"`
	MinSize          string                `toml:"min-size"`
	MaxSize          string                `toml:"max-size"`
	MinSeeders       int                   `toml:"min-seeders"`
	Uploaders        struct {
		Whitelist []string `toml:"whitelist"`
		Blacklist []string `toml:"blacklist"`
	} `toml:"uploaders"`
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

	return nil
}

// ApplyFilters applies the filters that have been specified
// on the given filters object into the Config variable.
func ApplyFilters(filters torrents.SearchFilters) {

	Config.VerifiedUploader = filters.VerifiedUploader
	if filters.MinQuality != "" {
		Config.MinQuality = filters.MinQuality
	}
	if filters.MaxQuality != "" {
		Config.MaxQuality = filters.MaxQuality
	}
	if filters.MinSize != "" {
		Config.MinSize = filters.MinSize
	}
	if filters.MaxSize != "" {
		Config.MaxSize = filters.MaxSize
	}
	Config.MinSeeders = filters.MinSeeders

	if len(Config.Uploaders.Whitelist) == 0 {
		Config.Uploaders.Whitelist = filters.UploaderWhitelist
	}
	if len(Config.Uploaders.Blacklist) == 0 {
		Config.Uploaders.Blacklist = filters.UploaderBlacklist
	}
}

// ApplyConfig applies the Config variable to the given filters object.
func ApplyConfig(filters torrents.SearchFilters) {
	filters.VerifiedUploader = Config.VerifiedUploader
	if Config.MinQuality != "" {
		filters.MinQuality = Config.MinQuality
	}
	if Config.MaxQuality != "" {
		filters.MaxQuality = Config.MaxQuality
	}
	if Config.MinSize != "" {
		filters.MinSize = Config.MinSize
	}
	if Config.MaxSize != "" {
		filters.MaxSize = Config.MaxSize
	}
	filters.MinSeeders = Config.MinSeeders

	if len(filters.UploaderWhitelist) == 0 {
		filters.UploaderWhitelist = Config.Uploaders.Whitelist
	}
	if len(filters.UploaderBlacklist) == 0 {
		filters.UploaderBlacklist = Config.Uploaders.Blacklist
	}
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

func configPath() string {

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return path.Join(usr.HomeDir, ".goirate/config.toml")
}
