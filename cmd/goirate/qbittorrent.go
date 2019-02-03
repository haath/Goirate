package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jaredlmosley/go-qbittorrent/qbt"
)

// QBittorrentConfig holds the configuration and credentials for communicating with the
// transmission daemon RPC service.
type QBittorrentConfig struct {
	URL      string `toml:"url"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

// QBittorrentClient is a simple wrapper of the underlying transmissionrpc.Client.
type QBittorrentClient struct {
	*qbt.Client
}

// GetQBittorrentClient returns a transmission RPC client connected with the given configuration.
func (cfg *QBittorrentConfig) GetClient() (client *QBittorrentClient, err error) {

	qb := qbt.NewClient(cfg.URL)

	if cfg.Username != "" {

		var loggedIn bool

		loggedIn, err = qb.Login(cfg.Username, cfg.Password)

		if os.Getenv("GOIRATE_DEBUG") == "true" {

			log.Printf("qBittorrent login: %t\n", loggedIn)
		}
	}

	client = &QBittorrentClient{Client: qb}

	return
}

// AddTorrent sends the given magnet link to the transmission daemon and begins its download,
// configuring the downloaded files to be placed at the specified output directory.
func (client *QBittorrentClient) AddTorrent(magnetLink, downloadDir string) error {

	options := map[string]string{
		"savepath": downloadDir,
	}

	resp, err := client.DownloadFromLink(magnetLink, options)

	if err != nil {

		return err
	}

	if resp.StatusCode != 200 {

		err = fmt.Errorf("qBittorrent HTTP error: %s", resp.Status)
	}

	return err
}
