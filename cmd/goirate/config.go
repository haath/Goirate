package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"

	"git.gmantaos.com/haath/Goirate/pkg/series"
	"git.gmantaos.com/haath/Goirate/pkg/torrents"
	"github.com/BurntSushi/toml"
)

// Config holds the global goirate configuration
var Config struct {
	torrents.SearchFilters
	TVDBCredentials series.TVDBCredentials `toml:"tvdb"`
	RPCConfig       RPCConfig              `toml:"transmission_rpc"`
	DownloadDir     struct {
		General string `toml:"general"`
		Movies  string `toml:"movies"`
		Series  string `toml:"series"`
		Music   string `toml:"music"`
	} `toml:"download_dirs"`
	Watchlist struct {
		Email    bool `toml:"email"`
		Download bool `toml:"download"`
	} `toml:"watchlist"`
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
func ApplyConfig(filters *torrents.SearchFilters) {

	applyFilters(filters, &Config.SearchFilters)
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

		if Config.RPCConfig.Host == "" {
			Config.RPCConfig = DefaultTransmissionRPCConfig()
		}

		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}

		setOrDefault := func(val *string, env string) {
			if os.Getenv(env) != "" {
				*val = os.Getenv(env)
			} else if *val == "" {
				*val = path.Join(usr.HomeDir, "Downloads")
			}
		}

		setOrDefault(&Config.DownloadDir.General, "GOIRATE_DOWNLOADS_DIR")
		setOrDefault(&Config.DownloadDir.Movies, "GOIRATE_DOWNLOADS_MOVIES")
		setOrDefault(&Config.DownloadDir.Series, "GOIRATE_DOWNLOADS_SERIES")
		setOrDefault(&Config.DownloadDir.Music, "GOIRATE_DOWNLOADS_MUSIC")
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
	dst.VerifiedUploader = dst.VerifiedUploader || src.VerifiedUploader
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
