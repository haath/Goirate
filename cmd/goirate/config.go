package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"gitlab.com/haath/goirate/pkg/series"
	"gitlab.com/haath/goirate/pkg/torrents"
	"gitlab.com/haath/goirate/pkg/utils"
)

// Config holds the global goirate configuration
var Config struct {
	torrents.SearchFilters
	KodiMediaPaths    bool                   `toml:"kodi_media_paths"`
	TPBMirrors        torrents.MirrorFilters `toml:"tpb_mirrors"`
	TVDBCredentials   series.TVDBCredentials `toml:"tvdb"`
	QBittorrentConfig QBittorrentConfig      `toml:"qbittorrent"`
	SMTPConfig        SMTPConfig             `toml:"smtp"`
	Watchlist         utils.WatchlistActions `toml:"actions"`
	DownloadDir       struct {
		General string `toml:"general"`
		Movies  string `toml:"movies"`
		Series  string `toml:"series"`
		Music   string `toml:"music"`
	} `toml:"download_dirs"`
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
// This function also replaces nil values with default and then exports the configuration
// before returning. It is implemented this way, so that any execution of the program
// generates the default configuration, including all of its properties, as opposed
// to the default behavior of BurntSushi/toml which excludes nil values.
// This way a user can begin manually configuring the application through config.toml immediately.
func ImportConfig() {

	if _, err := os.Stat(configPath()); err == nil {

		/*
			Import config.toml
		*/
		tomlBytes, err := ioutil.ReadFile(configPath())

		if err != nil {
			log.Fatal(err)
		}

		tomlString := string(tomlBytes)

		if _, err := toml.Decode(tomlString, &Config); err != nil {
			log.Fatal(err)
		}

		/*
			Property-setting closures
		*/
		setOrDefault := func(val *string, env string, defaultVal string) {
			if os.Getenv(env) != "" {
				*val = os.Getenv(env)
			} else if *val == "" {
				*val = defaultVal
			}
		}
		setOrDefaultQuality := func(val *torrents.VideoQuality, env string, defaultVal torrents.VideoQuality) {
			if os.Getenv(env) != "" {
				*val = torrents.VideoQuality(os.Getenv(env))
			} else if *val == "" {
				*val = defaultVal
			}
		}
		setOrDefaultInt := func(val *int, env string, defaultVal int) {
			if os.Getenv(env) != "" {
				num, err := strconv.ParseInt(os.Getenv(env), 10, 32)
				if err != nil {
					log.Fatal(err)
				}
				*val = int(num)
			} else if *val == 0 {
				*val = defaultVal
			}
		}
		setOrDefaultUint := func(val *uint16, env string, defaultVal uint16) {
			if os.Getenv(env) != "" {
				num, err := strconv.ParseUint(os.Getenv(env), 10, 16)
				if err != nil {
					log.Fatal(err)
				}
				*val = uint16(num)
			} else if *val == 0 {
				*val = defaultVal
			}
		}
		setOptionalBool := func(val *utils.OptionalBoolean, env string, defaultVal utils.OptionalBoolean) {
			if os.Getenv(env) == "true" {
				*val = utils.True
			} else if *val == "" {
				*val = defaultVal
			}
		}
		setBool := func(val *bool, env string) {
			if os.Getenv(env) == "true" {
				*val = true
			}
		}

		/*
			Search filters
		*/
		setBool(&Config.VerifiedUploader, "GOIRATE_VERIFIED_UPLOADER")
		setOrDefaultQuality(&Config.MinQuality, "GOIRATE_MIN_QUALITY", "")
		setOrDefaultQuality(&Config.MaxQuality, "GOIRATE_MAX_QUALITY", "")
		setOrDefault(&Config.MinSize, "GOIRATE_MIN_SIZE", "")
		setOrDefault(&Config.MaxSize, "GOIRATE_MAX_SIZE", "")
		setOrDefaultInt(&Config.MinSeeders, "GOIRATE_MIN_SEEDERS", 0)

		/*
			Download directory options
		*/
		var defaultDownloadsDir string
		usr, usrErr := user.Current()
		if usrErr == nil {
			defaultDownloadsDir = path.Join(usr.HomeDir, "Downloads")
		} else {
			defaultDownloadsDir = path.Join("~", "Downloads")
		}
		setOrDefault(&Config.DownloadDir.General, "GOIRATE_DOWNLOADS_DIR", defaultDownloadsDir)
		setOrDefault(&Config.DownloadDir.Movies, "GOIRATE_DOWNLOADS_MOVIES", defaultDownloadsDir)
		setOrDefault(&Config.DownloadDir.Series, "GOIRATE_DOWNLOADS_SERIES", defaultDownloadsDir)
		setOrDefault(&Config.DownloadDir.Music, "GOIRATE_DOWNLOADS_MUSIC", defaultDownloadsDir)

		/*
			Transmission RPC configurations
		*/
		setOrDefault(&Config.QBittorrentConfig.URL, "GOIRATE_QBT_URL", "http://localhost:8080")
		setOrDefault(&Config.QBittorrentConfig.Username, "GOIRATE_QBT_USERNAME", "")
		setOrDefault(&Config.QBittorrentConfig.Password, "GOIRATE_QBT_PASSWORD", "")

		/*
			SMTP configurations
		*/
		setOrDefault(&Config.SMTPConfig.Host, "GOIRATE_SMTP_HOST", "smtp.gmail.com")
		setOrDefaultUint(&Config.SMTPConfig.Port, "GOIRATE_SMTP_PORT", 587)
		setOrDefault(&Config.SMTPConfig.Username, "GOIRATE_SMTP_USERNAME", "")
		setOrDefault(&Config.SMTPConfig.Password, "GOIRATE_SMTP_PASSWORD", "")

		/*
			Watchlist options
		*/
		if os.Getenv("GOIRATE_ACTIONS_NOTIFY") != "" {

			Config.Watchlist.Emails = strings.Split(os.Getenv("GOIRATE_ACTIONS_NOTIFY"), ",")

		} else if Config.Watchlist.Emails == nil {

			Config.Watchlist.Emails = []string{}
		}
		setOptionalBool(&Config.Watchlist.SendEmail, "GOIRATE_ACTIONS_EMAIL", "")
		setOptionalBool(&Config.Watchlist.Download, "GOIRATE_ACTIONS_DOWNLOAD", "")

		/*
			Pirate Bay mirror filters
		*/
		if Config.TPBMirrors.Whitelist == nil {
			Config.TPBMirrors.Whitelist = []string{}
		}
		if Config.TPBMirrors.Blacklist == nil {
			Config.TPBMirrors.Blacklist = []string{}
		}

		/*
			Misc.
		*/
		setBool(&Config.KodiMediaPaths, "GOIRATE_KODI_MEDIA_PATHS")
	}

	ExportConfig()
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

// GetMirrorScraper returns a scraper for Pirate Bay mirrors, with the appropriate
// configuration passed to it from the Config variable.
func GetMirrorScraper() torrents.MirrorScraper {

	return torrents.MirrorScraper{
		MirrorFilters: Config.TPBMirrors,
	}
}
